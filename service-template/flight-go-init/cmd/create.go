package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	//go:embed all:template
	templateFS embed.FS

	//go:embed all:template2
	template2FS embed.FS

	//go:embed all:template3
	template3FS embed.FS

	provider             string
	moduleName           string
	serviceName          string
	serviceTitle         string
	basePath             string
	configPropertyStruct string
	dbName               string
	supplierCode         string
	templateName         string
)

type templateInfo struct {
	fs      embed.FS
	rootDir string
	desc    string
}

var templates = map[string]templateInfo{
	"template1": {
		fs:      templateFS,
		rootDir: "template",
		desc:    "Legacy template using hotel-utilities library with custom config structs",
	},
	"template2": {
		fs:      template2FS,
		rootDir: "template2",
		desc:    "Latest template using common-deps (COMMON-LIB-GO, COMMON-MODEL-GO) with dig DI, Echo, MongoDB migrations, BAU analytics, and error mapper",
	},
	"template3": {
		fs:      template3FS,
		rootDir: "template3",
		desc:    "Latest template with rule-supplier registry pattern, processor-based services, smart swagger gen, and streamlined project structure",
	},
}

var createCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "Create a new microservice from template",
	Long: `Create a new microservice from a selected template.

Available templates:
  template1  -  Legacy template using hotel-utilities library with custom config structs
  template2  -  Latest template using common-deps (COMMON-LIB-GO, COMMON-MODEL-GO) with
                 dig DI, Echo, MongoDB migrations, BAU analytics, and error mapper
  template3  -  Latest template with rule-supplier registry pattern, processor-based services,
                 smart swagger gen, and streamlined project structure`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]

		tmpl, ok := templates[templateName]
		if !ok {
			return fmt.Errorf("unknown template %q, available: template1, template2, template3", templateName)
		}

		fmt.Printf("Generating project: %s (using %s)\n", project, templateName)
		fmt.Printf("  Template: %s\n", tmpl.desc)

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

		err := copyTemplate(project, replacements, tmpl.fs, tmpl.rootDir)
		if err != nil {
			return fmt.Errorf("failed to copy template: %w", err)
		}

		reqPath := filepath.Join(project, "requirements.txt")
		if _, err := os.Stat(reqPath); err == nil {
			err = generateGoModFromRequirements(reqPath, project)
			if err != nil {
				return fmt.Errorf("failed to generate go.mod: %w", err)
			}
			fmt.Println("go.mod created from requirements.txt")
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

	createCmd.Flags().StringVar(&templateName, "template", "template1", `Template to use for project generation:
  template1  - Legacy template using hotel-utilities library with custom config structs
  template2  - Latest template using common-deps (COMMON-LIB-GO) with dig DI, Echo, MongoDB migrations, BAU analytics
  template3  - Latest template with rule-supplier registry, processor services, smart swagger gen`)
	createCmd.Flags().StringVar(&provider, "provider", "", "Provider name (e.g. trigana-go)")
	createCmd.Flags().StringVar(&moduleName, "module-name", "", "Go module name (e.g. trigana-go-be)")
	createCmd.Flags().StringVar(&serviceName, "service-name", "", "Service/binary name (e.g. TIX-FLIGHT-TRIGANA-INTEGRATOR-GO-BE)")
	createCmd.Flags().StringVar(&serviceTitle, "service-title", "", "Service title for Swagger (e.g. TRIGANA GO INTEGRATOR BE)")
	createCmd.Flags().StringVar(&basePath, "base-path", "", "Base API path (e.g. tix-flight-trigana-go-integrator)")
	createCmd.Flags().StringVar(&configPropertyStruct, "config-property-struct", "", "Config struct name (e.g. TriganaConfigProperties)")
	createCmd.Flags().StringVar(&dbName, "db-name", "", "Database name (e.g. flight_trigana_integrator)")
	createCmd.Flags().StringVar(&supplierCode, "supplier-code", "", "Supplier code (e.g. tiketcomTrigana)")

	createCmd.MarkFlagRequired("provider")
	createCmd.MarkFlagRequired("module-name")
	createCmd.MarkFlagRequired("service-name")
	createCmd.MarkFlagRequired("service-title")
	createCmd.MarkFlagRequired("base-path")
	createCmd.MarkFlagRequired("config-property-struct")
	createCmd.MarkFlagRequired("db-name")
	createCmd.MarkFlagRequired("supplier-code")
}

func copyTemplate(dest string, vars map[string]string, fsys embed.FS, rootDir string) error {
	return fs.WalkDir(fsys, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(path, rootDir)
		rel = applyReplacements(rel, vars)
		rel = strings.TrimSuffix(rel, ".tmpl")

		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}

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
