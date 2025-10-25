package configs

import "os"

func GetConfig(key string) string {
	return os.Getenv(key)
}
