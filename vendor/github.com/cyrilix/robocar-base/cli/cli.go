package cli

import (
	"log"
	"os"
	"strconv"
)

func SetDefaultValueFromEnv(value *string, key string, defaultValue string) {
	if os.Getenv(key) != "" {
		*value = os.Getenv(key)
	} else {
		*value = defaultValue
	}
}
func SetIntDefaultValueFromEnv(value *int, key string, defaultValue int) error {
	var sVal string
	if os.Getenv(key) != "" {
		sVal = os.Getenv(key)
		val, err := strconv.Atoi(sVal)
		if err != nil {
			log.Printf("unable to convert string to int: %v", err)
			return err
		}
		*value = val
	} else {
		*value = defaultValue
	}
	return nil
}
