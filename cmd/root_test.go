package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		inputPath   string
		expectError bool
		errorMsg    string
		expectMsg   string
		expectWarn  string
	}{
		{
			name: "valid_file_with_media_extension",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "movie.mp4")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
				return tmpFile
			},
			expectError: false,
			expectMsg:   "File path validated:",
		},
		{
			name: "valid_directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
			expectMsg:   "Directory path validated:",
		},
		{
			name: "non_existent_path",
			setupFunc: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent", "path")
			},
			expectError: true,
			errorMsg:    "path does not exist:",
		},
		{
			name: "file_with_non_media_extension",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "document.txt")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
				return tmpFile
			},
			expectError: false,
			expectMsg:   "File path validated:",
			expectWarn:  "File extension '.txt' may not be a supported media format",
		},
		{
			name: "relative_path_resolution",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "subdir")
				require.NoError(t, os.Mkdir(subDir, 0755))
				require.NoError(t, os.Chdir(tmpDir))
				return "subdir"
			},
			expectError: false,
			expectMsg:   "Directory path validated:",
		},
		{
			name: "path_with_spaces",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "my movie file.mkv")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
				return tmpFile
			},
			expectError: false,
			expectMsg:   "File path validated:",
		},
		{
			name: "path_with_special_characters",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "movie[2024].avi")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
				return tmpFile
			},
			expectError: false,
			expectMsg:   "File path validated:",
		},
		{
			name: "all_supported_media_extensions",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
			expectMsg:   "Directory path validated:",
		},
	}

	for ext := range mediaExtensions {
		ext := ext
		tests = append(tests, struct {
			name        string
			setupFunc   func(t *testing.T) string
			inputPath   string
			expectError bool
			errorMsg    string
			expectMsg   string
			expectWarn  string
		}{
			name: "media_extension_" + strings.TrimPrefix(ext, "."),
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "video"+ext)
				require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))
				return tmpFile
			},
			expectError: false,
			expectMsg:   "File path validated:",
		})
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.inputPath
			if tt.setupFunc != nil {
				path = tt.setupFunc(t)
			}

			cli := &CLI{Path: path}
			result, err := cli.validatePath()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.True(t, result.Success)

				if tt.expectMsg != "" {
					assert.Contains(t, result.Message, tt.expectMsg)
				}

				if tt.expectWarn != "" {
					assert.Equal(t, tt.expectWarn, result.Warning)
				}
			}
		})
	}
}

func TestValidateLanguages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		languages   []string
		expectError bool
		errorMsg    string
		expected    []string
	}{
		{
			name:        "valid_single_language",
			languages:   []string{"en"},
			expectError: false,
			expected:    []string{"en"},
		},
		{
			name:        "valid_multiple_languages",
			languages:   []string{"en", "es", "fr"},
			expectError: false,
			expected:    []string{"en", "es", "fr"},
		},
		{
			name:        "valid_locale_format",
			languages:   []string{"pt-BR", "zh-CN"},
			expectError: false,
			expected:    []string{"pt-BR", "zh-CN"},
		},
		{
			name:        "mixed_formats",
			languages:   []string{"en", "pt-BR", "es"},
			expectError: false,
			expected:    []string{"en", "pt-BR", "es"},
		},
		{
			name:        "three_letter_code",
			languages:   []string{"eng", "spa"},
			expectError: false,
			expected:    []string{"eng", "spa"},
		},
		{
			name:        "language_with_spaces",
			languages:   []string{" en ", " pt-BR "},
			expectError: false,
			expected:    []string{"en", "pt-BR"},
		},
		{
			name:        "empty_language_list",
			languages:   []string{},
			expectError: true,
			errorMsg:    "at least one language must be specified",
		},
		{
			name:        "empty_string_language",
			languages:   []string{"", "en", ""},
			expectError: false,
			expected:    []string{"en"},
		},
		{
			name:        "all_empty_strings",
			languages:   []string{"", "", ""},
			expectError: true,
			errorMsg:    "no valid language codes provided",
		},
		{
			name:        "invalid_too_short",
			languages:   []string{"e"},
			expectError: true,
			errorMsg:    "invalid language code 'e': must be 2-5 characters",
		},
		{
			name:        "invalid_too_long",
			languages:   []string{"english"},
			expectError: true,
			errorMsg:    "invalid language code 'english': must be 2-5 characters",
		},
		{
			name:        "invalid_format_numbers",
			languages:   []string{"e1"},
			expectError: true,
			errorMsg:    "invalid language code format 'e1'",
		},
		{
			name:        "invalid_format_special_chars",
			languages:   []string{"en!"},
			expectError: true,
			errorMsg:    "invalid language code format 'en!'",
		},
		{
			name:        "invalid_locale_format",
			languages:   []string{"en_US"},
			expectError: true,
			errorMsg:    "invalid language code format 'en_US'",
		},
		{
			name:        "invalid_locale_too_short",
			languages:   []string{"e-BR"},
			expectError: true,
			errorMsg:    "invalid language code format 'e-BR'",
		},
		{
			name:        "invalid_locale_too_long",
			languages:   []string{"eng-BR"},
			expectError: true,
			errorMsg:    "invalid language code 'eng-BR': must be 2-5 characters",
		},
		{
			name:        "case_insensitive",
			languages:   []string{"EN", "PT-br", "Es"},
			expectError: false,
			expected:    []string{"EN", "PT-br", "Es"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli := &CLI{Language: tt.languages}
			result, err := cli.validateLanguages()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Equal(t, tt.expected, cli.Language)
				assert.Contains(t, result.Message, "Language codes validated:")
			}
		})
	}
}

func TestValidateConfigFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid_config_file",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "config.yaml")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test: value"), 0644))
				return tmpFile
			},
			expectError: false,
		},
		{
			name: "non_existent_config",
			setupFunc: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.yaml")
			},
			expectError: true,
			errorMsg:    "config file does not exist:",
		},
		{
			name: "relative_path_config",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tmpFile := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test: value"), 0644))
				require.NoError(t, os.Chdir(tmpDir))
				return "config.yaml"
			},
			expectError: false,
		},
		{
			name: "config_with_spaces_in_path",
			setupFunc: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "my config.yaml")
				require.NoError(t, os.WriteFile(tmpFile, []byte("test: value"), 0644))
				return tmpFile
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			configPath := ""
			if tt.setupFunc != nil {
				configPath = tt.setupFunc(t)
			}

			cli := &CLI{Config: configPath}
			result, err := cli.validateConfigFile()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.True(t, result.Success)
				assert.Contains(t, result.Message, "Config file validated:")
				assert.True(t, filepath.IsAbs(cli.Config))
			}
		})
	}
}

func TestValidateModeConsistency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cli         CLI
		expectError bool
		errorMsg    string
		expectMsgs  []string
	}{
		{
			name: "search_mode_with_custom_path",
			cli: CLI{
				Search: "Breaking Bad S01E01",
				Path:   "/custom/path",
			},
			expectError: false,
			expectMsgs:  []string{"Manual search mode enabled: path argument '/custom/path' will be ignored"},
		},
		{
			name: "search_mode_with_default_path",
			cli: CLI{
				Search: "Breaking Bad S01E01",
				Path:   ".",
			},
			expectError: false,
			expectMsgs:  []string{},
		},
		{
			name: "empty_search_query",
			cli: CLI{
				Search: "   ",
				Path:   ".",
			},
			expectError: true,
			errorMsg:    "search query cannot be empty when using search mode",
		},
		{
			name: "interactive_mode",
			cli: CLI{
				Interactive: true,
			},
			expectError: false,
			expectMsgs:  []string{"Interactive mode enabled: you'll be able to select from multiple subtitle options"},
		},
		{
			name: "dry_run_mode",
			cli: CLI{
				DryRun: true,
			},
			expectError: false,
			expectMsgs:  []string{"Dry run mode: no files will be downloaded, only preview what would happen"},
		},
		{
			name: "all_modes_combined",
			cli: CLI{
				Search:      "Breaking Bad",
				Path:        "/movies",
				Interactive: true,
				DryRun:      true,
			},
			expectError: false,
			expectMsgs: []string{
				"Manual search mode enabled: path argument '/movies' will be ignored",
				"Interactive mode enabled: you'll be able to select from multiple subtitle options",
				"Dry run mode: no files will be downloaded, only preview what would happen",
			},
		},
		{
			name:        "normal_mode_no_flags",
			cli:         CLI{},
			expectError: false,
			expectMsgs:  []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli := tt.cli
			result, err := cli.validateModeConsistency()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.True(t, result.Success)

				if len(tt.expectMsgs) > 0 {
					for _, msg := range tt.expectMsgs {
						assert.Contains(t, result.Message, msg)
					}
				} else {
					assert.Empty(t, result.Message)
				}
			}
		})
	}
}

func TestIsValidLanguageCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"two_letter_lowercase", "en", true},
		{"two_letter_uppercase", "EN", true},
		{"three_letter_lowercase", "eng", true},
		{"locale_format_lowercase", "pt-br", true},
		{"locale_format_uppercase", "PT-BR", true},
		{"locale_format_mixed", "pt-BR", true},

		{"single_letter", "e", false},
		{"four_letters", "engl", false},
		{"six_letters", "englis", false},
		{"contains_numbers", "en1", false},
		{"contains_special_chars", "en!", false},
		{"underscore_separator", "en_US", false},
		{"locale_first_part_short", "e-BR", false},
		{"locale_first_part_long", "eng-BR", false},
		{"locale_second_part_short", "en-B", false},
		{"locale_second_part_long", "en-BRA", false},
		{"locale_missing_separator", "enBR", false},
		{"empty_string", "", false},
		{"spaces", "  ", false},
		{"locale_with_numbers", "en-B1", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := isValidLanguageCode(tt.code)
			assert.Equal(t, tt.expected, result, "isValidLanguageCode(%q) = %v, want %v", tt.code, result, tt.expected)
		})
	}
}

func TestValidateArguments(t *testing.T) {
	t.Parallel()

	t.Run("all_validations_pass", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "movie.mp4")
		require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))

		configFile := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(configFile, []byte("test: value"), 0644))

		cli := &CLI{
			Path:     tmpFile,
			Language: []string{"en", "pt-BR"},
			Config:   configFile,
		}

		err := cli.validateArguments()
		assert.NoError(t, err)
	})

	t.Run("path_validation_fails", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{
			Path:     "/nonexistent/path",
			Language: []string{"en"},
		}

		err := cli.validateArguments()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path does not exist")
	})

	t.Run("language_validation_fails", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		cli := &CLI{
			Path:     tmpDir,
			Language: []string{"invalid!"},
		}

		err := cli.validateArguments()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid language code")
	})

	t.Run("config_validation_fails", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		cli := &CLI{
			Path:     tmpDir,
			Language: []string{"en"},
			Config:   "/nonexistent/config.yaml",
		}

		err := cli.validateArguments()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config file does not exist")
	})

	t.Run("mode_consistency_fails", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		cli := &CLI{
			Path:     tmpDir,
			Language: []string{"en"},
			Search:   "   ",
		}

		err := cli.validateArguments()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search query cannot be empty")
	})

	t.Run("search_mode_skips_path_validation", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{
			Path:     "/nonexistent/path",
			Language: []string{"en"},
			Search:   "Breaking Bad S01E01",
		}

		err := cli.validateArguments()
		assert.NoError(t, err)
	})
}

func TestPrintValidationResults(t *testing.T) {

	results := []*ValidationResult{
		{Success: true, Message: "Test success message"},
		{Success: false, Message: "Test info message"},
		{Success: true, Warning: "Test warning message"},
		{Success: true, Message: "Success with message", Warning: "And a warning"},
	}

	cli := &CLI{}
	cli.printValidationResults(results)
}

func TestCLIRun(t *testing.T) {
	t.Parallel()

	t.Run("version_flag", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{Version: true}
		err := cli.Run()
		assert.NoError(t, err)
	})

	t.Run("validation_error", func(t *testing.T) {
		t.Parallel()

		cli := &CLI{
			Path:     "/nonexistent/path",
			Language: []string{"en"},
		}

		err := cli.Run()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("successful_validation", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		cli := &CLI{
			Path:     tmpDir,
			Language: []string{"en"},
		}

		err := cli.Run()
		assert.NoError(t, err)
	})
}
