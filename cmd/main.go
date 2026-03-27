package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"resource-manager/internal/state"
	"resource-manager/views"
)

func main() {
	cmd := flag.NewFlagSet("resource-manager", flag.ExitOnError)

	apiContract := cmd.String("api", "", "API contract file")
	resourceContract := cmd.String("resource", "", "Resource contract file")

	cmd.Parse(os.Args[1:])

	if *apiContract == "" || *resourceContract == "" {
		fmt.Println("Both --api and --resource flags are required")
		cmd.Usage()
		os.Exit(1)
	}

	// Initialize the application state
	err := state.AppState.InitializeAppState(*apiContract, *resourceContract)
	if err != nil {
		log.Fatal(err)
	}

	views.StartView()
}
