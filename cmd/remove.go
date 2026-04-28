package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	config "github.com/iamllcoolray/stone-cli/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	removeStone bool
	removeAll   bool
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove utiLITI and/or stone from the system",
	Long:  "Without flags removes the utiLITI installation.",
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&removeStone, "stone", false, "Remove stone itself from the system")
	removeCmd.Flags().BoolVar(&removeAll, "all", false, "Remove utiLITI and stone from the system")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	cfg, _ := config.Load()

	if removeAll {
		return removeEverything(reader, cfg)
	}

	if removeStone {
		return removeStoneOnly(reader, cfg)
	}

	return removeUtiLITI(reader, cfg)
}

func removeUtiLITI(reader *bufio.Reader, cfg *config.Config) error {
	if cfg.InstallPath == "" {
		return fmt.Errorf("install_path is not set — run: stone init")
	}

	fmt.Printf("This will remove utiLITI from: %s\n", cfg.InstallPath)
	if !confirm(reader) {
		fmt.Println("Aborted.")
		return nil
	}

	if err := os.RemoveAll(cfg.InstallPath); err != nil {
		return fmt.Errorf("removing utiLITI: %w", err)
	}

	// clear last_version from config
	cfg.LastVersion = ""
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("updating config: %w", err)
	}

	fmt.Println("utiLITI removed.")
	return nil
}

func removeStoneOnly(reader *bufio.Reader, cfg *config.Config) error {
	stonePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding stone binary: %w", err)
	}

	configDir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("finding config dir: %w", err)
	}

	fmt.Printf("This will remove:\n")
	fmt.Printf("  stone binary: %s\n", stonePath)
	fmt.Printf("  config dir:   %s\n", configDir)

	if !confirm(reader) {
		fmt.Println("Aborted.")
		return nil
	}

	if err := os.Remove(stonePath); err != nil {
		return fmt.Errorf("removing stone binary: %w", err)
	}

	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("removing config: %w", err)
	}

	fmt.Println("stone removed.")
	return nil
}

func removeEverything(reader *bufio.Reader, cfg *config.Config) error {
	stonePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding stone binary: %w", err)
	}

	configDir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("finding config dir: %w", err)
	}

	fmt.Println("This will remove:")
	if cfg.InstallPath != "" {
		fmt.Printf("  utiLITI:      %s\n", cfg.InstallPath)
	}
	fmt.Printf("  stone binary: %s\n", stonePath)
	fmt.Printf("  config dir:   %s\n", configDir)

	if !confirm(reader) {
		fmt.Println("Aborted.")
		return nil
	}

	if cfg.InstallPath != "" {
		if err := os.RemoveAll(cfg.InstallPath); err != nil {
			return fmt.Errorf("removing utiLITI: %w", err)
		}
	}

	if err := os.Remove(stonePath); err != nil {
		return fmt.Errorf("removing stone binary: %w", err)
	}

	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("removing config: %w", err)
	}

	fmt.Println("Everything removed.")
	return nil
}

func confirm(r *bufio.Reader) bool {
	fmt.Print("Are you sure? [y/N] ")
	answer, _ := r.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(answer)) == "y"
}
