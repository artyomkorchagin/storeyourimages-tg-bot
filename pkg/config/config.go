package config

import (
	"fmt"

	"github.com/artyomkorchagin/effectivemobile/pkg/helpers"
)

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		helpers.GetEnv("SERVER_HOST", ""),
		helpers.GetEnv("SERVER_PORT", ""),
		helpers.GetEnv("DB_NAME", ""),
		helpers.GetEnv("DB_USER", ""),
		helpers.GetEnv("DB_PASSWORD", ""),
		helpers.GetEnv("DB_SSLMODE", ""))
}
