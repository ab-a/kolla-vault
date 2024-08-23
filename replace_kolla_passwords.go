package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	passwordsFile = "etc/kolla/passwords.yml"  // Path to your passwords.yml file
	vaultURLVar   = "http://127.0.0.1:8200"     // Placeholder for Vault URL variable
	basePath      = "secret/data/kolla"         // Base Vault path including '/data'
	defaultPath   = "default"                   // Default subdirectory for Vault path
)

func main() {
	// Determine the Vault path from the environment variable, or use the default if not set
	vaultSubDir := getEnv("VAULT_PATH", defaultPath)
	vaultPath := fmt.Sprintf("%s/%s", basePath, vaultSubDir)

	// Load the passwords.yml file from disk
	data, err := ioutil.ReadFile(passwordsFile)
	if err != nil {
		log.Fatalf("Failed to read passwords.yml: %v", err)
	}

	// Parse the YAML content into a map
	passwords := make(map[string]interface{})
	err = yaml.Unmarshal(data, &passwords)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Iterate over the parsed passwords map and replace each password with a Vault lookup
	for key, value := range passwords {
		switch v := value.(type) {
		case string:
			// If the value is a string and doesn't already contain a Vault lookup, update it
			if v == "" || !containsVaultLookup(v) {
				log.Printf("Updating key %s with Vault lookup.", key)
				passwords[key] = generateVaultLookup(vaultPath, key)
			}
		case map[interface{}]interface{}:
			// If the value is a nested map, process the nested structure
			log.Printf("Processing nested map for key %s", key)
			nestedVaultPath := fmt.Sprintf("%s/%s", vaultPath, key)
			passwords[key] = processNestedMap(v, nestedVaultPath)
		default:
			// Log unsupported types that aren't handled
			log.Printf("Unhandled type for key %s: %T", key, v)
		}
	}

	// Serialize the updated passwords map back to YAML with proper single-quote handling
	updatedData := manualYAMLSerialization(passwords)

	// Write the updated data back to passwords.yml file
	err = ioutil.WriteFile(passwordsFile, []byte(updatedData), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated passwords.yml: %v", err)
	}

	fmt.Println("passwords.yml has been updated with Vault secret references.")
}

// getEnv retrieves the value of an environment variable or returns a fallback value if not set
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// generateVaultLookup generates a Vault lookup string for a given variable name and path
func generateVaultLookup(vaultPath, variableName string) string {
	return fmt.Sprintf("{{ lookup('community.general.hashi_vault', '%s/%s', 'url=%s', token=lookup('env', 'VAULT_TOKEN')) }}",
		vaultPath, variableName, vaultURLVar)
}

// processNestedMap processes nested maps by generating Vault lookups for each nested key
func processNestedMap(nestedMap map[interface{}]interface{}, vaultPath string) map[string]string {
	updatedMap := make(map[string]string)
	for subKey, subValue := range nestedMap {
		subKeyStr := fmt.Sprintf("%v", subKey) // Convert key to string
		switch v := subValue.(type) {
		case string:
			// If the value is a string and doesn't already contain a Vault lookup, update it
			if v == "" || !containsVaultLookup(v) {
				updatedMap[subKeyStr] = generateVaultLookup(vaultPath, subKeyStr)
			}
		default:
			// Log unsupported nested types that aren't handled
			log.Printf("Skipping unsupported nested type for key %s: %T", subKeyStr, v)
		}
	}
	return updatedMap
}

// containsVaultLookup checks if the string already contains a Vault lookup
func containsVaultLookup(value string) bool {
	return strings.Contains(value, "hashi_vault")
}

// manualYAMLSerialization serializes the map back to YAML format with proper single-quote handling
func manualYAMLSerialization(data map[string]interface{}) string {
	var result strings.Builder
	for key, value := range data {
		// Convert each key-value pair to a YAML-formatted string
		result.WriteString(fmt.Sprintf("%s: '%v'\n", key, value))
	}
	return result.String()
}
