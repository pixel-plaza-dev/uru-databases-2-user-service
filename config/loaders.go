package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LoadUsersServicePort load Users Service port from environment variables
func LoadUsersServicePort() (string, string) {
	// Get environment variable
	port, exists := os.LookupEnv(UsersServicePortKey)
	if !exists {
		log.Fatalf("Port not found at '%s' environment variable", UsersServicePortKey)
	}

	var portBuilder strings.Builder
	portBuilder.WriteString(":")
	portBuilder.WriteString(port)

	fmt.Println(portBuilder.String())

	return port, portBuilder.String()
}
