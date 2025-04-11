package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

func loadTemplates() error {
	templates = template.New("")

	entries, err := os.ReadDir("templates")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}
		content, err := os.ReadFile(filepath.Join("templates", entry.Name()))
		if err != nil {
			return err
		}
		_, err = templates.Parse(string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
