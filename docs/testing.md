# Testing

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test ./... -cover

# Run library tests only
go test .

# Run CLI tests only
go test ./cmd/inkcheck

# Run semantic tests only (requires model)
go test . -run Semantic

# Skip semantic tests (faster, no model needed)
go test . -skip Semantic
```

## Test Data

Test files are located in `testdata/` and are licensed under CC0 1.0 Universal (Public Domain).

## Semantic Tests

Semantic tests will automatically skip if the word embedding model is not available. To download the model:

```bash
inkcheck model download
```
