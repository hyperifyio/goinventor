// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

var (
	natsServer = flag.String("nats", getEnvOrDefault("GOINVENTOR_NATS_URL", nats.DefaultURL), "The NATS server URL")
	kvSource   = flag.String("source", getEnvOrDefault("GOINVENTOR_SOURCE", "env"), "The key value store source (env/nats)")
	list       = flag.Bool("list", false, "List all hosts")
	hostName   = flag.String("host", "", "Get host specific values")
	envPrefix  = flag.String("env-prefix", "INVENTORY_", "Prefix for ENV variables")
)

// InventoryItem represents a parsed inventory item with group, host, key, and value.
type InventoryItem struct {
	Group string
	Host  string
	Key   string
	Value string
}

func main() {

	flag.Parse()

	// Get NATS URL from environment variable or use default
	if *natsServer == "" {
		defaultURL := nats.DefaultURL
		natsServer = &defaultURL
	}

	var keyValues []string
	var inventory map[string]map[string]interface{}

	// Initialize appropriate trigger
	switch *kvSource {

	case "env":
		keyValues = filterKeyValuePairs(os.Environ(), *envPrefix)
		inventoryItems, parseErrors := parseInventoryItems(keyValues)
		if len(parseErrors) > 0 {
			for _, err := range parseErrors {
				log.Println("Parse error:", err)
			}
		}
		inventory = convertToAnsibleInventory(inventoryItems)

	//case "nats":

	default:
		log.Fatalf("Unknown event type: %s", *kvSource)
	}

	if *list {

		// Convert inventory to JSON
		inventoryJSON, err := json.MarshalIndent(inventory, "", "    ")
		if err != nil {
			log.Fatalf("Error creating JSON inventory: %v", err)
		}

		if _, err = os.Stdout.Write(inventoryJSON); err != nil {
			fmt.Println("Error writing the JSON:", err)
			return
		}

	} else if *hostName != "" {

		// Output host variables as JSON for the specified host
		var hostVarsJSON []byte
		var err error

		if meta, ok := inventory["_meta"]; ok {
			if hostvars, ok := meta["hostvars"]; ok {
				if hostVars, ok := hostvars.(map[string]map[string]interface{})[*hostName]; ok {
					hostVarsJSON, err = json.MarshalIndent(hostVars, "", "    ")
					if err == nil {
						if _, err = os.Stdout.Write(hostVarsJSON); err != nil {
							fmt.Println("Error writing the JSON:", err)
						}
						return
					}
				}
			}
		}

		// If the host is not found or an error occurred, return empty JSON object
		fmt.Println("{}")
	} else {
		fmt.Println("No operation specified (--list or --host)")
	}

}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// filterKeyValuePairs filters a slice of key-value strings by a prefix and
// returns a slice of strings with the prefix removed from each key.
func filterKeyValuePairs(keyValues []string, prefix string) []string {
	var filtered []string
	for _, kv := range keyValues {

		// Split each key-value pair into key and value
		keyVal := strings.SplitN(kv, "=", 2)
		if len(keyVal) != 2 {
			continue
		}
		key, value := keyVal[0], keyVal[1]

		// Check if the key starts with the given prefix
		if strings.HasPrefix(key, prefix) {
			trimmedKey := strings.TrimPrefix(key, prefix)
			filtered = append(filtered, trimmedKey+"="+value)
		}

	}
	return filtered
}

// parseInventoryItem parses a string in the format "group_hostname_key=value"
// and returns the corresponding InventoryItem struct.
//
// The format can be:
//
//	"group_hostname_key=value" (set "key=value" for a named hostname in the named group)
//	"_hostname_key=value" is same as "ungrouped_hostname_key=value" (set "key=value" for hostname without group)
//	"__key=value" is same as "ungrouped__key=value" (set "key=value" for any hostname without group)
//	"group__key=value" is same as "group_ungrouped_key=value" (set "key=value" in a named group)
func parseInventoryItem(item string) (InventoryItem, error) {

	parts := strings.SplitN(item, "=", 2)
	if len(parts) != 2 {
		return InventoryItem{}, fmt.Errorf("invalid format, expected key=value, got '%s'", item)
	}

	key, value := parts[0], parts[1]
	keyParts := strings.Split(key, "_")

	var group, host string

	switch {

	case len(keyParts) >= 3 && keyParts[0] == "":
		// Format: "_hostname_key=value" equivalent to "ungrouped_hostname_key=value"
		group, host = "ungrouped", keyParts[1]
		key = strings.Join(keyParts[2:], "_")

	case len(keyParts) >= 3 && keyParts[0] != "":
		// Format: "group_hostname_key=value"
		group, host = keyParts[0], keyParts[1]
		key = strings.Join(keyParts[2:], "_")

	default:
		return InventoryItem{}, fmt.Errorf("unrecognized key format: '%s'", key)
	}

	return InventoryItem{
		Group: group,
		Host:  host,
		Key:   key,
		Value: value,
	}, nil

}

// parseInventoryItems takes a slice of key-value strings and returns a slice of InventoryItem structs.
func parseInventoryItems(keyValues []string) ([]InventoryItem, []error) {
	var items []InventoryItem
	var errors []error

	for _, kv := range keyValues {
		item, err := parseInventoryItem(kv)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		items = append(items, item)
	}
	return items, errors
}

func convertToAnsibleInventory(items []InventoryItem) map[string]map[string]interface{} {
	groups := make(map[string]map[string]interface{})
	hostGroups := make(map[string]string) // Map to track which group a host belongs to
	groupVars := make(map[string]map[string]interface{})
	hostVars := make(map[string]map[string]interface{})

	// Initialize groups, host vars, and group vars
	for _, item := range items {
		if item.Host != "" {
			// Host-specific variables
			if _, exists := hostVars[item.Host]; !exists {
				hostVars[item.Host] = make(map[string]interface{})
			}
			hostVars[item.Host][item.Key] = item.Value

			// Track host's group
			if item.Group != "" && item.Group != "all" {
				hostGroups[item.Host] = item.Group
			}
		} else {
			// Group-specific variables
			if _, exists := groupVars[item.Group]; !exists {
				groupVars[item.Group] = make(map[string]interface{})
			}
			groupVars[item.Group][item.Key] = item.Value
		}
	}

	// Build groups with their hosts and vars
	for host, group := range hostGroups {
		if _, exists := groups[group]; !exists {
			groups[group] = map[string]interface{}{"hosts": []string{}, "vars": map[string]interface{}{}, "children": []string{}}
		}
		groups[group]["hosts"] = append(groups[group]["hosts"].([]string), host)
	}

	// Assign variables to groups
	for group, vars := range groupVars {
		if _, exists := groups[group]; !exists {
			groups[group] = map[string]interface{}{"hosts": []string{}, "vars": vars, "children": []string{}}
		} else {
			groups[group]["vars"] = vars
		}
	}

	// Include host vars in _meta
	groups["_meta"] = map[string]interface{}{"hostvars": hostVars}

	return groups
}

// extractHosts extracts host names from the inventory structure.
func extractHosts(hostsMap map[string]map[string]interface{}) []string {
	var hosts []string
	for host := range hostsMap {
		hosts = append(hosts, host)
	}
	return hosts
}
