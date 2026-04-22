package configs

import "os"

type ClerkConfig struct {
	ClerkSecret        string
	ClerkWebhookSecret string
}

func GetClerkConfig() *ClerkConfig {
	return &ClerkConfig{
		ClerkSecret:        os.Getenv("CLERK_SECRET"),
		ClerkWebhookSecret: os.Getenv("CLERK_WEBHOOK_SECRET"),
	}
}
