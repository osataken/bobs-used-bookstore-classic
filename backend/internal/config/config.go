package config

import "os"

// Config holds all application configuration
type Config struct {
	Port                   string
	AuthenticationMode     string // "local" or "aws"
	DatabaseMode           string // "local" or "aws"
	FileServiceMode        string // "local" or "aws"
	ImageValidationMode    string // "local" or "aws"
	LoggingMode            string // "local" or "aws"

	// Database
	DatabaseDSN string

	// AWS Cognito
	CognitoClientID    string
	CognitoUserPoolID  string
	CognitoDomain      string
	CognitoRegion      string
	CognitoRedirectURL string

	// AWS S3
	S3Bucket       string
	CloudFrontURL  string

	// Local file storage
	LocalStoragePath string

	// Server
	FrontendURL string
}

func Load() *Config {
	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		AuthenticationMode:  getEnv("SERVICES_AUTHENTICATION", "local"),
		DatabaseMode:        getEnv("SERVICES_DATABASE", "local"),
		FileServiceMode:     getEnv("SERVICES_FILESERVICE", "local"),
		ImageValidationMode: getEnv("SERVICES_IMAGEVALIDATIONSERVICE", "local"),
		LoggingMode:         getEnv("SERVICES_LOGGINGSERVICE", "local"),

		DatabaseDSN: getEnv("DATABASE_DSN", "bookstore.db"),

		CognitoClientID:    getEnv("COGNITO_CLIENT_ID", ""),
		CognitoUserPoolID:  getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoDomain:      getEnv("COGNITO_DOMAIN", ""),
		CognitoRegion:      getEnv("COGNITO_REGION", "us-east-1"),
		CognitoRedirectURL: getEnv("COGNITO_REDIRECT_URL", "http://localhost:8080/api/auth/callback"),

		S3Bucket:      getEnv("S3_BUCKET", ""),
		CloudFrontURL: getEnv("CLOUDFRONT_URL", ""),

		LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "./uploads"),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
