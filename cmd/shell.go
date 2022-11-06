package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "interactive shell for races",
	Args:  cobra.ExactArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkOutputFile(output)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("bye! 👋")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: use MultiWriter to write in a backup file
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
	loop:
		for {
			t := prompt.Input("⏱ ", completer)
			switch t {
			case "stop", "quit", "exit":
				break loop
			default:
				now := time.Now().Local().Format("2006-01-02:15:04:05")
				if _, err = fmt.Fprintf(f, "%s:\t%s\n", t, now); err != nil {
					return err
				}
			}
		}
		return nil
	},
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

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVar(&output, "output", "", "path to write the arrivals")
	shellCmd.MarkFlagRequired("output")
}
