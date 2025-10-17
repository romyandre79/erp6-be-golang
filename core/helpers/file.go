package helpers

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TemplateData struct {
	Name        string
	PackageName string
	GeneratedAt string
}

func ensureDir(path string, mode os.FileMode) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists and is not a directory: %s", path)
		}
		return nil // dir exists
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, mode)
}

func writeFileIfNotExist(path string, data []byte, mode os.FileMode, force bool) error {
	_, err := os.Stat(path)
	if err == nil {
		if !force {
			return fmt.Errorf("file already exists: %s", path)
		}
		// remove old file before writing so permission reset works predictably
		if err = os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return err
	}
	return os.Chmod(path, mode)
}

func renderTemplate(tmplStr string, d TemplateData) ([]byte, error) {
	tmpl, err := template.New("t").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, d); err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}

func isWritable(path string) bool {
	// quick check: try opening for append (without trunc) — will fail if not writable
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

func checkFilePermission(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}
