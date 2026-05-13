package template

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sylphy/git-switch/core/config"
	"gopkg.in/yaml.v3"
)

type Template struct {
	Template TemplateMeta `yaml:"template"`
	Config   TemplateCfg  `yaml:"config"`
	Init     TemplateInit `yaml:"init,omitempty"`
}

type TemplateMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type TemplateCfg struct {
	Profile string            `yaml:"profile"`
	Git     map[string]string `yaml:"git,omitempty"`
}

type TemplateInit struct {
	Gitignore string `yaml:"gitignore,omitempty"`
}

func templatesDir() (string, error) {
	dir, err := config.DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "templates"), nil
}

func ListTemplates() ([]string, error) {
	tdir, err := templatesDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(tdir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read templates directory %s: %w", tdir, err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

func CreateTemplate(name, profileName, description, gitignore string) error {
	if name == "" {
		return errors.New("template name is required")
	}
	tdir, err := templatesDir()
	if err != nil {
		return err
	}
	tmplDir := filepath.Join(tdir, name)
	if err := os.MkdirAll(tmplDir, 0o700); err != nil {
		return fmt.Errorf("create template directory: %w", err)
	}

	tmpl := Template{
		Template: TemplateMeta{Name: name, Description: description},
		Config:   TemplateCfg{Profile: profileName},
	}
	if gitignore != "" {
		tmpl.Init = TemplateInit{Gitignore: strings.TrimSpace(gitignore) + "\n"}
	}

	data, err := yaml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("marshal template: %w", err)
	}
	path := filepath.Join(tmplDir, "template.yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write template: %w", err)
	}
	return nil
}

func LoadTemplate(name string) (Template, error) {
	tdir, err := templatesDir()
	if err != nil {
		return Template{}, err
	}
	path := filepath.Join(tdir, name, "template.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Template{}, fmt.Errorf("template %q not found", name)
		}
		return Template{}, fmt.Errorf("read template %s: %w", path, err)
	}
	var tmpl Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return Template{}, fmt.Errorf("parse template %s: %w", path, err)
	}
	return tmpl, nil
}

func ApplyTemplate(name, targetDir string) error {
	tmpl, err := LoadTemplate(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	// Write git config
	if len(tmpl.Config.Git) > 0 {
		var sb strings.Builder
		for k, v := range tmpl.Config.Git {
			section, key, found := strings.Cut(k, ".")
			if !found {
				key = section
				section = ""
			}
			if section != "" {
				sb.WriteString(fmt.Sprintf("[%s]\n", section))
			}
			sb.WriteString(fmt.Sprintf("\t%s = %s\n", key, v))
		}
		configPath := filepath.Join(targetDir, ".gitconfig")
		if err := os.WriteFile(configPath, []byte(sb.String()), 0o600); err != nil {
			return fmt.Errorf("write .gitconfig: %w", err)
		}
	}

	// Write .gitignore
	if tmpl.Init.Gitignore != "" {
		gitignorePath := filepath.Join(targetDir, ".gitignore")
		if err := os.WriteFile(gitignorePath, []byte(tmpl.Init.Gitignore), 0o644); err != nil {
			return fmt.Errorf("write .gitignore: %w", err)
		}
	}

	return nil
}
