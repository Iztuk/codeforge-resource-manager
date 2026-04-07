package resources

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteTable struct {
	Name string
	Cols []SQLiteColumn

	UniqueKeys [][]string
}

type SQLiteColumn struct {
	CID        int
	Name       string
	Type       string
	NotNull    bool
	Default    sql.NullString
	PrimaryKey bool
}

type sqliteIndexRow struct {
	Seq     int
	Name    string
	Unique  int
	Origin  string
	Partial int
}

type sqliteIndexInfoRow struct {
	SeqNo int
	CID   int
	Name  string
}

func CheckSQLiteConnection(path string) error {
	dsn := fmt.Sprintf("file:%s?mode=rw", path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Connection validation
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func GetSQLiteTables(path string) ([]SQLiteTable, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return []SQLiteTable{}, err
	}
	defer db.Close()

	var tableNames []string
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err != nil {
		return []SQLiteTable{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return []SQLiteTable{}, err
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return []SQLiteTable{}, err
	}

	var sqliteTables []SQLiteTable
	for _, tableName := range tableNames {
		var table SQLiteTable = SQLiteTable{
			Name: tableName,
		}
		rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
		if err != nil {
			return []SQLiteTable{}, err
		}

		for rows.Next() {
			var cid int
			var name, typ string
			var notnull int
			var dflt sql.NullString
			var pk int

			if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
				return []SQLiteTable{}, err
			}

			var col SQLiteColumn = SQLiteColumn{
				CID:        cid,
				Name:       name,
				Type:       typ,
				NotNull:    intToBool(notnull),
				Default:    dflt,
				PrimaryKey: intToBool(pk),
			}

			table.Cols = append(table.Cols, col)
		}
		rows.Close()

		// --- unique keys (excluding pk) ---
		idxRows, err := db.Query(fmt.Sprintf("PRAGMA index_list(%s)", tableName))
		if err != nil {
			return []SQLiteTable{}, err
		}

		var uniqueIndexNames []string
		for idxRows.Next() {
			var r sqliteIndexRow
			if err := idxRows.Scan(&r.Seq, &r.Name, &r.Unique, &r.Origin, &r.Partial); err != nil {
				return []SQLiteTable{}, err
			}

			if r.Unique == 1 && r.Origin != "pk" {
				uniqueIndexNames = append(uniqueIndexNames, r.Name)
			}
		}
		idxRows.Close()

		for _, idxName := range uniqueIndexNames {
			infoRows, err := db.Query(fmt.Sprintf("PRAGMA index_info(%s)", idxName))
			if err != nil {
				return []SQLiteTable{}, err
			}

			colsBySeq := map[int]string{}
			maxSeq := -1
			for infoRows.Next() {
				var ir sqliteIndexInfoRow
				if err := infoRows.Scan(&ir.SeqNo, &ir.CID, &ir.Name); err != nil {
					return []SQLiteTable{}, err
				}
				colsBySeq[ir.SeqNo] = ir.Name
				if ir.SeqNo > maxSeq {
					maxSeq = ir.SeqNo
				}
			}
			if err := infoRows.Err(); err != nil {
				return []SQLiteTable{}, err
			}
			infoRows.Close()

			// ordered columns for this unique key
			keyCols := make([]string, 0, maxSeq+1)
			for i := 0; i <= maxSeq; i++ {
				if c, ok := colsBySeq[i]; ok {
					keyCols = append(keyCols, c)
				}
			}

			if len(keyCols) > 0 {
				table.UniqueKeys = append(table.UniqueKeys, keyCols)
			}
		}

		sqliteTables = append(sqliteTables, table)
	}

	if len(sqliteTables) == 0 {
		return []SQLiteTable{}, nil
	} else {
		return sqliteTables, nil
	}
}

func intToBool(i int) bool {
	return i != 0
}
