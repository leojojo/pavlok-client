package main

import (
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8888", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
  code, ok := r.URL.Query()["code"]
  if !ok || len(code[0]) < 1 {
    log.Println("Url Param 'code' is missing")
    return
  }
  log.Println("Url Param 'code' is: " + string(code[0]))
}
