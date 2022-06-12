package utils

import "os"

func GetEnv(name, value string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return value
}
