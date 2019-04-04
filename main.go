package main

import (
  "./handlers"
  "./middleware"

  "net/http"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", handlers.GetJobs)

  // Wrap the servemux with the limit middleware.
  http.ListenAndServe(":4000", middleware.Limit(mux))
}
