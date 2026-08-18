package middleware

import (
	"net/http"
	"time"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	LocalAuthCookieName = "LocalAuthentication"
	LocalUserSub        = "FB6135C7-1464-4A72-B74E-4B63D343DD09"
	LocalUsername       = "bookstoreuser"
	JWTSecret           = "local-dev-secret-key-do-not-use-in-production"
)

// AuthRequired middleware - checks authentication based on mode
func AuthRequired(cfg *config.Config, customerService *service.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.AuthenticationMode == "local" {
			handleLocalAuth(c, customerService)
		} else {
			handleCognitoAuth(c, customerService)
		}
	}
}

func handleLocalAuth(c *gin.Context, customerService *service.CustomerService) {
	cookie, err := c.Cookie(LocalAuthCookieName)
	if err != nil || cookie == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		c.Abort()
		return
	}

	// Validate JWT token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		c.Abort()
		return
	}

	sub := claims["sub"].(string)
	username := claims["username"].(string)
	roles := []string{}
	if r, ok := claims["roles"].([]interface{}); ok {
		for _, role := range r {
			roles = append(roles, role.(string))
		}
	}

	// CreateOrUpdate customer on every authenticated request (matching source behavior)
	customer, err := customerService.CreateOrUpdate(sub, username, username, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user"})
		c.Abort()
		return
	}

	c.Set(string(UserSubKey), sub)
	c.Set(string(UserNameKey), username)
	c.Set(string(UserRolesKey), roles)
	c.Set(string(CustomerKey), customer)
	c.Next()
}

func handleCognitoAuth(c *gin.Context, customerService *service.CustomerService) {
	// Check for Authorization header with Bearer token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		c.Abort()
		return
	}

	tokenStr := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
		return
	}

	// In production, validate against Cognito JWKS
	// For now, parse without full verification (JWKS validation would be added)
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		c.Abort()
		return
	}

	sub := getStringClaim(claims, "sub")
	username := getStringClaim(claims, "cognito:username")
	firstName := getStringClaim(claims, "given_name")
	lastName := getStringClaim(claims, "family_name")
	email := getStringClaim(claims, "email")
	roles := getStringSliceClaim(claims, "cognito:groups")

	// CreateOrUpdate customer on every authenticated request
	customer, err := customerService.CreateOrUpdate(sub, username, firstName, lastName, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user"})
		c.Abort()
		return
	}

	c.Set(string(UserSubKey), sub)
	c.Set(string(UserNameKey), username)
	c.Set(string(UserRolesKey), roles)
	c.Set(string(CustomerKey), customer)
	c.Next()
}

// AdminRequired middleware - checks for Administrators role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GenerateLocalToken creates a JWT for local auth
func GenerateLocalToken() (string, time.Time, error) {
	expiry := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":      LocalUserSub,
		"username": LocalUsername,
		"roles":    []string{"Administrators"},
		"exp":      expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWTSecret))
	return tokenStr, expiry, err
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getStringSliceClaim(claims jwt.MapClaims, key string) []string {
	if v, ok := claims[key]; ok {
		switch val := v.(type) {
		case []interface{}:
			result := make([]string, 0, len(val))
			for _, item := range val {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		case []string:
			return val
		}
	}
	return nil
}
