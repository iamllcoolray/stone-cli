package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/iamllcoolray/stone-cli/internal/template"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new LITIengine project",
	Long:  "Scaffolds a new LITIengine project with Gradle build files, a Main.java entry point, and the standard resource directory structure.",
	RunE:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Creating a new LITIengine project...")
	fmt.Println()

	projectName := prompt(reader, "Project name", "")
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	groupID := prompt(reader, "Group ID", "com.example")
	if groupID == "" {
		return fmt.Errorf("group ID cannot be empty")
	}

	defaultArtifact := strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))
	artifactID := prompt(reader, "Artifact ID", defaultArtifact)
	if artifactID == "" {
		return fmt.Errorf("artifact ID cannot be empty")
	}

	version := prompt(reader, "Version", "1.0.0")
	if version == "" {
		version = "1.0.0"
	}

	litiVersion := prompt(reader, "LITIengine version", "0.11.1")
	if litiVersion == "" {
		litiVersion = "0.11.1"
	}

	data := template.NewProjectData(projectName, groupID, artifactID, version, litiVersion)

	outputDir := artifactID

	if _, err := os.Stat(outputDir); err == nil {
		fmt.Printf("Directory '%s' already exists. Overwrite? [y/N] ", outputDir)
		answer, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	fmt.Println()
	fmt.Printf("Generating project in ./%s/\n", outputDir)

	if err := template.Generate(outputDir, data); err != nil {
		return fmt.Errorf("generating project: %w", err)
	}

	fmt.Println()
	fmt.Println("Project created successfully.")
	fmt.Println()
	fmt.Printf("  Package:      %s\n", data.PackageName)
	fmt.Printf("  LITIengine:   %s\n", data.LitiVersion)
	fmt.Printf("  Entry point:  %s/src/main/java/%s/Main.java\n", outputDir, data.PackagePath)
	fmt.Println()
	fmt.Printf("To get started:\n")
	fmt.Printf("  cd %s\n", outputDir)
	fmt.Printf("  gradle run\n")

	return nil
}
