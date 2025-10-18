package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"gorm.io/gorm"
)

// ==========================
// Struct metadata
// ==========================
type Column struct {
	Name       string
	DataType   string
	IsPrimary  bool
	IsNullable bool
}

type Relation struct {
	TableName    string
	ForeignKey   string
	TargetTable  string
	TargetColumn string
	RelationType string // one-to-many, many-to-one, many-to-many
}

// ==========================
// MAIN ENTRY POINT
// ==========================
func GenerateStructWithMeta(db *gorm.DB, table string) error {
	cols, err := GetColumns(db, table)
	if err != nil {
		return err
	}

	rels, err := GetRelations(db, table)
	if err != nil {
		return err
	}

	fmt.Println("Generating struct for table:", table)
	return GenerateStructTable(table, cols, rels, "./models")
}

func GenerateStructWithMultiMeta(db *gorm.DB, table []string) error {
	for _, val := range table {

		cols, err := GetColumns(db, val)
		if err != nil {
			return err
		}

		rels, err := GetRelations(db, val)
		if err != nil {
			return err
		}

		fmt.Println("Generating struct for table:", table)
		GenerateStructTable(val, cols, rels, "./models")
	}
	return nil
}

// ==========================
// GET COLUMN META
// ==========================
func GetColumns(db *gorm.DB, table string) ([]Column, error) {
	driver := db.Dialector.Name()
	var cols []Column
	query := ""

	switch driver {
	case "mysql":
		query = `
			SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY = 'PRI' AS is_primary, IS_NULLABLE = 'YES' AS is_nullable
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`
	case "postgres":
		query = `
			SELECT column_name, data_type, (SELECT EXISTS (
				SELECT 1 FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
				WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_name = c.table_name AND kcu.column_name = c.column_name
			)) AS is_primary, is_nullable = 'YES' AS is_nullable
			FROM information_schema.columns c
			WHERE table_name = $1`
	case "sqlite":
		query = fmt.Sprintf("PRAGMA table_info(%s);", table)
	case "sqlserver":
		query = `
			SELECT c.name AS column_name, t.name AS data_type,
			CASE WHEN i.is_primary_key = 1 THEN 1 ELSE 0 END AS is_primary,
			CASE WHEN c.is_nullable = 1 THEN 1 ELSE 0 END AS is_nullable
			FROM sys.columns c
			JOIN sys.types t ON c.user_type_id = t.user_type_id
			LEFT JOIN sys.index_columns ic ON ic.column_id = c.column_id AND ic.object_id = c.object_id
			LEFT JOIN sys.indexes i ON i.object_id = ic.object_id AND i.index_id = ic.index_id
			WHERE c.object_id = OBJECT_ID(?)`
	case "oracle":
		query = `
			SELECT COLUMN_NAME, DATA_TYPE,
			CASE WHEN CONSTRAINT_TYPE = 'P' THEN 1 ELSE 0 END AS is_primary,
			CASE WHEN NULLABLE = 'Y' THEN 1 ELSE 0 END AS is_nullable
			FROM USER_TAB_COLUMNS
			LEFT JOIN USER_CONS_COLUMNS USING (TABLE_NAME, COLUMN_NAME)
			LEFT JOIN USER_CONSTRAINTS USING (CONSTRAINT_NAME)
			WHERE TABLE_NAME = UPPER(:1)`
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	rows, err := db.Raw(query, table).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Column
		switch driver {
		case "sqlite":
			var cid int
			var notnull, pk int
			var dfltValue, colType string
			rows.Scan(&cid, &c.Name, &colType, &notnull, &dfltValue, &pk)
			c.DataType = colType
			c.IsPrimary = pk == 1
			c.IsNullable = notnull == 0
		default:
			rows.Scan(&c.Name, &c.DataType, &c.IsPrimary, &c.IsNullable)
		}
		cols = append(cols, c)
	}

	return cols, nil
}

// ==========================
// GET RELATIONS
// ==========================
func GetRelations(db *gorm.DB, table string) ([]Relation, error) {
	driver := db.Dialector.Name()
	var rels []Relation
	query := ""

	switch driver {
	case "mysql":
		query = `
			SELECT TABLE_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
			WHERE TABLE_SCHEMA = DATABASE() AND REFERENCED_TABLE_NAME IS NOT NULL`
	case "postgres":
		query = `
			SELECT
				tc.table_name, kcu.column_name, ccu.table_name AS foreign_table_name, ccu.column_name AS foreign_column_name
			FROM information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
			WHERE constraint_type = 'FOREIGN KEY'`
	case "sqlite":
		query = fmt.Sprintf("PRAGMA foreign_key_list(%s);", table)
	case "sqlserver":
		query = `
			SELECT fk_tab.name AS table_name, fk_col.name AS column_name,
			       pk_tab.name AS referenced_table_name, pk_col.name AS referenced_column_name
			FROM sys.foreign_keys fk
			JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
			JOIN sys.tables fk_tab ON fk_tab.object_id = fkc.parent_object_id
			JOIN sys.columns fk_col ON fk_col.column_id = fkc.parent_column_id AND fk_col.object_id = fk_tab.object_id
			JOIN sys.tables pk_tab ON pk_tab.object_id = fkc.referenced_object_id
			JOIN sys.columns pk_col ON pk_col.column_id = fkc.referenced_column_id AND pk_col.object_id = pk_tab.object_id`
	case "oracle":
		query = `
			SELECT a.table_name, a.column_name, c_pk.table_name AS referenced_table_name, b.column_name AS referenced_column_name
			FROM user_constraints c
			JOIN user_cons_columns a ON c.constraint_name = a.constraint_name
			JOIN user_cons_columns b ON c.r_constraint_name = b.constraint_name
			JOIN user_constraints c_pk ON c.r_constraint_name = c_pk.constraint_name
			WHERE c.constraint_type = 'R'`
	}

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Relation
		switch driver {
		case "sqlite":
			var id, seq int
			var tableName, fromCol, toCol, onUpdate, onDelete, match string
			rows.Scan(&id, &seq, &tableName, &fromCol, &toCol, &onUpdate, &onDelete, &match)
			r = Relation{TableName: table, ForeignKey: fromCol, TargetTable: tableName, TargetColumn: toCol}
		default:
			rows.Scan(&r.TableName, &r.ForeignKey, &r.TargetTable, &r.TargetColumn)
		}
		r.RelationType = "many-to-one"
		rels = append(rels, r)
	}

	return rels, nil
}

// ==========================
// GENERATE STRUCT FILE
// ==========================
func GenerateStructTable(table string, columns []Column, relations []Relation, outputDir string) error {
	structName := toCamel(table)
	filename := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(structName))

	os.MkdirAll(outputDir, 0755)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := `package models

{{.Import}}

type {{.StructName}} struct {
{{- range .Columns }}
	{{ .FieldName }} {{ .FieldType }} ` + "`gorm:\"column:{{ .ColumnName }}{{ if .IsPrimary }};primaryKey{{ end }}\" json:\"{{ .ColumnName }}\"`" + `
{{- end }}
{{- if .Relations }}
{{- range .Relations }}
	{{ .FieldName }} {{ .TargetStruct }} ` + "`gorm:\"foreignKey:{{ .ForeignKey }};references:{{ .TargetColumn }}\"`" + `
{{- end }}
{{- end }}
}

func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

	type fieldMeta struct {
		FieldName  string
		FieldType  string
		ColumnName string
		IsPrimary  bool
	}

	type relMeta struct {
		FieldName    string
		ForeignKey   string
		TargetStruct string
		TargetColumn string
	}

	var importStatement string = "import (\n"

	var fields []fieldMeta
	for _, c := range columns {
		fields = append(fields, fieldMeta{
			FieldName:  toCamel(c.Name),
			FieldType:  mapSQLTypeToGo(c.DataType, c.IsNullable),
			ColumnName: c.Name,
			IsPrimary:  c.IsPrimary,
		})
		switch c.DataType {
		case "timestamp", "time", "date", "datetime":
			if !strings.Contains(importStatement, "time") {
				importStatement += `  "time"`
			}
			continue
		}
	}

	importStatement += "\n)"

	var rels []relMeta
	for _, r := range relations {
		if r.TableName == table && r.ForeignKey != "" && r.TargetTable != "" {
			rels = append(rels, relMeta{
				FieldName:    toCamel(r.TargetTable),
				ForeignKey:   r.ForeignKey,
				TargetStruct: toCamel(r.TargetTable),
				TargetColumn: r.TargetColumn,
			})
		}
	}

	data := map[string]interface{}{
		"StructName": structName,
		"Columns":    fields,
		"Relations":  rels,
		"Import":     importStatement,
		"TableName":  table, // <--- ini penting untuk TableName()
	}

	t := template.Must(template.New("struct").Parse(tmpl))
	return t.Execute(f, data)
}

// ==========================
// HELPERS
// ==========================
func toCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	re := regexp.MustCompile(`[_\s]+`)
	parts := re.Split(s, -1)
	for i, p := range parts {
		if p == "" {
			continue
		}
		switch strings.ToLower(p) {
		case "id":
			parts[i] = "ID"
			continue
		case "url":
			parts[i] = "URL"
			continue
		case "api":
			parts[i] = "API"
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}

func mapSQLTypeToGo(sqlType string, nullable bool) string {
	t := strings.ToLower(sqlType)
	var typ string

	switch {
	case strings.Contains(t, "int"):
		if strings.Contains(t, "tinyint") {
			typ = "int8"
		} else if strings.Contains(t, "bigint") {
			typ = "int64"
		} else {
			typ = "int"
		}

	case strings.Contains(t, "decimal"), strings.Contains(t, "numeric"),
		strings.Contains(t, "float"), strings.Contains(t, "double"), strings.Contains(t, "real"):
		typ = "float64"

	case strings.Contains(t, "bool"):
		typ = "bool"

	case strings.Contains(t, "date"), strings.Contains(t, "time"), strings.Contains(t, "timestamp"):
		typ = "time.Time"

	case strings.Contains(t, "char"), strings.Contains(t, "text"), strings.Contains(t, "clob"), strings.Contains(t, "uuid"):
		typ = "string"

	default:
		typ = "string"
	}

	// gunakan pointer kalau kolom nullable
	if nullable {
		if !strings.HasPrefix(typ, "*") {
			typ = "*" + typ
		}
	}

	return typ
}
