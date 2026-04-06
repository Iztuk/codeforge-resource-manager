package resources

import (
	"fmt"
	"resource-manager/internal/state"
	"strings"
)

func AddDb(name, address string) []error {
	var errors []error

	if name == "" || address == "" {
		if name == "" {
			errors = append(errors, fmt.Errorf("Please provide a name"))
		}
		if address == "" {
			errors = append(errors, fmt.Errorf("Please provide an address"))
		}

		return errors
	}

	s := strings.Split(address, "://")
	dialect, addr := strings.ToLower(s[0]), strings.ToLower(s[1])

	errors = append(errors, isValidDbConfig(name, dialect, addr)...)
	if len(errors) > 0 {
		return errors
	}

	switch dialect {
	case "sqlite":
		errors = append(errors, CheckSQLiteConnection(strings.TrimPrefix(addr, "sqlite://")))
	default:
		errors = append(errors, fmt.Errorf("Invalid database dialect: %s, select from the following: sqlite", dialect))
	}

	return errors
}

func isValidDbConfig(name, dialect, addr string) []error {
	var errors []error
	if name == "" {
		errors = append(errors, fmt.Errorf("Please provide a name"))
	} else if !isValidResourceName(name) {
		errors = append(errors, fmt.Errorf("Please provide a unique resource name"))
	}

	if dialect == "sqlite" || dialect == "postgres" {

	} else {
		errors = append(errors, fmt.Errorf("Please provide a valid DB dialect (sqlite || postgres)"))
	}

	if addr == "" {
		errors = append(errors, fmt.Errorf("Please provide a valid address"))
	}

	return errors
}

func isValidResourceName(name string) bool {
	_, ok := state.AppState.ResourceContract.Resources[name]
	if !ok {
		return false
	}

	return true
}
