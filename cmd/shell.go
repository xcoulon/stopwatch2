package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewShellCmd() *cobra.Command {
	var outputFilename string
	shellCmd := &cobra.Command{
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
			fmt.Fprintln(cmd.OutOrStdout(), "**********************************")
			fmt.Fprintln(cmd.OutOrStdout(), "type 'start' when the race begins!")
			fmt.Fprintln(cmd.OutOrStdout(), "**********************************")
		loop:
			for {
				t := prompt.Input("⏱ ", completer)
				switch strings.TrimSpace(t) {
				case "stop", "quit", "exit":
					break loop
				case "start":
					now := time.Now().Local().Format(TimeFormat)
					if _, err = fmt.Fprintf(output, "- %s: start\n", now); err != nil {
						return err
					}
				case "":
					continue
				default: // teams
					now := time.Now().Local().Format(TimeFormat)
					if _, err := strconv.Atoi(t); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "'%s' is not a valid bib number\n", t)
						continue
					}
					if _, err = fmt.Fprintf(output, "- %s: %s\n", now, t); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	shellCmd.Flags().StringVar(&outputFilename, "output", "", "path to write the arrivals (YAML)")
	shellCmd.MarkFlagRequired("output")
	return shellCmd
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
