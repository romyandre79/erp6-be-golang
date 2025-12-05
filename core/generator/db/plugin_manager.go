package generator

import (
	"archive/zip"
	"encoding/json"
	"erp6-be-golang/models"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	if err := tx.Where("componentname = ?", manifest.Componentname).First(&existingComp).Error; err == nil {
		tx.Rollback()
		return c.Status(409).JSON(fiber.Map{
			"error":   "plugin_exists",
			"message": fmt.Sprintf("Plugin '%s' already exists", manifest.Componentname),
		})
	}

	// Create Component
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

	// Create Details
	for _, detail := range manifest.Details {
		detail.Componentid = newComp.Componentid
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

		// Special handling for .exe files (Windows) -> register as external plugin
		if strings.HasSuffix(f.Name, ".exe") {
			// Save to plugins/bin/
			binDir := "./plugins/bin"
			if err := os.MkdirAll(binDir, 0755); err != nil {
				continue
			}

			binPath := filepath.Join(binDir, f.Name)
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
			key := strings.ToLower(manifest.Componentclass)
			if key == "" {
				key = strings.ToLower(manifest.Componentname)
			}

			// Use absolute path for safety
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

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "Plugin registered successfully",
		"componentid": newComp.Componentid,
	})
}
