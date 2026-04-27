package cmd

import "github.com/spf13/cobra"

var (
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update utiLITI to the latest version",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false, "Force update even if already on latest version")
	rootCmd.AddCommand(updateCmd)
}

// runUpdate loads config, checks for a new version,
// fetches the download URL and upload ID for current platform,
// runs the updater, and saves the new version to config
func runUpdate(cmd *cobra.Command, args []string) error {
	return nil
}
