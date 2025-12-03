package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	//go:embed template/**
	templateFS embed.FS

	provider             string
	moduleName           string
	serviceName          string
	serviceTitle         string
	basePath             string
	configPropertyStruct string
	dbName               string
	supplierCode         string
)

var createCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "Create a new microservice from template",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		fmt.Println("Generating project:", project)

		// placeholder replacements
		replacements := map[string]string{
			"{{PROVIDER}}":               provider,
			"{{MODULE_NAME}}":            moduleName,
			"{{SERVICE_NAME}}":           serviceName,
			"{{SERVICE_TITLE}}":          serviceTitle,
			"{{BASE_PATH}}":              basePath,
			"{{CONFIG_PROPERTY_STRUCT}}": configPropertyStruct,
			"{{DB_NAME}}":                dbName,
			"{{SUPPLIER_CODE}}":          supplierCode,
		}

		// Step 1: copy template files
		err := copyTemplate(project, replacements)
		if err != nil {
			return fmt.Errorf("failed to copy template: %w", err)
		}

		// Step 2: generate go.mod from requirements.txt if it exists
		reqPath := filepath.Join(project, "requirements.txt")
		if _, err := os.Stat(reqPath); err == nil {
			err = generateGoModFromRequirements(reqPath, project)
			if err != nil {
				return fmt.Errorf("failed to generate go.mod: %w", err)
			}
			fmt.Println("go.mod created from requirements.txt")
		} else {
			fmt.Println("No requirements.txt found, skipping go.mod generation")
		}
		fmt.Println("Project generated successfully!")
		return nil
	},
}

func generateGoModFromRequirements(reqFile, projectDir string) error {
	data, err := os.ReadFile(reqFile)
	if err != nil {
		return fmt.Errorf("cannot read requirements.txt: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var moduleLine, goLine string
	var deps []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			moduleLine = line
		} else if strings.HasPrefix(line, "go ") {
			goLine = line
		} else {
			deps = append(deps, line)
		}
	}

	if moduleLine == "" {
		return fmt.Errorf("requirements.txt must contain a 'module' line")
	}
	if goLine == "" {
		goLine = "go 1.21" // default
	}

	var sb strings.Builder
	sb.WriteString(moduleLine + "\n")
	sb.WriteString(goLine + "\n\n")
	if len(deps) > 0 {
		//sb.WriteString("require (\n")
		for _, dep := range deps {
			sb.WriteString("\t" + dep + "\n")
		}
		//sb.WriteString(")\n")
	}

	goModPath := filepath.Join(projectDir, "go.mod")
	return os.WriteFile(goModPath, []byte(sb.String()), 0644)
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&provider, "provider", "", "Provider name")
	createCmd.Flags().StringVar(&moduleName, "module-name", "", "Module name")
	createCmd.Flags().StringVar(&serviceName, "service-name", "", "Service name")
	createCmd.Flags().StringVar(&serviceTitle, "service-title", "", "Service Title")
	createCmd.Flags().StringVar(&basePath, "base-path", "", "Base API Path")
	createCmd.Flags().StringVar(&configPropertyStruct, "config-property-struct", "", "Config struct name")
	createCmd.Flags().StringVar(&dbName, "db-name", "", "Database name")
	createCmd.Flags().StringVar(&supplierCode, "supplier-code", "", "Supplier code")

	// required flags
	createCmd.MarkFlagRequired("provider")
	createCmd.MarkFlagRequired("module-name")
	createCmd.MarkFlagRequired("service-name")
	createCmd.MarkFlagRequired("service-title")
	createCmd.MarkFlagRequired("base-path")
	createCmd.MarkFlagRequired("config-property-struct")
	createCmd.MarkFlagRequired("db-name")
	createCmd.MarkFlagRequired("supplier-code")
}

func copyTemplate(dest string, vars map[string]string) error {
	return fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path + apply replacements
		rel := strings.TrimPrefix(path, "template")
		rel = applyReplacements(rel, vars)

		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		// Read file content
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace variables inside file
		content := applyReplacements(string(data), vars)

		return os.WriteFile(target, []byte(content), 0644)
	})
}

func applyReplacements(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
