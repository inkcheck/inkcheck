# Test Data

This directory contains test data for the inkcheck project.

## Test Files

### Argumentative Essay (`argumentative_essay.md`)
- **Type**: Argumentative/persuasive writing
- **Purpose**: Tests claim/support ratio, counterargument engagement, thesis detection
- **License**: CC0 1.0 Universal (Public Domain)
- **Characteristics**: Contains thesis, evidence, counterarguments, conclusion

### Narrative Prose (`narrative_prose.md`)
- **Type**: Narrative/creative writing
- **Purpose**: Tests figurative language, tension/resolution, voice consistency
- **License**: CC0 1.0 Universal (Public Domain)
- **Characteristics**: Story arc, descriptive language, character perspective

### Technical Writing (`technical_writing.md`)
- **Type**: Technical/expository
- **Purpose**: Tests vocabulary sophistication, jargon detection, specificity
- **License**: CC0 1.0 Universal (Public Domain)
- **Characteristics**: Technical terminology, structured explanations, examples

### Edge Cases (`edge_cases.md`)
- **Type**: Mixed/edge cases
- **Purpose**: Tests handling of unusual text patterns
- **License**: CC0 1.0 Universal (Public Domain)
- **Characteristics**: Single sentences, empty paragraphs, mixed content, code blocks

## Licensing

All test files in this directory are dedicated to the public domain under CC0 1.0 Universal,
except where otherwise noted. See `LICENSE.txt` for full details.

## Adding New Test Data

When adding new test files:

1. Use CC0, public domain, or permissively licensed content
2. Include license information in file footer
3. Choose text that exercises specific metrics
4. Add corresponding test cases in `*_test.go` files
5. Document the purpose and characteristics in this README

## Usage in Tests

These files are used by:
- `inkcheck_test.go` - Library-level end-to-end tests
- `cmd/inkcheck/cli_test.go` - CLI end-to-end tests

Tests verify that metrics produce reasonable values and handle edge cases gracefully,
without asserting exact metric values (which would be brittle).
