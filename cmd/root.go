package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
)

type CLI struct {
	Path        string   `arg:"" default:"." help:"Path to media file or directory to search for subtitles. Supports files (.mp4, .mkv, etc.) and directories."`
	Language    []string `short:"l" long:"language" default:"en" help:"Subtitle language codes (ISO 639-1/locale format). Examples: en, pt-BR, es, fr. Supports multiple comma-separated values."`
	Interactive bool     `short:"i" long:"interactive" help:"Enable interactive fuzzy finder mode for subtitle selection. Allows browsing and previewing multiple subtitle options."`
	Config      string   `short:"c" long:"config" type:"existingfile" help:"Path to custom YAML configuration file. Default location: ~/.subs-cli/config.yaml"`
	DryRun      bool     `long:"dry-run" help:"Preview mode: displays what subtitles would be downloaded without actually downloading them. Useful for testing."`
	Search      string   `short:"s" long:"search" help:"Manual search query mode. Use instead of filename parsing (e.g., 'Breaking Bad S01E01'). Overrides path-based search."`
	Version     bool     `short:"v" long:"version" help:"Display detailed version information including build details, Git commit, and platform info."`
}

func (c *CLI) Run() error {
	if c.Version {
		c.printVersionInfo()
		return nil
	}

	if err := c.validateArguments(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	c.displayConfiguration()

	// TODO: Implement subtitle search and download logic
	fmt.Println("\n[TODO] Core subtitle search and download logic not yet implemented")

	return nil
}

func (c *CLI) printVersionInfo() {
	fmt.Printf("subs-cli version %s\n", Version)
	if BuildTime != "unknown" {
		fmt.Printf("Built: %s\n", BuildTime)
	}
	if GitCommit != "unknown" {
		fmt.Printf("Commit: %s\n", GitCommit)
	}
	if GoVersion != "unknown" {
		fmt.Printf("Go version: %s\n", GoVersion)
	} else {
		fmt.Printf("Go version: %s\n", runtime.Version())
	}
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func (c *CLI) validateArguments() error {
	var results []*ValidationResult

	if c.Search == "" {
		result, err := c.validatePath()
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	langResult, err := c.validateLanguages()
	if err != nil {
		return err
	}
	results = append(results, langResult)

	if c.Config != "" {
		configResult, err := c.validateConfigFile()
		if err != nil {
			return err
		}
		results = append(results, configResult)
	}

	modeResult, err := c.validateModeConsistency()
	if err != nil {
		return err
	}
	results = append(results, modeResult)

	c.printValidationResults(results)

	return nil
}

func (c *CLI) printValidationResults(results []*ValidationResult) {
	for _, result := range results {
		if result.Success && result.Message != "" {
			fmt.Printf("✓ %s\n", result.Message)
		}
		if result.Warning != "" {
			fmt.Printf("⚠ Warning: %s\n", result.Warning)
		}
		if result.Message != "" && !result.Success {
			fmt.Printf("ℹ %s\n", result.Message)
		}
	}
}

type ValidationResult struct {
	Success bool
	Message string
	Warning string
}

var mediaExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".m4v":  true,
	".mpg":  true,
	".mpeg": true,
	".3gp":  true,
}

func (c *CLI) validatePath() (*ValidationResult, error) {
	cleanPath := filepath.Clean(c.Path)

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%s': %w", c.Path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", absPath)
		}
		return nil, fmt.Errorf("cannot access path '%s': %w", absPath, err)
	}

	c.Path = absPath

	result := &ValidationResult{Success: true}

	if info.IsDir() {
		result.Message = fmt.Sprintf("Directory path validated: %s", c.Path)
	} else {
		result.Message = fmt.Sprintf("File path validated: %s", c.Path)

		ext := strings.ToLower(filepath.Ext(c.Path))
		if !mediaExtensions[ext] && ext != "" {
			result.Warning = fmt.Sprintf("File extension '%s' may not be a supported media format", ext)
		}
	}

	return result, nil
}

func (c *CLI) validateLanguages() (*ValidationResult, error) {
	if len(c.Language) == 0 {
		return nil, fmt.Errorf("at least one language must be specified")
	}

	validLanguages := make([]string, 0, len(c.Language))

	for _, lang := range c.Language {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}

		if len(lang) < 2 || len(lang) > 5 {
			return nil, fmt.Errorf("invalid language code '%s': must be 2-5 characters (e.g., 'en', 'pt-BR')", lang)
		}

		if !isValidLanguageCode(lang) {
			return nil, fmt.Errorf("invalid language code format '%s': expected format like 'en' or 'pt-BR'", lang)
		}

		validLanguages = append(validLanguages, lang)
	}

	if len(validLanguages) == 0 {
		return nil, fmt.Errorf("no valid language codes provided")
	}

	c.Language = validLanguages
	return &ValidationResult{
		Success: true,
		Message: fmt.Sprintf("Language codes validated: %v", c.Language),
	}, nil
}

func (c *CLI) validateConfigFile() (*ValidationResult, error) {
	absPath, err := filepath.Abs(c.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid config file path '%s': %w", c.Config, err)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", absPath)
		}
		return nil, fmt.Errorf("cannot access config file '%s': %w", absPath, err)
	}

	c.Config = absPath
	return &ValidationResult{
		Success: true,
		Message: fmt.Sprintf("Config file validated: %s", c.Config),
	}, nil
}

func (c *CLI) validateModeConsistency() (*ValidationResult, error) {
	result := &ValidationResult{Success: true}
	var messages []string

	if c.Search != "" {
		if c.Path != "." {
			messages = append(messages, fmt.Sprintf("Manual search mode enabled: path argument '%s' will be ignored", c.Path))
		}

		if strings.TrimSpace(c.Search) == "" {
			return nil, fmt.Errorf("search query cannot be empty when using search mode")
		}
	}

	if c.Interactive {
		messages = append(messages, "Interactive mode enabled: you'll be able to select from multiple subtitle options")
	}

	if c.DryRun {
		messages = append(messages, "Dry run mode: no files will be downloaded, only preview what would happen")
	}

	if len(messages) > 0 {
		result.Message = strings.Join(messages, "\n")
	}

	return result, nil
}

func (c *CLI) displayConfiguration() {
	fmt.Println("\n--- Configuration ---")

	if c.Search != "" {
		fmt.Printf("Mode: Manual search\n")
		fmt.Printf("Search query: %s\n", c.Search)
	} else {
		fmt.Printf("Mode: Path-based search\n")
		fmt.Printf("Target path: %s\n", c.Path)
	}

	fmt.Printf("Languages: %v\n", c.Language)
	fmt.Printf("Interactive: %t\n", c.Interactive)
	fmt.Printf("Dry run: %t\n", c.DryRun)

	if c.Config != "" {
		fmt.Printf("Config file: %s\n", c.Config)
	} else {
		fmt.Printf("Config file: default (~/.subs-cli/config.yaml)\n")
	}
}

func isValidLanguageCode(code string) bool {
	code = strings.ToLower(code)

	if len(code) == 2 || len(code) == 3 {
		for _, r := range code {
			if r < 'a' || r > 'z' {
				return false
			}
		}
		return true
	}

	if len(code) == 5 && code[2] == '-' {
		firstPart := code[:2]
		secondPart := code[3:]

		for _, r := range firstPart {
			if r < 'a' || r > 'z' {
				return false
			}
		}

		for _, r := range secondPart {
			if r < 'a' || r > 'z' {
				return false
			}
		}

		return true
	}

	return false
}

func Execute() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("subs"),
		kong.Description("A powerful CLI tool for automatically finding and downloading subtitles for your media files.\n\n"+
			"Examples:\n"+
			"  subs /path/to/movie.mp4                    # Find subtitles for a specific file\n"+
			"  subs /path/to/movies/ -l en,pt-BR          # Search directory for multiple languages\n"+
			"  subs . -i -l es                           # Interactive mode with Spanish subtitles\n"+
			"  subs --search \"Breaking Bad S01E01\"        # Manual search query\n"+
			"  subs /path/to/series/ --dry-run           # Preview mode without downloading\n"+
			"  subs -c ~/.config/subs.yaml /movies/      # Use custom config file\n\n"+
			"Supported languages: en, es, pt-BR, fr, de, it, ru, ja, ko, zh, and many more.\n"+
			"Use standard ISO 639-1 codes (en) or locale codes (pt-BR, zh-CN)."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: false,
		}),
	)

	err := cli.Run()
	ctx.FatalIfErrorf(err)
}
