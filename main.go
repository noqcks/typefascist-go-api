package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"

  "github.com/gorilla/mux"
)

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", convert)
  log.Fatal(http.ListenAndServe(":8080", router))
}

type ErrorBody struct {
  Status  string  `json:"status"`
  Message string  `json:"message"`
}

func convert(res http.ResponseWriter, req *http.Request) {

  // return 405 error
  if req.Method != "POST" {
    resBody := &ErrorBody{
      Status: "405",
      Message: "The only available method is POST",
    }
    response, err := json.Marshal(resBody)
    if err != nil {
      http.Error(res, err.Error(), http.StatusInternalServerError)
      return
    }
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusMethodNotAllowed)
    res.Write(response)
  }

}
