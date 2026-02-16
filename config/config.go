package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds all configurable thresholds and paths for inkcheck.
type Config struct {
	RedundancyThreshold       float64 `toml:"redundancy_threshold"`
	MATTRWindowSize           int     `toml:"mattr_window_size"`
	UniformityDeviation       float64 `toml:"uniformity_deviation"`
	JargonRankThreshold       int     `toml:"jargon_rank_threshold"`
	TransitionRepeatThreshold int     `toml:"transition_repeat_threshold"`
	OpenerWordCount           int     `toml:"opener_word_count"`
	ReadabilityFormula        string  `toml:"readability_formula"`
	ModelDir                  string  `toml:"model_dir"`
	ModelURL                  string  `toml:"model_url"`
}

// DefaultConfig returns a Config with all current hardcoded defaults.
func DefaultConfig() Config {
	return Config{
		RedundancyThreshold:       0.90,
		MATTRWindowSize:           50,
		UniformityDeviation:       0.2,
		JargonRankThreshold:       6000,
		TransitionRepeatThreshold: 3,
		OpenerWordCount:           2,
		ReadabilityFormula:        "flesch_kincaid_grade",
		ModelDir:                  ".inkcheck/models",
		ModelURL:                  "https://conceptnet.s3.amazonaws.com/downloads/2019/numberbatch/numberbatch-en-19.08.txt.gz",
	}
}

// Path returns the config file path.
//
//	Linux/macOS: ~/.config/inkcheck/config.toml
//	Windows:     %AppData%\inkcheck\config.toml
func Path() (string, error) {
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine config directory: %w", err)
		}
		return filepath.Join(dir, "inkcheck", "config.toml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "inkcheck", "config.toml"), nil
}

// Load reads the config file. If the file does not exist, it returns defaults.
// Returns an error if the path cannot be determined or if the config file
// exists but cannot be parsed.
func Load() (Config, error) {
	cfg := DefaultConfig()

	path, err := Path()
	if err != nil {
		return cfg, fmt.Errorf("getting config path: %w", err)
	}

	_, err = toml.DecodeFile(path, &cfg)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}

func defaultConfigTemplate() string {
	d := DefaultConfig()
	return fmt.Sprintf(`# inkcheck configuration
# https://github.com/inkcheck/inkcheck

# Cosine similarity threshold for detecting redundant paragraphs
redundancy_threshold = %g

# Window size for Moving-Average Type-Token Ratio
mattr_window_size = %d

# Deviation threshold for paragraph position uniformity
uniformity_deviation = %g

# Words ranked above this are considered jargon
jargon_rank_threshold = %d

# Minimum occurrences to flag a transition phrase as repeated
transition_repeat_threshold = %d

# Number of leading words used to identify sentence openers
opener_word_count = %d

# Readability formula to use (flesch_kincaid_grade, gunning_fog, flesch_reading_ease, etc.)
# See: github.com/inkcheck/readability for all available formulas
readability_formula = "%s"

# Word embedding model directory (relative to home)
model_dir = "%s"

# Word embedding model download URL
model_url = "%s"
`, d.RedundancyThreshold, d.MATTRWindowSize, d.UniformityDeviation,
		d.JargonRankThreshold, d.TransitionRepeatThreshold, d.OpenerWordCount,
		d.ReadabilityFormula, d.ModelDir, d.ModelURL)
}

// Init writes the default config to disk. Returns an error if the file
// already exists.
func Init() error {
	path, err := Path()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config already exists: %s", path)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(defaultConfigTemplate()), 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// List formats the config as key=value lines for display, including the file path.
func List(cfg Config) string {
	var b strings.Builder

	path, err := Path()
	if err != nil {
		path = "(unknown)"
	}
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		fmt.Fprintf(&b, "config file: %s (not found, using defaults)\n", path)
	} else {
		fmt.Fprintf(&b, "config file: %s\n", path)
	}

	fmt.Fprintf(&b, "\n")
	fmt.Fprintf(&b, "redundancy_threshold       = %.2f\n", cfg.RedundancyThreshold)
	fmt.Fprintf(&b, "mattr_window_size          = %d\n", cfg.MATTRWindowSize)
	fmt.Fprintf(&b, "uniformity_deviation       = %.2f\n", cfg.UniformityDeviation)
	fmt.Fprintf(&b, "jargon_rank_threshold      = %d\n", cfg.JargonRankThreshold)
	fmt.Fprintf(&b, "transition_repeat_threshold = %d\n", cfg.TransitionRepeatThreshold)
	fmt.Fprintf(&b, "opener_word_count          = %d\n", cfg.OpenerWordCount)
	fmt.Fprintf(&b, "readability_formula        = %s\n", cfg.ReadabilityFormula)
	fmt.Fprintf(&b, "model_dir                  = %s\n", cfg.ModelDir)
	fmt.Fprintf(&b, "model_url                  = %s\n", cfg.ModelURL)

	return b.String()
}
