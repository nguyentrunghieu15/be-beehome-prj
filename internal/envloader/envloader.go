package envloader

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func getEnv(config map[string]interface{}) map[string]interface{} {
	env := make(map[string]interface{})
	for k := range config {
		if v, ok := os.LookupEnv(k); ok {
			env[k] = v
		}
	}
	return env
}

func MustLoad(file string, config map[string]interface{}) error {
	if err := godotenv.Load(file); err != nil {
		return err
	}
	// Validate Env
	data := getEnv(config)
	if err := validator.ValidateMap(config, data); err != nil {
		return err
	}
	return nil
}
