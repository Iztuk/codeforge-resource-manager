package resources

import (
	"fmt"
	"resource-manager/internal/contracts"
	"resource-manager/internal/state"
	"strings"
)

func AddDb(name, address string) []error {
	var errors []error

	if name == "" || address == "" {
		if name == "" {
			errors = append(errors, fmt.Errorf("please provide a name"))
		}
		if address == "" {
			errors = append(errors, fmt.Errorf("please provide an address"))
		}

		return errors
	}

	s := strings.Split(address, "://")
	dialect, addr := strings.ToLower(s[0]), strings.ToLower(s[1])

	errors = append(errors, isValidDbConfig(name, dialect, addr)...)
	if len(errors) > 0 {
		return errors
	}

	var normalizedTables map[string]contracts.DBTable
	switch dialect {
	case "sqlite":
		addr = strings.TrimPrefix(addr, "sqlite://")
		if err := CheckSQLiteConnection(addr); err != nil {
			errors = append(errors, err)
			return errors
		}

		tables, err := GetSQLiteTables(addr)
		if err != nil {
			errors = append(errors, err)
			return errors
		}

		if len(tables) > 0 {
			normalizedTables = NormalizeSQLiteTables(tables)
		}
	default:
		errors = append(errors, fmt.Errorf("Invalid database dialect: %s, select from the following: sqlite", dialect))
	}

	state.AppState.ResourceContract.Resources[name] = contracts.Resource{
		Name: name,
		Type: "database",
		DB: &contracts.DB{
			Dialect: dialect,
			Tables:  normalizedTables,
		},
	}

	if err := state.WriteToResourceFile(); err != nil {
		errors = append(errors, err)
		return errors
	}

	return nil
}

func isValidDbConfig(name, dialect, addr string) []error {
	var errors []error
	if name == "" {
		errors = append(errors, fmt.Errorf("please provide a name"))
	} else if !isValidResourceName(name) {
		errors = append(errors, fmt.Errorf("please provide a unique resource name"))
	}

	if dialect == "sqlite" || dialect == "postgres" {

	} else {
		errors = append(errors, fmt.Errorf("please provide a valid DB dialect (sqlite || postgres)"))
	}

	if addr == "" {
		errors = append(errors, fmt.Errorf("please provide a valid address"))
	}

	return errors
}

func isValidResourceName(name string) bool {
	_, ok := state.AppState.ResourceContract.Resources[name]
	if ok {
		return false
	}

	return true
}
