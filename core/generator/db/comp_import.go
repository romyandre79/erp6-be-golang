package generator

import (
	"erp6-be-golang/core/helpers"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("importdata", func(ctx *WorkflowContext) error {
		return handleImportData(ctx.FiberCtx, ctx.Params, ctx.DB, ctx.FileHeader, ctx.Search)
	})
}

func cSaveFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func handleImportData(c *fiber.Ctx, params []WorkflowDetailResult, db *gorm.DB, fileHeader *multipart.FileHeader, search bool) error {
	var (
		column          string
		insertParameter string
		spInsert        string
		updateParameter string
		spUpdate        string
		method          string
		enable          = true
	)

	// Extract parameters
	for _, p := range params {
		switch p.InputName {
		case "column":
			column = strings.TrimSpace(p.CompValue)
		case "insertparameter":
			insertParameter = strings.TrimSpace(p.CompValue)
		case "spinsert":
			spInsert = strings.TrimSpace(p.CompValue)
		case "updateparameter":
			updateParameter = strings.TrimSpace(p.CompValue)
		case "importtype":
			method = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "spupdate":
			spUpdate = strings.TrimSpace(p.CompValue)
		case "enableimport":
			if p.CompValue == "false" {
				enable = false
			}
		}
	}

	// Save file to "uploads/"
	saveDir := "uploads"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	savePath := filepath.Join(saveDir, filepath.Base(fileHeader.Filename))
	if err := cSaveFile(fileHeader, savePath); err != nil {
		return err
	}

	// Read Excel
	f, err := excelize.OpenFile(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	fmt.Printf("Sheet %s", sheet)
	rows, err := f.GetRows(sheet)
	fmt.Printf("Rows %v", rows)
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return fmt.Errorf("no data found in Excel")
	}
	var sqlStr string

	// Run transaction (rollback on error)
	return db.Transaction(func(tx *gorm.DB) error {
		for i := 1; i < len(rows); i++ {
			colMap := map[string]string{}
			colDefs := strings.Split(column, ",")
			for _, def := range colDefs {
				parts := strings.Split(def, "=")
				if len(parts) != 2 {
					continue
				}
				colName := strings.TrimSpace(parts[0])
				colIndex := strings.TrimSpace(parts[1])
				val, err := f.GetCellValue(sheet, fmt.Sprintf("%s%d", colIndex, i+1))

				if err != nil {
					continue
				}
				colMap[colName] = val
			}

			isEmpty := false
			for _, v := range colMap {
				if strings.TrimSpace(v) == "" {
					isEmpty = true
					break
				}
			}

			var params []string

			if isEmpty {
				params = strings.Split(insertParameter, ",")
				if method == "table" {
					fields, vals := []string{}, []string{}
					for _, p := range params {
						p = strings.TrimSpace(p)
						fields = append(fields, p)
						vals = append(vals, fmt.Sprintf("'%s'", colMap[p]))
					}
					sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", spInsert, strings.Join(fields, ","), strings.Join(vals, ","))
				} else {
					bind := []string{}
					for _, p := range params {
						bind = append(bind, p)
					}
					sqlStr = fmt.Sprintf("CALL %s(%s)", spInsert, strings.Join(bind, ","))
				}
			} else {
				params = strings.Split(updateParameter, ",")
				if method == "table" {
					id := params[0]
					set := []string{}
					for _, p := range params[1:] {
						set = append(set, fmt.Sprintf("%s='%s'", p, colMap[p]))
					}
					sqlStr = fmt.Sprintf("UPDATE %s SET %s WHERE %s='%s'", spUpdate, strings.Join(set, ","), id, colMap[id])
				} else {
					bind := []string{}
					for _, p := range params {
						bind = append(bind, p)
					}
					sqlStr = fmt.Sprintf("CALL %s(%s)", spUpdate, strings.Join(bind, ","))
				}
			}

			if !enable {
				fmt.Println("Preview SQL:", sqlStr)
				continue
			}

			// Execute with GORM
			args := []interface{}{}
			for _, p := range params {
				args = append(args, colMap[p])
			}
			if err := tx.Exec(sqlStr).Error; err != nil {
				log.Printf("Row %d failed: %v", i+1, err)
				return err // rollback
			}
		}
		if !enable {
			helpers.FailResponse(c, 401, "INVALID DATA UPLOADED", sqlStr)
		} else {
			helpers.SuccessResponse(c, "DATA UPLOADED", "Filename "+fileHeader.Filename+" Size "+fmt.Sprintf("%d", fileHeader.Size))
		}
		return nil
	})
}
