package config

import "os"

// Config holds application configuration
type Config struct {
	Authentication         string
	Database               string
	FileService            string
	ImageValidationService string
	LoggingService         string
	DatabaseDSN            string
	Port                   string
	CognitoClientID        string
	CognitoMetadataAddress string
	CognitoDomain          string
	S3BucketName           string
	CloudFrontDomain       string
	LocalImagePath         string
}

// Load creates a Config from environment variables with local-first defaults
func Load() *Config {
	cfg := &Config{
		Authentication:         getEnvOrDefault("SERVICES_AUTHENTICATION", "local"),
		Database:               getEnvOrDefault("SERVICES_DATABASE", "local"),
		FileService:            getEnvOrDefault("SERVICES_FILESERVICE", "local"),
		ImageValidationService: getEnvOrDefault("SERVICES_IMAGEVALIDATIONSERVICE", "local"),
		LoggingService:         getEnvOrDefault("SERVICES_LOGGINGSERVICE", "local"),
		DatabaseDSN:            getEnvOrDefault("DATABASE_DSN", "bookstore.db"),
		Port:                   getEnvOrDefault("PORT", "8080"),
		CognitoClientID:        os.Getenv("COGNITO_CLIENT_ID"),
		CognitoMetadataAddress: os.Getenv("COGNITO_METADATA_ADDRESS"),
		CognitoDomain:          os.Getenv("COGNITO_DOMAIN"),
		S3BucketName:           os.Getenv("S3_BUCKET_NAME"),
		CloudFrontDomain:       os.Getenv("CLOUDFRONT_DOMAIN"),
		LocalImagePath:         getEnvOrDefault("LOCAL_IMAGE_PATH", "./uploads/coverimages"),
	}
	return cfg
}

// IsAWS returns true if the given service toggle is set to "aws"
func (c *Config) IsAWS(service string) bool {
	switch service {
	case "authentication":
		return c.Authentication == "aws"
	case "database":
		return c.Database == "aws"
	case "fileservice":
		return c.FileService == "aws"
	case "imagevalidation":
		return c.ImageValidationService == "aws"
	case "logging":
		return c.LoggingService == "aws"
	default:
		return false
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
