package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"resource-manager/internal/contracts"
)

type ApplicationState struct {
	ApiContract      contracts.OpenApiDoc
	ResourceContract contracts.ResourceDoc
}

var AppState ApplicationState

var (
	apiContractFile      string
	resourceContractFile string
)

// NOTE: Add some validation logic in the future to prevent silent bugs involving broken contract structure
func (as *ApplicationState) InitializeAppState(api, resource string) error {
	a, err := os.ReadFile(api)
	if err != nil {
		return err
	}

	var apiContract contracts.OpenApiDoc
	if err := decodeStrictJSON(a, &apiContract); err != nil {
		return err
	}

	r, err := os.ReadFile(resource)
	if err != nil {
		return err
	}

	var resourceContract contracts.ResourceDoc
	if err := decodeStrictJSON(r, &resourceContract); err != nil {
		return err
	}

	as.ApiContract = apiContract
	as.ResourceContract = resourceContract

	apiContractFile = api
	resourceContractFile = resource

	return nil
}

func decodeStrictJSON(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(v); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("unexpected trailing data after JSON object")
	}

	return nil
}

func WriteToResourceFile() error {
	jsonData, err := json.MarshalIndent(AppState.ResourceContract, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(resourceContractFile, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
