package config

import "os"

type Config struct {
	AuthMode            string
	DatabaseMode        string
	FileServiceMode     string
	ImageValidationMode string
	LoggingMode         string
	Port                string
	DatabaseDSN         string
	UploadDir           string
	CognitoClientID     string
	CognitoMetadataURL  string
	CognitoDomain       string
	S3Bucket            string
	CloudFrontDomain    string
}

func Load() *Config {
	cfg := &Config{
		AuthMode:            getEnv("SERVICES_AUTHENTICATION", "local"),
		DatabaseMode:        getEnv("SERVICES_DATABASE", "local"),
		FileServiceMode:     getEnv("SERVICES_FILESERVICE", "local"),
		ImageValidationMode: getEnv("SERVICES_IMAGEVALIDATIONSERVICE", "local"),
		LoggingMode:         getEnv("SERVICES_LOGGINGSERVICE", "local"),
		Port:                getEnv("PORT", "8080"),
		DatabaseDSN:         getEnv("DATABASE_DSN", "bookstore.db"),
		UploadDir:           getEnv("UPLOAD_DIR", "./uploads"),
		CognitoClientID:     getEnv("COGNITO_CLIENT_ID", ""),
		CognitoMetadataURL:  getEnv("COGNITO_METADATA_URL", ""),
		CognitoDomain:       getEnv("COGNITO_DOMAIN", ""),
		S3Bucket:            getEnv("S3_BUCKET", ""),
		CloudFrontDomain:    getEnv("CLOUDFRONT_DOMAIN", ""),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
