package handlers

import (
  "../models"

  "net/http"
  "encoding/json"
)

var jobs []models.Job
///add some job to the slice

func GetJobs(w http.ResponseWriter, r *http.Request) {
  jobs = append(jobs, models.Job{ID: 1, Name: "Accounting"})
  jobs = append(jobs, models.Job{ID: 2, Name: "Programming"})
  json.NewEncoder(w).Encode(jobs)
}
