package main

import (
  "fmt"
  "time"
  "bytes"
  "net/http"
  "net/http/httputil"
  "crypto/rand"
  "encoding/base64"
  "golang.org/x/oauth2"
)

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

func testSendVib(w http.ResponseWriter, code string, token string) {
  var buf bytes.Buffer
  api_url := "http://pavlok-mvp.herokuapp.com/api/v1/stimuli/vibration/255?access_token=" + token
  fmt.Fprintf(w, "<a href='"+api_url+"'></a>")
  resp, err := http.Post(api_url, code, &buf)
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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
  code_query, ok := r.URL.Query()["code"]
  code := code_query[0]
  if !ok || len(code) < 1 {
    fmt.Errorf("Url Param 'code' is missing")
    return
  }

  token_exchanged, err := pavlok.Exchange(oauth2.NoContext, code)
  token := string(token_exchanged.AccessToken)
  if err != nil {
    fmt.Errorf("code exchange failed: %s", err.Error())
    return
  }
  fmt.Fprintf(w, "code is: %s\ntoken is: %s\n", string(code), token)

  testSendVib(w, code, token)
}
