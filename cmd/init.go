package cmd

import (
	"bufio"

	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up stone and install utiLITI",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	config.Load()
	return nil
}

func prompt(r *bufio.Reader, question, defaultVal string) string {
	return ""
}
