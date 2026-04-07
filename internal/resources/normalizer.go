package resources

import (
	"resource-manager/internal/contracts"
	"strings"
)

func NormalizeSQLiteTables(sqliteTables []SQLiteTable) map[string]contracts.DBTable {
	tables := make(map[string]contracts.DBTable)

	for _, t := range sqliteTables {
		dbTable := contracts.DBTable{
			PrimaryKey: []string{},
			Fields:     make(map[string]contracts.FieldSpec),
			UniqueKeys: t.UniqueKeys,
		}

		for _, c := range t.Cols {
			// Nullable rule:
			nullable := !c.NotNull
			if c.PrimaryKey {
				nullable = false
			}

			// Default rule:
			var def *string
			if c.Default.Valid {
				s := c.Default.String
				def = &s
			}

			fs := contracts.FieldSpec{
				ColumnName: c.Name,
				Type:       CanonicalTypeSQLite(c.Type, c.Name),
				Nullable:   nullable,
				Default:    def,
				Read:       true,
				Write:      true,
				Mutable:    true,
			}

			if c.PrimaryKey {
				dbTable.PrimaryKey = append(dbTable.PrimaryKey, c.Name)
			}

			dbTable.Fields[c.Name] = fs
		}

		tables[t.Name] = dbTable
	}

	return tables
}

func CanonicalTypeSQLite(typeRaw, colName string) string {
	t := strings.ToUpper(strings.TrimSpace(typeRaw))
	n := strings.ToLower(strings.TrimSpace(colName))

	// Explicit bool types
	if strings.Contains(t, "BOOL") {
		return "boolean"
	}

	// Numeric affinity
	if strings.Contains(t, "INT") ||
		strings.Contains(t, "REAL") ||
		strings.Contains(t, "FLOA") ||
		strings.Contains(t, "DOUB") ||
		strings.Contains(t, "NUM") ||
		strings.Contains(t, "DEC") {

		// Boolean-by-name heuristic on numeric columns (optional)
		if strings.HasPrefix(n, "is_") ||
			strings.HasPrefix(n, "has_") ||
			strings.HasPrefix(n, "can_") ||
			strings.HasPrefix(n, "should_") ||
			strings.HasPrefix(n, "enabled_") ||
			strings.HasPrefix(n, "active_") ||
			strings.HasSuffix(n, "_flag") ||
			strings.HasSuffix(n, "_enabled") ||
			strings.HasSuffix(n, "_active") {
			return "boolean"
		}

		return "number"
	}

	// Text/date-ish affinity
	if strings.Contains(t, "CHAR") ||
		strings.Contains(t, "CLOB") ||
		strings.Contains(t, "TEXT") ||
		strings.Contains(t, "VARCHAR") ||
		strings.Contains(t, "DATE") ||
		strings.Contains(t, "TIME") ||
		strings.Contains(t, "UUID") ||
		strings.Contains(t, "JSON") {
		return "string"
	}

	// Blob / none / unknown
	if strings.Contains(t, "BLOB") || t == "" {
		return "string"
	}

	// Default fallback
	return "string"
}
