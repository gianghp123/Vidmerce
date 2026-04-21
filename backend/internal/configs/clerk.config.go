package configs

import "os"

type ClerkConfig struct {
	ClerkSecret string
}

func GetClerkConfig() *ClerkConfig {
	return &ClerkConfig{
		ClerkSecret: os.Getenv("CLERK_SECRET"),
	}
}
