package cfg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var cfgFromEnv cfg

func Config() *cfg {
	return &cfgFromEnv
}

func init() {
	var errs []string

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfgFromEnv = cfg{
		PgHost:       envToString("PG_HOST", &errs),
		PgPort:       envToString("PG_PORT", &errs),
		PgDbName:     envToString("PG_DATABASE", &errs),
		PgUser:       envToString("PG_USER", &errs),
		PgPassword:   envToString("PG_PASSWORD", &errs),
		ServerPort:   envToString("SERVER_PORT", &errs),
	}

	if len(errs) > 0 {
		res := "\n" + strings.Join(errs, "\n")
		panic(res)
	}
}

type cfg struct {
	PgHost       string
	PgPort       string
	PgUser       string
	PgPassword   string
	PgDbName     string
	ServerPort   string
}

func envToString(key string, errs *[]string) string {
	_, exists := os.LookupEnv(key)

	if !exists {
		value := ""
		return value
	}

	value := os.Getenv(key)
	if len(value) == 0 {
		*errs = append(*errs, fmt.Sprintf("env variable not found. Key : %v", key))
	}
	return value
}
