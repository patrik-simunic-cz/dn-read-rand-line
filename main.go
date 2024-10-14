package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"

	"readline/generator"
	"readline/reader"
)

func main() {
	app := &cli.App{
		Name:  "basicread",
		Usage: "Read (basic)",
		Action: func(context *cli.Context) (err error) {
			var line int

			args := context.Args()

			path := args.Get(0)
			if path == "" {
				return fmt.Errorf("Missing path to the file to read")
			}

			lineRaw := args.Get(1)
			if line, err = strconv.Atoi(lineRaw); err != nil {
				return fmt.Errorf("Invalid line (\"%s\"): expected an integer", lineRaw)
			}

			if line < 0 {
				return fmt.Errorf("Line index to print must be greater than or equal to 0")
			}

			return reader.ReadLine("./index.idx", path, line+1, false)
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
			Name:  "read",
			Usage: "Read line",
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
				&cli.IntFlag{
					Name:     "line",
					Usage:    "Line to print",
					Required: true,
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

				line := context.Int("line")
				if line < 1 {
					return fmt.Errorf("Line to print must be greater than 0")
				}

				enableVerbose := context.Bool("verbose")

				return reader.ReadLine(indexPath, path, line, enableVerbose)
			}},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
