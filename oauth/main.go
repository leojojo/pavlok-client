package main

import (
  "fmt"
  "os"
  "net/http"
  "golang.org/x/oauth2"
  "github.com/joho/godotenv"
  "github.com/gorilla/mux"
)

var pavlok = newConfig()

func newConfig() *oauth2.Config {
  err := godotenv.Load()
  if err != nil {
    fmt.Errorf("Cannot load .env: %s", err);
    os.Exit(1)
  }
  client_id, exists := os.LookupEnv("CLIENT_ID")
  secret_key, exists := os.LookupEnv("SECRET_KEY")
  if !exists {
    fmt.Errorf("Cannot find CLIENT_ID and SECRET_KEY in .env")
    os.Exit(1)
  }
  c := &oauth2.Config{
    RedirectURL: "https://zap.leojojo.me/oauth2/callback",
    ClientID: client_id,
    ClientSecret: secret_key,
    Endpoint: oauth2.Endpoint{
      AuthURL: "http://pavlok-mvp.herokuapp.com/oauth/authorize",
      TokenURL: "http://pavlok-mvp.herokuapp.com/oauth/token",
    },
  }
  return c
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", indexHandler)
  router.HandleFunc("/login", loginHandler)
  router.HandleFunc("/oauth2/callback", callbackHandler)
  http.ListenAndServe(":8888", router)
}
