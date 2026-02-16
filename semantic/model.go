package semantic

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

const (
	defaultModelFile = "numberbatch-en-19.08.txt.gz"

	// modelDownloadTimeout is the maximum time allowed for downloading the model.
	modelDownloadTimeout = 30 * time.Minute

	// bufferSize is the buffer size for reading model files (1MB).
	bufferSize = 1 << 20

	// maxReasonableWords prevents memory exhaustion from corrupted model files.
	maxReasonableWords = 10_000_000
)

// ModelManager holds a loaded word embedding model.
type ModelManager struct {
	embeddings map[string][]float32
	dim        int
}

var (
	// ErrModelNotFound is returned by LoadModel when the model file does not exist.
	ErrModelNotFound = fmt.Errorf("model not found")
	// ErrModelCorrupt is returned by LoadModel when the file exists but cannot be loaded.
	ErrModelCorrupt = fmt.Errorf("model corrupt")
)

// ModelPath returns the resolved path to the model file.
// Checks INKCHECK_MODEL_PATH env var first, then uses config.
func ModelPath(cfg config.Config) (string, error) {
	path := os.Getenv("INKCHECK_MODEL_PATH")
	if path != "" {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, cfg.ModelDir, defaultModelFile), nil
}

// LoadModel loads the word embedding model from disk.
// Returns ErrModelNotFound if the file does not exist.
func LoadModel(cfg config.Config) (*ModelManager, error) {
	path, err := ModelPath(cfg)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrModelNotFound, path)
	}

	m, err := loadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrModelCorrupt, path, err)
	}
	return m, nil
}

// DownloadModel downloads the word embedding model to the configured path.
func DownloadModel(cfg config.Config) error {
	return DownloadModelWithContext(context.Background(), cfg)
}

// DownloadModelWithContext downloads the word embedding model with cancellation support.
func DownloadModelWithContext(ctx context.Context, cfg config.Config) error {
	path, err := ModelPath(cfg)
	if err != nil {
		return err
	}
	return downloadModelWithContext(ctx, path, cfg.ModelURL)
}

func loadFromFile(path string) (*ModelManager, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var reader io.Reader
	if strings.HasSuffix(path, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("gzip open: %w", err)
		}
		defer gz.Close()
		reader = bufio.NewReaderSize(gz, bufferSize)
	} else {
		reader = bufio.NewReaderSize(f, bufferSize)
	}

	return readWord2VecText(reader)
}

// readWord2VecText reads word2vec text format:
//
//	first line: num_words dimensions
//	subsequent lines: word float1 float2 ... floatN
func readWord2VecText(r io.Reader) (*ModelManager, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, bufferSize), bufferSize)

	// Read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty model file")
	}
	header := strings.Fields(scanner.Text())
	if len(header) < 2 {
		return nil, fmt.Errorf("invalid header: %q", scanner.Text())
	}
	numWords, err := strconv.Atoi(header[0])
	if err != nil {
		return nil, fmt.Errorf("invalid word count in header: %w", err)
	}
	if numWords > maxReasonableWords {
		return nil, fmt.Errorf("suspiciously large word count: %d (max %d)", numWords, maxReasonableWords)
	}
	dim, err := strconv.Atoi(header[1])
	if err != nil {
		return nil, fmt.Errorf("invalid dimension in header: %w", err)
	}

	embeddings := make(map[string][]float32, numWords)
	lineNum := 1

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		// Find the first space to split word from vector
		idx := strings.IndexByte(line, ' ')
		if idx < 0 {
			continue
		}
		word := line[:idx]
		rest := line[idx+1:]

		parts := strings.Fields(rest)
		if len(parts) != dim {
			continue
		}

		vec := make([]float32, dim)
		valid := true
		for i, s := range parts {
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				valid = false
				break
			}
			vec[i] = float32(v)
		}
		if !valid {
			continue
		}

		// Numberbatch uses /c/en/word format for multilingual,
		// but English-only file uses plain words with underscores
		embeddings[strings.ToLower(word)] = vec
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read error at line %d: %w", lineNum, err)
	}

	return &ModelManager{embeddings: embeddings, dim: dim}, nil
}

func downloadModelWithContext(ctx context.Context, path, url string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Downloading word embedding model to %s...\n", path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: modelDownloadTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	tmpPath := path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	// Wrap response body with context-aware reader
	reader := &contextReader{ctx: ctx, r: resp.Body}
	written, err := io.Copy(f, reader)
	f.Close()
	if err != nil {
		_ = os.Remove(tmpPath) // Cleanup on error, ignore remove errors
		return err
	}

	if resp.ContentLength > 0 && written != resp.ContentLength {
		_ = os.Remove(tmpPath) // Cleanup on error, ignore remove errors
		return fmt.Errorf("incomplete download: got %d bytes, expected %d", written, resp.ContentLength)
	}

	fmt.Fprintf(os.Stderr, "Downloaded %d MB\n", written/1024/1024)
	return os.Rename(tmpPath, path)
}

// contextReader wraps an io.Reader and checks context cancellation on each read.
type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr *contextReader) Read(p []byte) (int, error) {
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
		return cr.r.Read(p)
	}
}

// Embedding returns the word vector for a single word.
// Returns nil if the word is not in the vocabulary or the model is nil.
func (m *ModelManager) Embedding(word string) []float32 {
	if m == nil {
		return nil
	}
	vec, ok := m.embeddings[strings.ToLower(word)]
	if !ok {
		return nil
	}
	return vec
}

// SentenceEmbedding returns the average word vector for all words in a sentence.
func (m *ModelManager) SentenceEmbedding(sentence string) []float32 {
	return m.averageEmbedding(shared.ListWords(sentence))
}

// ParagraphEmbedding returns the average word vector for all words in a paragraph.
func (m *ModelManager) ParagraphEmbedding(paragraph string) []float32 {
	return m.averageEmbedding(shared.ListWords(paragraph))
}

// Similarity returns the cosine similarity between embeddings of two texts.
// Returns 0 if the model is nil.
func (m *ModelManager) Similarity(text1, text2 string) float64 {
	if m == nil {
		return 0
	}
	e1 := m.averageEmbedding(shared.ListWords(text1))
	e2 := m.averageEmbedding(shared.ListWords(text2))
	return shared.CosineSimilarity(e1, e2)
}

func (m *ModelManager) averageEmbedding(words []string) []float32 {
	var sum []float32
	count := 0
	for _, w := range words {
		vec := m.Embedding(w)
		if vec == nil {
			continue
		}
		if sum == nil {
			sum = make([]float32, len(vec))
		}
		for i, v := range vec {
			sum[i] += v
		}
		count++
	}
	if count == 0 {
		return nil
	}
	for i := range sum {
		sum[i] /= float32(count)
	}
	return sum
}
