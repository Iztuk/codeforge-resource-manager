package state

import (
	"bytes"
	"encoding/json"
	"os"
	"resource-manager/internal/contracts"
)

type ApplicationState struct {
	ApiContract      contracts.OpenApiDoc
	ResourceContract contracts.ResourceDoc
}

var AppState ApplicationState

// NOTE: Add some validation logic in the future to prevent silent bugs involving broken contract structure
func (as *ApplicationState) InitializeAppState(api, resource string) error {
	a, err := os.ReadFile(api)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(a))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&as.ApiContract); err != nil {
		return err
	}

	r, err := os.ReadFile(resource)
	if err != nil {
		return err
	}

	dec = json.NewDecoder(bytes.NewReader(r))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&as.ResourceContract); err != nil {
		return err
	}

	return nil
}
