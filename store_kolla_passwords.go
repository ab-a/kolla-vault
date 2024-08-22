package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	passwordsFile = "etc/kolla/passwords.yml" // Path to the passwords.yml file
	vaultAddress  = "http://127.0.0.1:8200"    // Address of the Vault server
	basePath      = "secret/data/kolla"        // Base Vault path, includes the 'data' section
	defaultPath   = "default"                  // Default subdirectory to use in Vault path
)

func main() {
	// Determine the Vault subdirectory from the environment variable or use the default
	vaultSubDir := getEnv("VAULT_PATH", defaultPath)
	vaultPath := fmt.Sprintf("%s/%s", basePath, vaultSubDir)
	vaultToken := os.Getenv("VAULT_TOKEN")

	// Ensure VAULT_TOKEN is set, or exit
	if vaultToken == "" {
		log.Println("VAULT_TOKEN is not set. Exiting.")
		return
	}

	// Initialize the Vault client
	client, err := initializeVaultClient(vaultAddress, vaultToken)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	// Load the passwords.yml file into memory
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

	// Store the parsed passwords in Vault
	if err := storePasswordsInVault(client, vaultPath, passwords); err != nil {
		log.Fatalf("Failed to store passwords: %v", err)
	}

	log.Println("All passwords have been stored in Vault.")
}

// getEnv retrieves the value of an environment variable or returns a fallback value if not set
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		fmt.Printf("Using Vault path '%s/%s'.\n", basePath, value)
		return value
	}
	fmt.Printf("VAULT_PATH is not set. Using default path '%s/%s'.\n", basePath, fallback)
	return fallback
}

// initializeVaultClient creates a new Vault client and sets the authentication token
func initializeVaultClient(address, token string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{Address: address})
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return client, nil
}

// storePasswordsInVault stores the passwords recursively in Vault
func storePasswordsInVault(client *api.Client, vaultPath string, data map[string]interface{}) error {
	// Delegate the storage process to storeNestedMap to handle both simple and nested values
	return storeNestedMap(client, vaultPath, data)
}

// storeNestedMap recursively processes and stores nested structures in Vault
func storeNestedMap(client *api.Client, path string, data map[string]interface{}) error {
	for key, value := range data {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			// Convert map[interface{}]interface{} to map[string]interface{} for compatibility with Vault
			nestedMap := make(map[string]interface{})
			for k, val := range v {
				strKey := fmt.Sprintf("%v", k)
				nestedMap[strKey] = val
			}
			// Recursively store nested maps
			if err := storeNestedMap(client, fmt.Sprintf("%s/%s", path, key), nestedMap); err != nil {
				return err
			}
		case string:
			// Store simple string values as secrets in Vault
			if err := storeSecret(client, path, key, v); err != nil {
				return err
			}
		case nil:
			// Handle nil values (currently skipping them)
			log.Printf("Skipping nil value for key %s", key)
		default:
			// Skip unsupported types
			log.Printf("Skipping unsupported type for key %s: %T", key, v)
		}
	}
	return nil
}

// storeSecret stores a single key-value pair in Vault at the specified path
func storeSecret(client *api.Client, vaultPath, key, value string) error {
	// Write the secret to Vault in the appropriate path
	_, err := client.Logical().Write(
		fmt.Sprintf("%s/%s", vaultPath, key),
		map[string]interface{}{"data": map[string]interface{}{"value": value}},
	)
	if err != nil {
		log.Printf("Error storing secret %s at path %s/%s: %v", key, vaultPath, key, err)
	}
	return err
}
