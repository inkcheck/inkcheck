# CLI Usage

## Installation

Build from source (macOS):

```bash
CGO_LDFLAGS="-framework Accelerate" go build -o inkcheck ./cmd/inkcheck
```

Build from source (Linux with OpenBLAS):

```bash
apt install libopenblas-dev
go build -o inkcheck ./cmd/inkcheck
```

## Synopsis

```
inkcheck [-m <metric>] [file or directory ...]
inkcheck [-m <metric>] < input.txt
```

## Options

- `-m <metric>` (optional) - the metric to compute. Omit to run all 21 metrics.
  Use `-m help` to see the full list.

## Input

The tool accepts one or more file paths or directory paths as arguments. When
given a directory it walks it recursively and processes files with these
extensions (case-insensitive):

- `.txt`
- `.md`
- `.rst`
- `.adoc`
- `.tex`

Shell glob patterns are also supported (e.g. `*.md`). If no arguments are
provided, text is read from standard input.

## Output

**Single file or stdin, single metric:** prints the metric name, a tab, and the
value.

```
$ inkcheck -m paragraph_variance essay.md
paragraph_variance	0.4523
```

**Single file, all metrics:** prints one line per metric.

```
$ inkcheck essay.md
paragraph_variance	0.4523
sentence_length_variance	0.5102
...
```

**Multiple files:** prints the file path as a prefix on each line.

```
$ inkcheck -m paragraph_variance docs/
docs/intro.md	paragraph_variance	0.3674
docs/conclusion.md	paragraph_variance	0.5210
```

Structured metrics include additional detail in parentheses:

```
$ inkcheck -m punctuation_profile essay.md
punctuation_profile	7/8 types	(. 57 , 58 ; 1 : 3 — 4 () 2 ? 1 ! 0)
```

## Examples

Run all metrics on a file:

```bash
inkcheck essay.md
```

Run a single structural metric:

```bash
inkcheck -m sentence_length_variance essay.md
```

Analyze all Markdown files in a directory:

```bash
inkcheck -m vocabulary_sophistication ./content/
```

Pipe text from another command:

```bash
curl -s https://example.com/article.txt | inkcheck -m hedging_analysis
```

Run only semantic metrics:

```bash
inkcheck -m topic_coherence article.md
inkcheck -m redundancy_detection article.md
```

List available metrics:

```bash
inkcheck -m help
```

## Environment Variables

- `INKCHECK_MODEL_PATH` - path to a custom word2vec binary model file. If not
  set, the model is downloaded to `~/.inkcheck/models/` on first use of a
  semantic metric.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid metric, file not found, model download failure) |
