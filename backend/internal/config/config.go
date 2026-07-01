package config

import "os"

type Config struct {
	Port                      string
	AuthenticationMode        string
	DatabaseMode              string
	FileServiceMode           string
	ImageValidationMode       string
	LoggingServiceMode        string
	DatabaseDSN               string
	CognitoClientID           string
	CognitoMetadataURL        string
	CognitoDomain             string
	CognitoRedirectURL        string
	S3Bucket                  string
	CloudFrontDomain          string
	LocalImagePath            string
	LocalUserSub              string
	LocalUsername             string
}

func Load() *Config {
	return &Config{
		Port:                      getEnv("PORT", "8080"),
		AuthenticationMode:        getEnv("SERVICES_AUTHENTICATION", "local"),
		DatabaseMode:              getEnv("SERVICES_DATABASE", "local"),
		FileServiceMode:           getEnv("SERVICES_FILESERVICE", "local"),
		ImageValidationMode:       getEnv("SERVICES_IMAGEVALIDATIONSERVICE", "local"),
		LoggingServiceMode:        getEnv("SERVICES_LOGGINGSERVICE", "local"),
		DatabaseDSN:               getEnv("DATABASE_DSN", "bookstore.db"),
		CognitoClientID:           getEnv("COGNITO_CLIENT_ID", ""),
		CognitoMetadataURL:        getEnv("COGNITO_METADATA_URL", ""),
		CognitoDomain:             getEnv("COGNITO_DOMAIN", ""),
		CognitoRedirectURL:        getEnv("COGNITO_REDIRECT_URL", "http://localhost:8080/api/auth/callback"),
		S3Bucket:                  getEnv("S3_BUCKET", ""),
		CloudFrontDomain:          getEnv("CLOUDFRONT_DOMAIN", ""),
		LocalImagePath:            getEnv("LOCAL_IMAGE_PATH", "./uploads/coverimages"),
		LocalUserSub:              getEnv("LOCAL_USER_SUB", "FB6135C7-1464-4A72-B74E-4B63D343DD09"),
		LocalUsername:             getEnv("LOCAL_USERNAME", "bookstoreuser"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
