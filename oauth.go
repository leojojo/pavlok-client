package main

import (
  "fmt"
  "os"
  "time"
  "bytes"
  "net/http"
  "net/http/httputil"
  "crypto/rand"
  "encoding/base64"
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
  u := pavlok.AuthCodeURL(oauthState)
  http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
  code_query, ok := r.URL.Query()["code"]
  code := code_query[0]
  if !ok || len(code) < 1 {
    fmt.Errorf("Url Param 'code' is missing")
    return
  }

  token, err := pavlok.Exchange(oauth2.NoContext, code)
  if err != nil {
    fmt.Errorf("code exchange failed: %s", err.Error())
    return
  }
  fmt.Fprintf(w, "code is: %s\ntoken is: %s\n", string(code), string(token.AccessToken))
  var buf bytes.Buffer
  resp, err := http.Post("http://pavlok-mvp.herokuapp.com/api/v1/stimuli/vibration/255", code, &buf)
  if err != nil {
    fmt.Errorf("failed post: %s", err.Error())
    return
  }
  dumpResp, err := httputil.DumpResponse(resp, true)
  if err != nil {
    fmt.Errorf("failed post: %s", err.Error())
    return
  }
  fmt.Printf("%s", dumpResp)
}
