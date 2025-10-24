//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var imports []string

	// scan folder plugins
	filepath.WalkDir("plugins", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != "plugins" {
			name := filepath.Base(path)
			imports = append(imports, fmt.Sprintf("_ \"erp6-be-golang/plugins/%s\"", name))
		}
		return nil
	})

	// buat isi file
	content := fmt.Sprintf(`package main

// AUTO-GENERATED FILE: DO NOT EDIT MANUALLY
import (
	%s
)
`, strings.Join(imports, "\n\t"))

	// tulis ke file
	if err := os.WriteFile("plugins_imports.go", []byte(content), 0644); err != nil {
		panic(err)
	}
	fmt.Println("✅ plugins_imports.go berhasil dibuat")
}
