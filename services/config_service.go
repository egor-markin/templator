package services

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

// Based on
// https://towardsdatascience.com/use-environment-variable-in-your-next-golang-project-39e17c3aaa66
// https://github.com/spf13/viper

func GetHttpPort() int {
	return getIntValue("HTTP_PORT", 80)
}

func GetAdminSessionTokenTTL() int {
	return getIntValue("ADMIN_SESSION_TOKEN_TTL", 86400)
}

func GetAdminPassword() string {
	return getStringValue("ADMIN_PASSWORD", "hJZfZ1wZwt0mRQSi")
}

func init() {
	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error while reading config file %s", err)
	}
}

func getIntValue(key string, defaultValue int) int {
	value, ok := viper.Get(key).(string)
	if !ok {
		return defaultValue
	}

	// Converting retrieved value into integer
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getStringValue(key string, defaultValue string) string {
	value, ok := viper.Get(key).(string)
	if !ok {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
