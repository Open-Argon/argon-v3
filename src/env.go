package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// get the environment variables
func getEnv() ArObject {
	env := make(anymap)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = ArString(pair[1])
	}
	err := godotenv.Load(".env")
	if err == nil {
		values, err := godotenv.Read()
		if err == nil {
			for k, v := range values {
				env[k] = ArString(v)
			}
		}
	}
	return Map(env)
}

var env = getEnv()
