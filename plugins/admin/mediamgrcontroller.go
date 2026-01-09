package admin

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MediaController struct {
	RootPath string
}

func NewMediaController(rootPath string) *MediaController {
	return &MediaController{RootPath: rootPath}
}

// 🟢 List semua file/folder
func (m *MediaController) ListMedia(c *fiber.Ctx) error {
	path := c.Query("path", "./public")
	target := filepath.Join(m.RootPath, path)

	files, err := os.ReadDir(target)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var list []fiber.Map
	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())
		info, _ := os.Stat(filepath.Join(target, f.Name()))

		list = append(list, fiber.Map{
			"name": f.Name(),
			"path": fullPath,
			"type": func() string {
				if f.IsDir() {
					return "folder"
				}
				return "file"
			}(),
			"size":  info.Size(),
			"mtime": info.ModTime(),
			"url":   fmt.Sprintf("/%s", fullPath),
		})
	}

	return c.JSON(list)
}

// 🟡 Upload file
func (m *MediaController) UploadMedia(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid form"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No files received"})
	}
	path := c.FormValue("path", "./public")
	saveDir := filepath.Join(m.RootPath, path)
	fmt.Printf("DEBUG: UploadMedia path='%s' root='%s' saveDir='%s' files=%d\n", path, m.RootPath, saveDir, len(files))
	os.MkdirAll(saveDir, os.ModePerm)

	for _, file := range files {
		savePath := filepath.Join(saveDir, file.Filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"message": "Upload success"})
}

// 🔴 Delete file atau folder
func (m *MediaController) DeleteMedia(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(400).JSON(fiber.Map{"error": "path required"})
	}

	target := filepath.Join(m.RootPath, path)
	if err := os.RemoveAll(target); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted"})
}

// 🟣 Rename file/folder
func (m *MediaController) RenameMedia(c *fiber.Ctx) error {
	var body struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	oldPath := filepath.Join(m.RootPath, body.OldPath)
	newPath := filepath.Join(filepath.Dir(oldPath), body.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Renamed"})
}

// 🟤 Buat folder baru
func (m *MediaController) CreateFolder(c *fiber.Ctx) error {
	var body struct {
		Path   string `json:"path"`
		Folder string `json:"folder"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	newFolder := filepath.Join(m.RootPath, body.Path, body.Folder)
	if err := os.MkdirAll(newFolder, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Folder created"})
}

// ⚪ Preview file text (txt, json, md, dll)
func (m *MediaController) PreviewMedia(c *fiber.Ctx) error {
	path := c.Query("path")
	if path == "" {
		return c.Status(400).JSON(fiber.Map{"error": "path required"})
	}

	target := filepath.Join(m.RootPath, path)
	file, err := os.Open(target)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ext := strings.ToLower(filepath.Ext(path))

	// Deteksi mime type
	mimeType := "text/plain"
	if strings.HasPrefix(ext, ".jpg") || strings.HasPrefix(ext, ".jpeg") {
		mimeType = "image/jpeg"
	} else if strings.HasPrefix(ext, ".png") {
		mimeType = "image/png"
	} else if strings.HasPrefix(ext, ".gif") {
		mimeType = "image/gif"
	} else if strings.HasPrefix(ext, ".webp") {
		mimeType = "image/webp"
	} else if strings.HasPrefix(ext, ".pdf") {
		mimeType = "application/pdf"
	} else if strings.HasPrefix(ext, ".json") {
		mimeType = "application/json"
	}

	var content string
	// Kalau file teks, baca langsung
	if strings.HasPrefix(mimeType, "text/") || strings.Contains(mimeType, "json") {
		content = string(data)
		if len(content) > 2*1024*1024 {
			content = content[:2*1024*1024] + "\n\n...[truncated]"
		}
	} else {
		// Kalau file biner, encode base64
		encoded := base64.StdEncoding.EncodeToString(data)
		content = fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
	}

	return c.JSON(fiber.Map{
		"name":    filepath.Base(path),
		"ext":     ext,
		"content": content,
		"mime":    mimeType,
	})
}

// 💾 Save file content (create or update)
func (m *MediaController) SaveMedia(c *fiber.Ctx) error {
	var body struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Path == "" {
		return c.Status(400).JSON(fiber.Map{"error": "path required"})
	}

	target := filepath.Join(m.RootPath, body.Path)
	
	// Write content to file
	if err := os.WriteFile(target, []byte(body.Content), 0644); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Saved"})
}

// Helpers for backward compatibility if needed, or just use the struct directly
var DefaultMediaController = NewMediaController("./public")

func ListMedia(c *fiber.Ctx) error { return DefaultMediaController.ListMedia(c) }
func UploadMedia(c *fiber.Ctx) error { return DefaultMediaController.UploadMedia(c) }
func DeleteMedia(c *fiber.Ctx) error { return DefaultMediaController.DeleteMedia(c) }
func RenameMedia(c *fiber.Ctx) error { return DefaultMediaController.RenameMedia(c) }
func CreateFolder(c *fiber.Ctx) error { return DefaultMediaController.CreateFolder(c) }
func PreviewMedia(c *fiber.Ctx) error { return DefaultMediaController.PreviewMedia(c) }
