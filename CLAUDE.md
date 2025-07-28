# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building
```bash
go build -o subs main.go
```

### Running Tests
```bash
go test -v ./...
go test -race ./...
```

### Module Management
```bash
go mod tidy
go mod download
```

### Running the CLI
```bash
go run main.go --help
go run main.go . --language pt-BR --interactive
```

## Architecture Overview

This is a Go CLI application built with the Kong framework for parsing subtitles for media files. The project follows a clean architecture pattern with clear separation of concerns:

### Project Structure
```
subs-cli/
├── cmd/root.go          # CLI command definitions and Kong setup
├── internal/            # Private application packages
│   ├── api/            # OpenSubtitles API client
│   ├── cache/          # Subtitle caching layer
│   ├── config/         # YAML configuration management  
│   ├── interactive/    # Fuzzy finder UI for subtitle selection
│   └── parser/         # Media filename parsing logic
├── pkg/models/         # Public data models and types
└── main.go            # Application entry point
```

### Key Technologies
- **CLI Framework**: Kong for command-line argument parsing and validation
- **HTTP Client**: Resty for OpenSubtitles API communication
- **Configuration**: Koanf for YAML config file management with mapstructure
- **Interactive UI**: go-fuzzyfinder for subtitle selection (planned)
- **Caching**: BigCache for subtitle response caching (planned)

### Core Components

#### CLI Structure (cmd/root.go)
The CLI uses Kong with the following main flags:
- `Path` (positional): Media file or directory path 
- `Language` (-l): Subtitle language(s), defaults to "en"
- `Interactive` (-i): Enable interactive subtitle selection
- `Config` (-c): Custom config file path
- `DryRun`: Preview mode without downloading
- `Search` (-s): Manual search query mode

#### Internal Packages
- **api/**: OpenSubtitles API client with rate limiting and authentication
- **parser/**: Regex-based filename parsing for TV shows and movies
- **cache/**: Local caching to reduce API calls and improve performance
- **config/**: YAML configuration loading with defaults
- **interactive/**: Fuzzy finder interface for subtitle selection

#### Data Models (pkg/models/)
Core types for subtitles, search parameters, and media information that are shared across packages.

## Development Notes

### Current Implementation Status
The project is in early development with basic CLI structure implemented. The main `Run()` function currently only prints debug information - the core subtitle search and download logic is marked as TODO.

### Expected Implementation Flow
1. Parse media filename to extract title, season, episode, quality
2. Search OpenSubtitles API with parsed metadata
3. Present results via interactive fuzzy finder (if enabled)
4. Download and save selected subtitle files
5. Cache responses to minimize API usage

### Configuration File
The application expects a YAML config at `~/.subs-cli/config.yaml` with:
- OpenSubtitles API credentials
- Default language preferences  
- Cache settings
- Interactive mode preferences

### Error Handling
Should implement comprehensive error handling for:
- Invalid filenames that don't match parsing patterns
- API rate limiting and authentication failures
- Network connectivity issues
- File permission errors during subtitle downloads

### Testing Strategy
Focus testing on:
- Filename parsing with various media naming conventions
- API client with mocked responses
- Configuration loading and validation
- Interactive selection logic