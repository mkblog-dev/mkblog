package main

import (
	"context"
	"log"
	"os"

	mkblog "github.com/mkblog-dev/mkblog"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "d",
				Usage:    "Input directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "o",
				Usage:    "Output directory",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			inputDir := cmd.String("d")
			outputDir := cmd.String("o")

			if inputDir == "" || outputDir == "" {
				log.Println("Both -d (input) and -o (output) must be specified.")
				os.Exit(1)
			}

			err := mkblog.Build(inputDir, outputDir)
			if err != nil {
				log.Fatalf("build failed: %v", err)
			}

			log.Println("done")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
