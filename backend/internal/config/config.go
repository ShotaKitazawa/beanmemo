package config

import "os"

type Config struct {
	DSN  string
	Port string
}

func Load() Config {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "beanmemo:beanmemo@tcp(localhost:3306)/beanmemo?parseTime=true"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{DSN: dsn, Port: port}
}
