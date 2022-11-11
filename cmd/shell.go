package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newShellCmd() *cobra.Command {
	var outputFilename string
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "interactive shell for races",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if !force {
				return checkOutputFile(outputFilename)
			}
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("bye! 👋")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: use MultiWriter to write in a backup file
			output, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer output.Close()
		loop:
			for {
				t := prompt.Input("⏱ ", completer)
				switch t {
				case "stop", "quit", "exit":
					break loop
				case "start":
					now := time.Now().Local().Format(TimeFormat)
					if _, err = fmt.Fprintf(output, "start: %s\nteams:\n", now); err != nil {
						return err
					}
				default: // "start" and any team number
					now := time.Now().Local().Format(TimeFormat)
					if _, err = fmt.Fprintf(output, "- bibNumber: %s\n  scratch: %s\n", t, now); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&outputFilename, "output", "", "path to write the arrivals (YAML)")
	cmd.MarkFlagRequired("output")
	return cmd
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "start", Description: "start the timer"},
		{Text: "stop", Description: "exit"},
		{Text: "quit", Description: "exit"},
		{Text: "exit", Description: "exit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
