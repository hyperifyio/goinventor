// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var (
	natsServer = flag.String("nats", getEnvOrDefault("GOINVENTOR_NATS_URL", nats.DefaultURL), "The NATS server URL")
	kvSource   = flag.String("source", getEnvOrDefault("GOINVENTOR_SOURCE", "env"), "The key value store source (env/nats)")
	list       = flag.Bool("list", false, "List all hosts")
	hostName   = flag.String("host", "", "Get host specific values")
)

func main() {

	flag.Parse()

	// Get NATS URL from environment variable or use default
	if *natsServer == "" {
		defaultURL := nats.DefaultURL
		natsServer = &defaultURL
	}

	// Initialize appropriate trigger
	switch *kvSource {

	case "env":

	case "nats":

	default:
		log.Fatalf("Unknown event type: %s", *kvSource)
	}

	if *list {

		// Example inventory data
		inventory := map[string]interface{}{
			"all": map[string]interface{}{
				"hosts": []string{"host1", "host2"},
				"vars": map[string]interface{}{
					"ansible_user": "admin",
				},
			},
		}

		// Convert inventory to JSON
		inventoryJSON, err := json.MarshalIndent(inventory, "", "    ")
		if err != nil {
			log.Fatalf("Error creating JSON inventory: %v", err)
		}

		if _, err = os.Stdout.Write(inventoryJSON); err != nil {
			fmt.Println("Error writing the JSON:", err)
			return
		}

	} else {

		hostVars := make(map[string]interface{})

		hostVars = map[string]interface{}{
			"ansible_host": "10.0.0.1",
			"ansible_user": "user1",
		}

		// Output host variables as JSON
		hostVarsJSON, err := json.Marshal(hostVars)
		if err != nil {
			log.Fatalf("Error creating JSON for host variables: %v", err)
		}

		if _, err = os.Stdout.Write(hostVarsJSON); err != nil {
			fmt.Println("Error writing the JSON:", err)
			return
		}

	}

}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
