package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GitHubConfig struct {
	Username string
	APIKey   string
}

type DatabaseConfig struct {
	DBHost      string
	DBName     string
	DBUsername string
	DBPassword string
	DBType     string
	DBPort     string
}

type Config struct {
	DBConfig DatabaseConfig
	GitHub GitHubConfig
	SecretKey []byte
	DebugMode bool
	UserRoles []string
	MaxUsers  int
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		GitHub: GitHubConfig{
			Username: getEnv("GITHUB_USERNAME", ""),
			APIKey: getEnv("GITHUB_API_KEY", ""),
		},
		DBConfig: DatabaseConfig{
			DBHost:      getEnv("DB_HOST", ""),
			DBName:     getEnv("DB_NAME", ""),
			DBUsername: getEnv("DB_USER", ""),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBType:     getEnv("DB_TYPE", ""),
			DBPort:     getEnv("DB_PORT", ""),
		},
		SecretKey: getEnvAsByte("SECRET_KEY"),
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
		UserRoles: getEnvAsSlice("USER_ROLES", []string{"admin"}, ","),
		MaxUsers:  getEnvAsInt("MAX_USERS", 1),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		fmt.Println(value)
		return value
	}

	return defaultVal
}

func getEnvAsByte(key string) []byte {
	if value, exists := os.LookupEnv(key); exists {
		return []byte(value)
	}
	return []byte("")
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}