package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"go/build"
)

type DatabaseType string

const (
	MongoDB  DatabaseType = "mongodb"
	Postgres DatabaseType = "postgres"
)

type Secrets struct {
	Port             string      `json:"PORT"`
	Environment      Environment `json:"ENVIRONMENT"`
	DatabaseUrl      string      `json:"DATABASE_URL"`
	DatabaseName     string      `json:"DATABASE_NAME"`
	CurrentDatabase  string      `json:"CURRENT_DATABASE"`
	JwtSecret        string      `json:"JWT_SECRET"`
	ElasticSearchUrl string      `json:"ELASTIC_SEARCH_URL"`
}

var secrets Secrets

const ServiceName = "api"
const defaultPort = "4000"

func init() {
	loadEnvironmentVariables()
	initializeSecrets()
}

// loadEnvironmentVariables loads the environment variables from the .env file.
func loadEnvironmentVariables() {
	importPath := fmt.Sprintf("%s/config", strings.ReplaceAll(ServiceName, "-", "."))
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err == nil {
		env := filepath.Join(p.Dir, "../.env")
		_ = godotenv.Load(env)
	}
}

// InitializeSecrets initializes the secrets.
func initializeSecrets() {
	secrets = Secrets{
		DatabaseUrl:      os.Getenv("DATABASE_URL"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		ElasticSearchUrl: os.Getenv("ELASTIC_SEARCH_URL"),
	}

	setCurrentDatabase()
	setEnvironment()
	setPort()
}

// setCurrentDatabase sets the current database.
func setCurrentDatabase() {
	currentDB := os.Getenv("CURRENT_DATABASE")
	switch DatabaseType(currentDB) {
	case MongoDB, Postgres:
		secrets.CurrentDatabase = currentDB
	default:
		log.Fatal("Invalid value for CURRENT_DATABASE. It must be either 'mongodb' or 'postgres'.")
	}
}

// setEnvironment sets the environment.
func setEnvironment() {
	if envStr := os.Getenv("ENVIRONMENT"); envStr != "" {
		if env := Environment(envStr); env.IsValid() != nil {
			log.Fatal("Error in environment variables: ", env.IsValid())
		}
	} else {
		log.Fatal("ENVIRONMENT is not set.")
	}
}

// setPort sets the port.
func setPort() {
	if port := os.Getenv("PORT"); port != "" {
		secrets.Port = port
	} else {
		secrets.Port = defaultPort
	}
}

// GetSecrets returns the initialized Secrets.
func GetSecrets() Secrets {
	return secrets
}
