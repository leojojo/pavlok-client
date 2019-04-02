package main

import (
  "fmt"
  "os"
  "log"
  "time"
  "net/http"
  "crypto/rand"
  "encoding/base64"
  "golang.org/x/oauth2"
  "github.com/joho/godotenv"
  "github.com/gorilla/mux"
)

var oauthConf = &oauth2.Config{
  RedirectURL: "https://zap.leojojo.me/oauth2/callback",
  ClientID: os.Getenv("CLIENT_ID"),
  ClientSecret: os.Getenv("SECRET_KEY"),
  Scopes: []string{
    "https://pavlok-mvp.herokuapp.com/oauth/",
    "https://pavlok-mvp.herokuapp.com/api/",
  },
}

type CallbackRequest struct {
  Code  string `form:"code"`
  State string `form:"state"`
}

func init() {
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", indexHandler)
  router.HandleFunc("/login", loginHandler)
  router.HandleFunc("/oauth2/callback", callbackHandler)
  http.ListenAndServe(":8888", router)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "<a href='/login'>Log in with Pavlok</a>")
}

func generateStateOauthCookie(w http.ResponseWriter) string {
  var expiration = time.Now().Add(365 * 24 * time.Hour)
  b := make([]byte, 16)
  rand.Read(b)
  state := base64.URLEncoding.EncodeToString(b)
  cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
  http.SetCookie(w, &cookie)
  return state
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
  oauthState := generateStateOauthCookie(w)
  u := oauthConf.AuthCodeURL(oauthState)
  http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
  code, ok := r.URL.Query()["code"]
  if !ok || len(code[0]) < 1 {
    log.Println("Url Param 'code' is missing")
    return
  }
  fmt.Fprintln(w, "Url Param 'code' is: " + string(code[0]))
}
