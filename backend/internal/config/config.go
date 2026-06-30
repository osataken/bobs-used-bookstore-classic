package config

import "os"

type Config struct {
	Authentication         string
	Database               string
	FileService            string
	ImageValidationService string
	LoggingService         string
	DatabaseDSN            string
	Port                   string
	CognitoRegion          string
	CognitoUserPoolID      string
	CognitoClientID        string
	S3Bucket               string
	LocalImagePath         string
}

func Load() *Config {
	return &Config{
		Authentication:         getEnvOrDefault("SERVICES_AUTHENTICATION", "local"),
		Database:               getEnvOrDefault("SERVICES_DATABASE", "local"),
		FileService:            getEnvOrDefault("SERVICES_FILESERVICE", "local"),
		ImageValidationService: getEnvOrDefault("SERVICES_IMAGEVALIDATIONSERVICE", "local"),
		LoggingService:         getEnvOrDefault("SERVICES_LOGGINGSERVICE", "local"),
		DatabaseDSN:            getEnvOrDefault("DATABASE_DSN", "bookstore.db"),
		Port:                   getEnvOrDefault("PORT", "8080"),
		CognitoRegion:          getEnvOrDefault("COGNITO_REGION", "us-east-1"),
		CognitoUserPoolID:      getEnvOrDefault("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:        getEnvOrDefault("COGNITO_CLIENT_ID", ""),
		S3Bucket:               getEnvOrDefault("S3_BUCKET", ""),
		LocalImagePath:         getEnvOrDefault("LOCAL_IMAGE_PATH", "./uploads"),
	}
}

func (c *Config) IsLocalAuth() bool {
	return c.Authentication != "aws"
}

func (c *Config) IsLocalDB() bool {
	return c.Database != "aws"
}

func (c *Config) IsLocalFileService() bool {
	return c.FileService != "aws"
}

func (c *Config) IsLocalImageValidation() bool {
	return c.ImageValidationService != "aws"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
