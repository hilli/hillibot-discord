package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Token string
)

func Load() {
	godotenv.Load(".env")
	Token = os.Getenv("DISCORD_TOKEN")
}
