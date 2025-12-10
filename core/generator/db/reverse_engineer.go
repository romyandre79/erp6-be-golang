package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ForeignKeyInfo represents a foreign key relationship
type ForeignKeyInfo struct {
	TableName        string
	ColumnName       string
	ReferencedTable  string
	ReferencedColumn string
}

// ReverseEngineeredTable represents a table extracted from the database
type ReverseEngineeredTable struct {
	Name        string
	Columns     []TableColumn
	ForeignKeys []ForeignKeyInfo
}

// GetAllTableNames retrieves all table names from the database
func GetAllTableNames(db *gorm.DB) ([]string, error) {
	driver := db.Dialector.Name()
	var tables []string
	var query string

	switch driver {
	case "mysql":
		query = "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'"
	case "postgres":
		query = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public'"
	case "sqlite", "sqlite3":
		query = "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
	case "sqlserver":
		query = "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'"
	case "oracle":
		query = "SELECT TABLE_NAME FROM USER_TABLES"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// MapSQLTypeToDesignerType converts SQL types to designer types
func MapSQLTypeToDesignerType(sqlType string, isPrimary bool, isAutoIncrement bool, driver string) string {
	sqlType = strings.ToLower(sqlType)

	// Check for auto-increment primary key
	if isPrimary && isAutoIncrement {
		return "auto"
	}

	// Map common types
	switch {
	case strings.Contains(sqlType, "int"):
		if strings.Contains(sqlType, "tinyint(1)") || strings.Contains(sqlType, "boolean") || strings.Contains(sqlType, "bool") {
			return "boolean"
		}
		return "number"
	case strings.Contains(sqlType, "decimal"), strings.Contains(sqlType, "numeric"),
		strings.Contains(sqlType, "float"), strings.Contains(sqlType, "double"), strings.Contains(sqlType, "real"):
		return "decimal"
	case strings.Contains(sqlType, "bool"):
		return "boolean"
	case strings.Contains(sqlType, "date") && !strings.Contains(sqlType, "time"):
		return "date"
	case strings.Contains(sqlType, "datetime"), strings.Contains(sqlType, "timestamp"):
		return "timestamp"
	case strings.Contains(sqlType, "time") && !strings.Contains(sqlType, "stamp"):
		return "time"
	case strings.Contains(sqlType, "text"), strings.Contains(sqlType, "clob"):
		return "longtext"
	case strings.Contains(sqlType, "char"), strings.Contains(sqlType, "varchar"), strings.Contains(sqlType, "varchar2"):
		return "text"
	default:
		return "text"
	}
}

// ReverseEngineerTable extracts the structure of a single table
func ReverseEngineerTable(db *gorm.DB, tableName string) (*ReverseEngineeredTable, error) {
	driver := db.Dialector.Name()
	var columns []TableColumn

	// Get column information
	type ColumnInfo struct {
		ColumnName    string
		DataType      string
		IsPrimary     bool
		IsNullable    bool
		ColumnDefault *string
		Extra         string
	}

	var colInfos []ColumnInfo
	var query string

	switch driver {
	case "mysql":
		query = `
			SELECT 
				COLUMN_NAME as column_name,
				DATA_TYPE as data_type,
				COLUMN_KEY = 'PRI' as is_primary,
				IS_NULLABLE = 'YES' as is_nullable,
				COLUMN_DEFAULT as column_default,
				EXTRA as extra
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION`
	case "postgres":
		query = `
			SELECT 
				column_name,
				data_type,
				(SELECT EXISTS (
					SELECT 1 FROM information_schema.table_constraints tc
					JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
					WHERE tc.constraint_type = 'PRIMARY KEY' 
					AND tc.table_name = c.table_name 
					AND kcu.column_name = c.column_name
				)) as is_primary,
				is_nullable = 'YES' as is_nullable,
				column_default,
				'' as extra
			FROM information_schema.columns c
			WHERE table_name = $1
			ORDER BY ordinal_position`
	case "sqlite", "sqlite3":
		// SQLite uses PRAGMA
		query = fmt.Sprintf("PRAGMA table_info('%s')", tableName)
		rows, err := db.Raw(query).Rows()
		if err != nil {
			return nil, fmt.Errorf("failed to query columns: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var cid, notnull, pk int
			var name, colType string
			var dfltValue *string
			if err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk); err != nil {
				return nil, fmt.Errorf("failed to scan column: %v", err)
			}

			isAutoIncrement := pk == 1 && strings.Contains(strings.ToLower(colType), "integer")
			designerType := MapSQLTypeToDesignerType(colType, pk == 1, isAutoIncrement, driver)

			col := TableColumn{
				Name: name,
				Type: designerType,
			}

			if notnull == 0 && designerType != "auto" {
				col.AllowNull = "true"
			} else {
				col.AllowNull = "false"
			}

			if dfltValue != nil && *dfltValue != "" {
				col.Default = *dfltValue
			}

			columns = append(columns, col)
		}
		goto ProcessForeignKeys

	case "sqlserver":
		query = `
			SELECT 
				c.name as column_name,
				t.name as data_type,
				CASE WHEN i.is_primary_key = 1 THEN 1 ELSE 0 END as is_primary,
				CASE WHEN c.is_nullable = 1 THEN 1 ELSE 0 END as is_nullable,
				dc.definition as column_default,
				CASE WHEN c.is_identity = 1 THEN 'auto_increment' ELSE '' END as extra
			FROM sys.columns c
			JOIN sys.types t ON c.user_type_id = t.user_type_id
			LEFT JOIN sys.index_columns ic ON ic.column_id = c.column_id AND ic.object_id = c.object_id
			LEFT JOIN sys.indexes i ON i.object_id = ic.object_id AND i.index_id = ic.index_id AND i.is_primary_key = 1
			LEFT JOIN sys.default_constraints dc ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
			WHERE c.object_id = OBJECT_ID(?)
			ORDER BY c.column_id`
	case "oracle":
		query = `
			SELECT 
				c.COLUMN_NAME as column_name,
				c.DATA_TYPE as data_type,
				CASE WHEN cc.CONSTRAINT_TYPE = 'P' THEN 1 ELSE 0 END as is_primary,
				CASE WHEN c.NULLABLE = 'Y' THEN 1 ELSE 0 END as is_nullable,
				c.DATA_DEFAULT as column_default,
				CASE WHEN c.IDENTITY_COLUMN = 'YES' THEN 'auto_increment' ELSE '' END as extra
			FROM USER_TAB_COLUMNS c
			LEFT JOIN USER_CONS_COLUMNS ucc ON c.TABLE_NAME = ucc.TABLE_NAME AND c.COLUMN_NAME = ucc.COLUMN_NAME
			LEFT JOIN USER_CONSTRAINTS cc ON ucc.CONSTRAINT_NAME = cc.CONSTRAINT_NAME AND cc.CONSTRAINT_TYPE = 'P'
			WHERE c.TABLE_NAME = UPPER(:1)
			ORDER BY c.COLUMN_ID`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err := db.Raw(query, tableName).Scan(&colInfos).Error; err != nil {
		return nil, fmt.Errorf("failed to query columns: %v", err)
	}

	// Convert to TableColumn format
	for _, info := range colInfos {
		isAutoIncrement := strings.Contains(strings.ToLower(info.Extra), "auto_increment") ||
			strings.Contains(strings.ToLower(info.Extra), "identity")

		designerType := MapSQLTypeToDesignerType(info.DataType, info.IsPrimary, isAutoIncrement, driver)

		col := TableColumn{
			Name: info.ColumnName,
			Type: designerType,
		}

		if info.IsNullable && designerType != "auto" {
			col.AllowNull = "true"
		} else {
			col.AllowNull = "false"
		}

		if info.ColumnDefault != nil && *info.ColumnDefault != "" {
			col.Default = strings.Trim(*info.ColumnDefault, "'")
		}

		columns = append(columns, col)
	}

ProcessForeignKeys:
	// Get foreign keys
	foreignKeys, err := ExtractForeignKeys(db, tableName)
	if err != nil {
		// Don't fail if foreign keys can't be extracted, just log and continue
		fmt.Printf("Warning: failed to extract foreign keys for table %s: %v\n", tableName, err)
		foreignKeys = []ForeignKeyInfo{}
	}

	return &ReverseEngineeredTable{
		Name:        tableName,
		Columns:     columns,
		ForeignKeys: foreignKeys,
	}, nil
}

// ExtractForeignKeys extracts foreign key relationships for a table
func ExtractForeignKeys(db *gorm.DB, tableName string) ([]ForeignKeyInfo, error) {
	driver := db.Dialector.Name()
	var foreignKeys []ForeignKeyInfo
	var query string

	switch driver {
	case "mysql":
		query = `
			SELECT 
				TABLE_NAME as table_name,
				COLUMN_NAME as column_name,
				REFERENCED_TABLE_NAME as referenced_table,
				REFERENCED_COLUMN_NAME as referenced_column
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
			WHERE TABLE_SCHEMA = DATABASE() 
			AND TABLE_NAME = ?
			AND REFERENCED_TABLE_NAME IS NOT NULL`
	case "postgres":
		query = `
			SELECT
				tc.table_name,
				kcu.column_name,
				ccu.table_name AS referenced_table,
				ccu.column_name AS referenced_column
			FROM information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
			WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1`
	case "sqlite", "sqlite3":
		query = fmt.Sprintf("PRAGMA foreign_key_list('%s')", tableName)
		rows, err := db.Raw(query).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var id, seq int
			var table, from, to, onUpdate, onDelete, match string
			if err := rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
				return nil, err
			}
			foreignKeys = append(foreignKeys, ForeignKeyInfo{
				TableName:        tableName,
				ColumnName:       from,
				ReferencedTable:  table,
				ReferencedColumn: to,
			})
		}
		return foreignKeys, nil
	case "sqlserver":
		query = `
			SELECT 
				fk_tab.name AS table_name,
				fk_col.name AS column_name,
				pk_tab.name AS referenced_table,
				pk_col.name AS referenced_column
			FROM sys.foreign_keys fk
			JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
			JOIN sys.tables fk_tab ON fk_tab.object_id = fkc.parent_object_id
			JOIN sys.columns fk_col ON fk_col.column_id = fkc.parent_column_id AND fk_col.object_id = fk_tab.object_id
			JOIN sys.tables pk_tab ON pk_tab.object_id = fkc.referenced_object_id
			JOIN sys.columns pk_col ON pk_col.column_id = fkc.referenced_column_id AND pk_col.object_id = pk_tab.object_id
			WHERE fk_tab.name = ?`
	case "oracle":
		query = `
			SELECT 
				a.table_name,
				a.column_name,
				c_pk.table_name AS referenced_table,
				b.column_name AS referenced_column
			FROM user_constraints c
			JOIN user_cons_columns a ON c.constraint_name = a.constraint_name
			JOIN user_cons_columns b ON c.r_constraint_name = b.constraint_name
			JOIN user_constraints c_pk ON c.r_constraint_name = c_pk.constraint_name
			WHERE c.constraint_type = 'R' AND a.table_name = UPPER(:1)`
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err := db.Raw(query, tableName).Scan(&foreignKeys).Error; err != nil {
		return nil, err
	}

	return foreignKeys, nil
}

// ReverseEngineerDatabase extracts all tables from the database
func ReverseEngineerDatabase(db *gorm.DB) ([]ReverseEngineeredTable, error) {
	tableNames, err := GetAllTableNames(db)
	if err != nil {
		return nil, err
	}

	var tables []ReverseEngineeredTable

	for _, tableName := range tableNames {
		// Skip system tables
		if strings.HasPrefix(strings.ToLower(tableName), "sys") ||
			strings.HasPrefix(strings.ToLower(tableName), "information_schema") {
			continue
		}

		table, err := ReverseEngineerTable(db, tableName)
		if err != nil {
			fmt.Printf("Warning: failed to reverse engineer table %s: %v\n", tableName, err)
			continue
		}

		tables = append(tables, *table)
	}

	return tables, nil
}

// ConvertToTableDefinition converts a ReverseEngineeredTable to TableDefinition format
func ConvertToTableDefinition(table *ReverseEngineeredTable, x, y int) *TableDefinition {
	tableDef := &TableDefinition{}
	tableDef.Table.Name = table.Name
	tableDef.Table.X = x
	tableDef.Table.Y = y
	tableDef.Table.Width = 240
	tableDef.Table.Columns = table.Columns
	tableDef.Table.IsPublished = false
	tableDef.Table.Comment = "Reverse engineered from database"

	return tableDef
}

// ConvertToJSON converts a TableDefinition to JSON string
func ConvertToJSON(tableDef *TableDefinition) (string, error) {
	jsonBytes, err := json.Marshal(tableDef)
	if err != nil {
		return "", fmt.Errorf("failed to marshal table definition: %v", err)
	}
	return string(jsonBytes), nil
}
