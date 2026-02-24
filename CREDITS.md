# Credits and Attributions

## Development

Built using [Claude Code](https://claude.ai/claude-code).

## Go Dependencies

All Go module dependencies are listed in `go.mod` and are licensed under MIT or BSD licenses compatible with this project:

- `github.com/BurntSushi/toml` - MIT License
- `github.com/tsawler/prose/v3` - MIT License
- `github.com/yuin/goldmark` - MIT License
- `gonum.org/v1/gonum` - BSD-3-Clause License
- `gopkg.in/neurosnap/sentences.v1` - MIT License

## Data Sources

### ConceptNet Numberbatch

The semantic analysis metrics use word embeddings from ConceptNet Numberbatch (downloaded separately by the user, not bundled with this software).

- **Version**: 19.08
- **License**: CC BY-SA 4.0
- **Source**: https://github.com/commonsense/conceptnet-numberbatch
- **Download URL**: https://conceptnet.s3.amazonaws.com/downloads/2019/numberbatch/
- **Citation**: Robyn Speer, Joshua Chin, and Catherine Havasi (2017). "ConceptNet 5.5: An Open Multilingual Graph of General Knowledge." In proceedings of AAAI 31.

The word embedding model is optional and only required for semantic metrics. The model file itself is licensed under CC BY-SA 4.0 and is not part of this software distribution.

### Google 10,000 English Words

Word frequency data used in vocabulary sophistication analysis (embedded from `wordlist/google-10000-english.txt`).

- **Source**: https://github.com/first20hours/google-10000-english
- **Derived from**: Google Trillion Word Corpus
- **Note**: Word frequency lists are factual data compilations and are used here for linguistic analysis purposes.

## License

This project (Inkcheck) is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

Note that third-party data sources listed above may have their own licenses. Users should comply with the respective licenses when using those components.
