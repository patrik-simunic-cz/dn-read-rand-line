package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"readline/generator"
	"readline/reader"
)

func main() {
	app := &cli.App{
		Name:  "help",
		Usage: "Help",
		Action: func(context *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{{
			Name:  "generate",
			Usage: "Generate test data",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Value:    "",
					Usage:    "Path to the generated file",
					Required: true,
				},
				&cli.IntFlag{
					Name:  "lines",
					Value: 1000,
					Usage: "Generated lines count",
				},
				&cli.IntFlag{
					Name:  "wordsPerLine",
					Value: 30,
					Usage: "Max words per generated line",
				},
			},
			Action: func(context *cli.Context) error {
				path := context.String("path")
				if path == "" {
					return fmt.Errorf("Argument --path is required")
				}

				lines := context.Int("lines")
				if lines < 1 {
					return fmt.Errorf("Argument --lines is required and must be grather than 0")
				}

				wordsPerLine := context.Int("wordsPerLine")
				if lines < 1 {
					return fmt.Errorf("Argument --wordsPerLine is required and must be grather than 0")
				}

				return generator.Generate(path, lines, wordsPerLine)
			},
		}, {
			Name:  "rand",
			Usage: "Read random line",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Value:    "",
					Usage:    "Path to the file to read",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "indexPath",
					Value: "./index.idx",
					Usage: "Path to the index file",
				},
				&cli.BoolFlag{
					Name:  "verbose",
					Usage: "Print stats",
				},
			},
			Action: func(context *cli.Context) error {
				indexPath := context.String("indexPath")
				if indexPath == "" {
					return fmt.Errorf("Missing index file path")
				}

				path := context.String("path")
				if path == "" {
					return fmt.Errorf("Missing path to the file to read")
				}

				enableVerbose := context.Bool("verbose")

				return reader.ReadRandomLine(indexPath, path, enableVerbose)
			}},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
