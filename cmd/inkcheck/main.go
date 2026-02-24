package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/inkcheck/config"
	"github.com/inkcheck/semantic"
)

var textExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".rst":  true,
	".adoc": true,
	".tex":  true,
}

var allMetricNames = []string{
	// structure
	"paragraph_variance",
	"sentence_length_variance",
	"sentence_opener_diversity",
	"sentence_type_distribution",
	"paragraph_position_analysis",
	"punctuation_profile",
	// rhetoric
	"transition_word_density",
	"vocabulary_sophistication",
	"hedging_analysis",
	"specificity_score",
	"voice_consistency",
	"figurative_language",
	"rhetorical_diversity",
	"claim_support_ratio",
	"counterargument_engagement",
	"audience_awareness",
	"argument_structure",
	"tension_and_resolution",
	"stance_analysis",
	"contraction_rate",
	"temporal_orientation",
	"economy_analysis",
	// readability
	"readability",
	// semantic
	"topic_coherence",
	"semantic_progression",
	"redundancy_detection",
	"information_novelty",
	"emotional_tone",
}

var semanticMetrics = map[string]bool{
	"topic_coherence":      true,
	"semantic_progression": true,
	"redundancy_detection": true,
	"information_novelty":  true,
	"emotional_tone":       true,
}

func main() {
	// Handle subcommands before flag parsing.
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "config":
			handleConfig(os.Args[2:])
			return
		case "model":
			handleModel(os.Args[2:])
			return
		case "analyse", "analyze":
			runAnalyse(os.Args[2:])
			return
		case "signature":
			runSignature(os.Args[2:])
			return
		}
	}

	// Default: analyse
	runAnalyse(os.Args[1:])
}

func runAnalyse(cliArgs []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	fs := flag.NewFlagSet("analyse", flag.ExitOnError)
	metric := fs.String("m", "", "metric to compute (omit for all)")
	format := fs.String("format", "text", "output format: text, json")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: inkcheck [analyse] [-m <metric>] [-format <format>] [files or directories...]\n")
		fmt.Fprintf(os.Stderr, "       inkcheck signature [-format <format>] [files or directories...]\n")
		fmt.Fprintf(os.Stderr, "       inkcheck config init|list\n")
		fmt.Fprintf(os.Stderr, "       inkcheck model download\n")
		fmt.Fprintf(os.Stderr, "       command | inkcheck [-m <metric>]\n\n")
		fmt.Fprintf(os.Stderr, "Available metrics:\n")
		for _, name := range allMetricNames {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fs.PrintDefaults()
	}
	fs.Parse(cliArgs)

	if *metric == "help" {
		fs.Usage()
		os.Exit(0)
	}

	if *metric != "" && !validMetric(*metric) {
		fmt.Fprintf(os.Stderr, "Error: unknown metric %q\n\nAvailable metrics:\n", *metric)
		for _, name := range allMetricNames {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		os.Exit(1)
	}

	// Validate format
	outputFormat := OutputFormat(*format)
	if outputFormat != FormatText && outputFormat != FormatJSON {
		fmt.Fprintf(os.Stderr, "Error: unknown format %q (supported: text, json)\n", *format)
		os.Exit(1)
	}

	args := fs.Args()

	// Show usage if no files and no piped input.
	if len(args) == 0 && isTerminal(os.Stdin) {
		fs.Usage()
		os.Exit(0)
	}

	// Load semantic model if needed.
	var model *semantic.ModelManager
	if needsSemanticModel(*metric) {
		model, err = loadSemanticModel(cfg)
		if err != nil {
			os.Exit(1)
		}
	}

	files := resolveInputs(args)
	if files == nil {
		// stdin mode
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		if err := printMetricsWithFormat(outputFormat, os.Stdout, string(data), "", *metric, model, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	multiFile := len(files) > 1
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
			continue
		}
		prefix := ""
		if multiFile {
			prefix = path
		}
		if err := printMetricsWithFormat(outputFormat, os.Stdout, string(data), prefix, *metric, model, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
	}
}

func runSignature(cliArgs []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	fs := flag.NewFlagSet("signature", flag.ExitOnError)
	format := fs.String("format", "text", "output format: text, json")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: inkcheck signature [-format <format>] [files or directories...]\n")
		fmt.Fprintf(os.Stderr, "       command | inkcheck signature [-format <format>]\n\n")
		fmt.Fprintf(os.Stderr, "Computes a 10-axis style signature for the text.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}
	fs.Parse(cliArgs)

	outputFormat := OutputFormat(*format)
	if outputFormat != FormatText && outputFormat != FormatJSON {
		fmt.Fprintf(os.Stderr, "Error: unknown format %q (supported: text, json)\n", *format)
		os.Exit(1)
	}

	args := fs.Args()

	if len(args) == 0 && isTerminal(os.Stdin) {
		fs.Usage()
		os.Exit(0)
	}

	// Signature uses semantic metrics when available; skip if model not found.
	var model *semantic.ModelManager
	model, err = semantic.LoadModel(cfg)
	if err != nil {
		model = nil // signature works without semantic metrics
	}

	files := resolveInputs(args)
	if files == nil {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		if err := printSignatureWithFormat(outputFormat, os.Stdout, string(data), "", model, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	multiFile := len(files) > 1
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
			continue
		}
		prefix := ""
		if multiFile {
			prefix = path
		}
		if err := printSignatureWithFormat(outputFormat, os.Stdout, string(data), prefix, model, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
	}
}

// loadSemanticModel loads the semantic model, prompting for download if needed.
func loadSemanticModel(cfg config.Config) (*semantic.ModelManager, error) {
	model, err := semantic.LoadModel(cfg)
	if errors.Is(err, semantic.ErrModelNotFound) {
		if !isTerminal(os.Stdin) {
			fmt.Fprintf(os.Stderr, "Semantic model not found. Run 'inkcheck model download' first.\n")
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "Semantic model not found. Download now? (~310 MB) [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Fprintf(os.Stderr, "Skipping. Run 'inkcheck model download' to download later.\n")
			return nil, errors.New("model download declined")
		}
		if err := semantic.DownloadModel(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading model: %v\n", err)
			return nil, err
		}
		model, err = semantic.LoadModel(cfg)
	} else if errors.Is(err, semantic.ErrModelCorrupt) {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'inkcheck model download --force' to re-download.\n")
		return nil, err
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading semantic model: %v\n", err)
		return nil, err
	}
	return model, nil
}

// resolveInputs expands args (files, directories, globs) into a list of file paths.
// Returns nil if args is empty (meaning stdin should be used).
func resolveInputs(args []string) []string {
	if len(args) == 0 {
		return nil
	}

	var files []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			matches, globErr := filepath.Glob(arg)
			if globErr != nil || len(matches) == 0 {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			files = append(files, matches...)
			continue
		}
		if info.IsDir() {
			if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && textExtensions[strings.ToLower(filepath.Ext(path))] {
					files = append(files, path)
				}
				return nil
			}); err != nil {
				fmt.Fprintf(os.Stderr, "Error walking %s: %v\n", arg, err)
				os.Exit(1)
			}
		} else {
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No files found.\n")
		os.Exit(1)
	}
	return files
}

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: inkcheck config <init|list>\n")
		os.Exit(1)
	}
	switch args[0] {
	case "init":
		if err := config.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		path, err := config.Path()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config written to %s\n", path)
	case "list":
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		fmt.Print(config.List(cfg))
	default:
		fmt.Fprintf(os.Stderr, "Unknown config command %q\nUsage: inkcheck config <init|list>\n", args[0])
		os.Exit(1)
	}
}

func handleModel(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: inkcheck model <download [-force]>\n")
		os.Exit(1)
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}
	switch args[0] {
	case "download":
		fs := flag.NewFlagSet("model download", flag.ExitOnError)
		force := fs.Bool("force", false, "re-download model even if it exists")
		fs.Parse(args[1:])
		path, err := semantic.ModelPath(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if _, err := os.Stat(path); err == nil && !*force {
			fmt.Fprintf(os.Stderr, "Model already exists: %s\n", path)
			fmt.Fprintf(os.Stderr, "Use -force to re-download.\n")
			os.Exit(0)
		}
		if *force {
			os.Remove(path)
		}
		if err := semantic.DownloadModel(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Model downloaded to %s\n", path)
	default:
		fmt.Fprintf(os.Stderr, "Unknown model command %q\nUsage: inkcheck model <download [-force]>\n", args[0])
		os.Exit(1)
	}
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func validMetric(name string) bool {
	for _, m := range allMetricNames {
		if m == name {
			return true
		}
	}
	return false
}

func needsSemanticModel(metric string) bool {
	if metric == "" {
		return true // "all" includes semantic metrics
	}
	return semanticMetrics[metric]
}

