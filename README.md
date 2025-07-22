# subs-cli

A powerful CLI tool for automatically finding and downloading subtitles for your media files.

## Overview

`subs-cli` intelligently parses media filenames to identify TV shows and movies, then searches for matching subtitles in your preferred language through the OpenSubtitles API.

## Features

- 🎯 **Smart Filename Parsing** - Automatically extracts series name, year, season, and episode from filenames
- 🌍 **Multi-language Support** - Download subtitles in any language (default: en)
- 🔍 **Interactive Search** - Use fuzzy finder to select from multiple subtitle options
- 📁 **Batch Processing** - Process entire directories of media files
- ⚡ **Caching** - Local cache to reduce API calls and improve performance
- 🔧 **Configurable** - Store API keys and preferences in config file
- 📊 **Smart Matching** - Prioritizes subtitles matching video quality and release group

## Installation

```bash
go install github.com/yourusername/subs-cli@latest
```

Or download pre-built binaries from the [releases page](https://github.com/yourusername/subs-cli/releases).

## Quick Start

```bash
# Find subtitles for current directory
subs .

# Find Portuguese (Brazil) subtitles
subs . --language pt-BR

# Interactive mode - choose from search results
subs "Dark Matter" --language pt-BR --interactive

# Process specific file
subs Dark.Matter.2024.S01E01.1080p.x265-ELiTE.mkv

# Batch process with custom config
subs /path/to/media --config ~/.subs-cli/config.yaml
```

## Configuration

Create a config file at `~/.subs-cli/config.yaml`:

```yaml
# OpenSubtitles API configuration
opensubtitles:
  api_key: your_api_key_here
  username: your_username
  password: your_password

# Default settings
defaults:
  language: pt-BR
  interactive: true
  auto_select: false

# Cache settings
cache:
  enabled: true
  ttl: 24h
  path: ~/.subs-cli/cache
```

## Filename Format

The tool expects media files to follow common naming conventions:

```
Series.Name.Year.SxxExx.Quality.Source.ext
Series.Name.SxxExx.Quality.Source.ext
Movie.Name.Year.Quality.Source.ext
```

Examples:
- `Dark.Matter.2024.S01E01.1080p.x265-ELiTE.mkv`
- `The.Office.S03E07.720p.BluRay.x264.mkv`
- `Inception.2010.1080p.BluRay.x264-SPARKS.mkv`

## API Limits

OpenSubtitles API has the following limits:

- **Free tier**: 20 downloads/day (with free account)
- **VIP tier**: 1000 downloads/day ($15/year)

The tool respects these limits and provides helpful messages when limits are reached.

## Advanced Usage

### Batch Processing

Process an entire season:
```bash
subs /media/series/Dark.Matter.2024.S01/ --language pt-BR
```

### Custom Patterns

For non-standard filenames, use manual search:
```bash
subs --search "Dark Matter S01E01" --language pt-BR
```

### Multiple Languages

Download subtitles in multiple languages:
```bash
subs . --language pt-BR,en,es
```

### Dry Run

Preview what would be downloaded:
```bash
subs . --dry-run
```

## Building from Source

```bash
git clone https://github.com/yourusername/subs-cli
cd subs-cli
go build -o subs main.go
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [OpenSubtitles](https://www.opensubtitles.com) for providing the subtitle database
- [Kong](https://github.com/alecthomas/kong) for the excellent CLI framework
- [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) for interactive selection