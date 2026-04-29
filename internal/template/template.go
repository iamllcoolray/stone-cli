package template

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFiles embed.FS

type ProjectData struct {
	ProjectName string
	GroupID     string
	ArtifactID  string
	Version     string
	LitiVersion string
	PackageName string
	PackagePath string
}

func NewProjectData(projectName, groupID, artifactID, version, litiVersion string) ProjectData {
	sanitized := strings.ReplaceAll(artifactID, "-", "")
	sanitized = strings.ReplaceAll(sanitized, "_", "")
	sanitized = strings.ToLower(sanitized)

	packageName := groupID + "." + sanitized
	packagePath := strings.ReplaceAll(packageName, ".", "/")

	return ProjectData{
		ProjectName: projectName,
		GroupID:     groupID,
		ArtifactID:  artifactID,
		Version:     version,
		LitiVersion: litiVersion,
		PackageName: packageName,
		PackagePath: packagePath,
	}
}

func Generate(outputDir string, data ProjectData) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating project dir: %w", err)
	}

	javaSrcDir := filepath.Join(outputDir, "src", "main", "java", filepath.FromSlash(data.PackagePath))
	if err := os.MkdirAll(javaSrcDir, 0755); err != nil {
		return fmt.Errorf("creating java src dir: %w", err)
	}

	resourceDirs := []string{
		"src/main/resources/audio",
		"src/main/resources/localization",
		"src/main/resources/maps",
		"src/main/resources/misc",
		"src/main/resources/sprites",
	}
	for _, dir := range resourceDirs {
		if err := os.MkdirAll(filepath.Join(outputDir, dir), 0755); err != nil {
			return fmt.Errorf("creating resource dir %s: %w", dir, err)
		}
	}

	type templateFile struct {
		src  string
		dest string
	}

	files := []templateFile{
		{"templates/build.gradle.ftl", "build.gradle"},
		{"templates/settings.gradle.ftl", "settings.gradle"},
		{"templates/game.litidata.ftl", "game.litidata"},
		{"templates/Main.java.ftl", filepath.Join("src", "main", "java", filepath.FromSlash(data.PackagePath), "Main.java")},
	}

	for _, f := range files {
		if err := renderTemplate(f.src, filepath.Join(outputDir, f.dest), data); err != nil {
			return fmt.Errorf("rendering %s: %w", f.src, err)
		}
	}

	return nil
}

func renderTemplate(src, dest string, data ProjectData) error {
	content, err := templateFiles.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading template %s: %w", src, err)
	}

	tmpl, err := template.New(src).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parsing template %s: %w", src, err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("creating dir for %s: %w", dest, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", dest, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing template %s: %w", src, err)
	}

	return nil
}
