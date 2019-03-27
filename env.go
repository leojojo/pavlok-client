package main

import (
  "fmt"
  "os"
  "log"
  "github.com/joho/godotenv"
)

func init() {
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}

func main() {
  client_id, exists := os.LookupEnv("CLIENT_ID")
  secret_key, exists := os.LookupEnv("SECRET_KEY")
  if exists {
    fmt.Println(client_id, secret_key)
  }
}
