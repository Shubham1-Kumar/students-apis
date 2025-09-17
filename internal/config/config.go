// Config Management Approaches:
// 1. Config Files (YAML/JSON) – good for dev/local.
// 2. Environment Variables – standard in Docker/K8s (best for production).
// 3. Command-line Flags – good for tools or overrides.
// 4. Remote Config Services – for dynamic/microservices setups.
// 5. Secrets Manager – for sensitive data (DB creds, API keys).
//
// Best Practice (Cloud/K8s):
// - Use env vars for core configs (via ConfigMap).
// - Use a secrets manager for sensitive data (Vault, K8s Secrets, AWS/GCP).
// - Optional config files only for structured/local configs.
// - Override order: defaults → file → env → flags.

package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"producton"` // these are called struct tags
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// if it's not there then we will check for the cmd flags/arguments
		// these are passed when we run the program

		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}

	}

	// now check is there any file available on that path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cfg
}
