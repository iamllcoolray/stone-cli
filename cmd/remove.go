package cmd

import (
	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	removeAll bool
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove utiLITI and/or stone from the system",
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&removeAll, "all", false, "Remove utiLITI and stone from the system")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	config.Load()
	return nil
}
