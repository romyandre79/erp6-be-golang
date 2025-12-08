package generator

import (
	"archive/zip"
	"encoding/json"
	"erp6-be-golang/models"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PluginManifest defines the structure of plugin.json
type PluginManifest struct {
	Componentname       string                   `json:"componentname"`
	Componenttitle      string                   `json:"componenttitle"`
	Componentcategoryid int                      `json:"componentcategoryid"`
	Componentclass      string                   `json:"componentclass"`
	Version             string                   `json:"version"`
	Createdby           string                   `json:"createdby"`
	Input               int8                     `json:"input"`
	Output              int8                     `json:"output"`
	Details             []models.Componentdetail `json:"details"`
}

// HandlePluginUpload handles the upload and registration of a plugin ZIP
func HandlePluginUpload(c *fiber.Ctx, db *gorm.DB) error {
	// 1. Get the file
	file, err := c.FormFile("plugin")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "missing_file",
			"message": "Plugin ZIP file is required",
		})
	}

	// 2. Save to temp
	tempDir := "./tmp/plugins"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	tempPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer os.Remove(tempPath) // Clean up zip

	// 3. Open ZIP
	r, err := zip.OpenReader(tempPath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "invalid_zip",
			"message": "Could not open ZIP file",
		})
	}
	defer r.Close()

	// 4. Find and parse plugin.json
	var manifest PluginManifest
	foundManifest := false

	for _, f := range r.File {
		if f.Name == "plugin.json" {
			rc, err := f.Open()
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			if err := json.Unmarshal(content, &manifest); err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error":   "invalid_manifest",
					"message": "Invalid plugin.json format",
				})
			}
			foundManifest = true
			break
		}
	}

	if !foundManifest {
		return c.Status(400).JSON(fiber.Map{
			"error":   "missing_manifest",
			"message": "plugin.json not found in ZIP",
		})
	}

	// 5. Register in Database
	tx := db.Begin()

	// Check if component exists
	var existingComp models.Component
	var componentID int

	err = tx.Where("componentname = ?", manifest.Componentname).First(&existingComp).Error
	if err == nil {
		// UPDATE
		componentID = existingComp.Componentid

		existingComp.Componenttitle = manifest.Componenttitle
		existingComp.Componentcategoryid = manifest.Componentcategoryid
		existingComp.Componentclass = manifest.Componentclass
		existingComp.Version = manifest.Version
		existingComp.Createdby = manifest.Createdby
		existingComp.Input = manifest.Input
		existingComp.Output = manifest.Output

		if err := tx.Save(&existingComp).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Delete old details
		if err := tx.Where("componentid = ?", componentID).Delete(&models.Componentdetail{}).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

	} else if err == gorm.ErrRecordNotFound {
		// CREATE
		newComp := models.Component{
			Componentname:       manifest.Componentname,
			Componenttitle:      manifest.Componenttitle,
			Componentcategoryid: manifest.Componentcategoryid,
			Componentclass:      manifest.Componentclass,
			Version:             manifest.Version,
			Createdby:           manifest.Createdby,
			Input:               manifest.Input,
			Output:              manifest.Output,
		}

		if err := tx.Create(&newComp).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		componentID = newComp.Componentid

	} else {
		// Error
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create Details (for both new and update)
	for _, detail := range manifest.Details {
		detail.Componentid = componentID
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// 6. Extract Assets (Optional - extract everything else to public/plugins/{name})
	assetsDir := filepath.Join("./public/plugins", strings.ToLower(manifest.Componentname))
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for _, f := range r.File {
		if f.Name == "plugin.json" || f.FileInfo().IsDir() {
			continue
		}

		// Prevent Zip Slip
		fpath := filepath.Join(assetsDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(assetsDir)+string(os.PathSeparator)) {
			continue
		}

		// Special handling for executables (Windows .exe, or Linux/Mac binaries)
		// Logic: Detect current OS/Arch and only save matching binary
		lowerName := strings.ToLower(f.Name)
		isBinary := false

		// Current System - using strings package instead of importing runtime for now to avoid large diffs if possible,
		// but best to just rely on heuristics we know from build scripts.
		// Patterns: _windows_amd64.exe, _linux_amd64, _darwin_amd64, _darwin_arm64
		// If simple name (chat.exe), assume it matches if OS matches extension.

		// Check if it looks like a binary we built
		if strings.HasSuffix(lowerName, ".exe") || strings.Contains(lowerName, "_linux_") || strings.Contains(lowerName, "_darwin_") {
			isBinary = true
		}

		if isBinary {
			// runtime.GOOS/GOARCH are needed for strict matching.
			// Let's import runtime. For now we assume the user's intent:
			// "move only the file same as OS and remove other"

			// We can't import runtime easily in middle of function without updating imports.
			// Assuming we will update imports in next step or use simple heuristic if we can't.
			// Ideally we use: runtime.GOOS, runtime.GOARCH

			// Let's implement strict check assuming runtime is available (will add import)
			target := "_" + runtime.GOOS + "_" + runtime.GOARCH

			// Exception: if file is just "chat.exe" and we are on windows
			isExactMatch := false
			if runtime.GOOS == "windows" && strings.HasSuffix(lowerName, ".exe") && !strings.Contains(lowerName, "_windows_") && !strings.Contains(lowerName, "_linux_") && !strings.Contains(lowerName, "_darwin_") {
				isExactMatch = true
			}

			if !strings.Contains(lowerName, target) && !isExactMatch {
				// Skip mismatching binaries
				continue
			}

			// Save to plugins/bin/
			binDir := "./plugins/bin"
			if err := os.MkdirAll(binDir, 0755); err != nil {
				continue
			}

			saveName := f.Name
			key := strings.ToLower(manifest.Componentname)
			if key != "" {
				saveName = key + filepath.Ext(f.Name)
			}

			binPath := filepath.Join(binDir, saveName)
			// Ensure we write with execute permissions
			outFile, err := os.OpenFile(binPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				continue
			}

			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				continue
			}
			io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			// Register in memory
			absPath, _ := filepath.Abs(binPath)
			GlobalRegistry.RegisterExternal(key, absPath)
			fmt.Printf("Registered external plugin: %s -> %s\n", key, absPath)
			continue
		}

		// Special handling for .go files -> register as dynamic plugin (Yaegi)
		if strings.HasSuffix(f.Name, ".go") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, _ := io.ReadAll(rc)
			rc.Close()

			// Register in memory immediately
			// Use the component class/name from manifest as the key
			key := strings.ToLower(manifest.Componentclass)
			if key == "" {
				key = strings.ToLower(manifest.Componentname)
			}

			GlobalRegistry.RegisterDynamic(key, string(content))
			fmt.Printf("Registered dynamic plugin: %s\n", key)

			// Also save to disk for persistence (optional, but good practice)
			goFilePath := filepath.Join("./core/generator/db", f.Name)
			if err := os.WriteFile(goFilePath, content, 0644); err != nil {
				fmt.Printf("Failed to write go file: %v\n", err)
			}
			continue
		}

		// Create dir if needed
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			continue
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}

	tx.Commit()

	// 7. Trigger Reload of Plugins (Hot Reload)
	go LoadPlugins(db)

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "Plugin registered successfully",
		"componentid": componentID,
	})
}

// LoadPlugins matches executable files in plugins/bin with registered components
func LoadPlugins(db *gorm.DB) {
	fmt.Println("Scanning for external plugins...")

	// 1. Get all components/classes
	var components []models.Component
	db.Find(&components)

	// Map class -> component
	classMap := make(map[string]bool)
	for _, c := range components {
		if c.Componentclass != "" {
			classMap[strings.ToLower(c.Componentclass)] = true
		}
	}

	// 2. Scan bin directory
	binDir := "./plugins/bin"
	files, err := os.ReadDir(binDir)
	if err != nil {
		fmt.Printf("Error reading plugin bin dir: %v\n", err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := strings.ToLower(f.Name())
		// Only executables
		if !strings.HasSuffix(name, ".exe") && !strings.HasSuffix(name, "") {
			// On linux, no extension, but let's assume valid binaries for now
			// To be safe, maybe just skip .json or .txt?
			// For now, let's accept all files that don't have known non-binary extensions if not windows
		}

		fullPath, _ := filepath.Abs(filepath.Join(binDir, f.Name()))

		// Registration Strategy:
		// 1. Exact match (minus extension)
		baseName := strings.TrimSuffix(name, filepath.Ext(name))

		// Register as baseName
		GlobalRegistry.RegisterExternal(baseName, fullPath)
		fmt.Printf("Loaded plugin: %s -> %s\n", baseName, fullPath)

		// 2. Heuristic: Split by _ (e.g. mysql_windows_amd64 -> mysql)
		parts := strings.Split(baseName, "_")
		if len(parts) > 1 {
			shortName := parts[0]
			// Always register the short alias
			GlobalRegistry.RegisterExternal(shortName, fullPath)
			fmt.Printf("Loaded plugin alias: %s -> %s\n", shortName, fullPath)
		}
	}
}
