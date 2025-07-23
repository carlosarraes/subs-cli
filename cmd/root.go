package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var Version = "dev"

type CLI struct {
	Path        string   `arg:"" default:"." help:"Path to media file or directory"`
	Language    []string `short:"l" default:"en" help:"Subtitle language(s)"`
	Interactive bool     `short:"i" help:"Interactive mode"`
	Config      string   `short:"c" type:"path" help:"Config file path"`
	DryRun      bool     `help:"Preview without downloading"`
	Search      string   `short:"s" help:"Manual search query"`
	Version     bool     `short:"v" help:"Show version"`
}

func (c *CLI) Run() error {
	if c.Version {
		fmt.Printf("subs-cli version %s\n", Version)
		return nil
	}

	// TODO: Implement subtitle search and download logic
	fmt.Printf("Searching for subtitles in path: %s\n", c.Path)
	fmt.Printf("Languages: %v\n", c.Language)
	
	return nil
}

func Execute() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("subs"),
		kong.Description("A powerful CLI tool for automatically finding and downloading subtitles for your media files"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err := cli.Run()
	ctx.FatalIfErrorf(err)
}