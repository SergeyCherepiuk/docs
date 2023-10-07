package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyCherepiuk/docs/pkg/http"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	e := http.Router{}.Build()
	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
