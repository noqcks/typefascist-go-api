package main

import (
  "fmt"
  "log"
  "strings"
  "net/http"
  "encoding/json"
  "path/filepath"

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

  /////////////
  /// FILE ////
  /////////////
  req.ParseMultipartForm(32 << 20)
  file, handler, err := req.FormFile("file")
  if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }
  // TODO: implement error handling for file upload
  // if file == nil {
  //     resBody := &ErrorBody{
  //     Status: "400",
  //     Message: "Please include a file to convert",
  //   }
  //   response, err := json.Marshal(resBody)
  //   if err != nil {
  //     http.Error(res, err.Error(), http.StatusInternalServerError)
  //     return
  //   }
  //   res.Header().Set("Content-Type", "application/json")
  //   res.WriteHeader(http.StatusBadRequest)
  //   res.Write(response)
  // }
  convert_from := strings.Trim(filepath.Ext(handler.Filename), ".")


  //////////////
  /// FORMAT ///
  //////////////
  convert_to := req.FormValue("format")
  if convert_to == "" {
      resBody := &ErrorBody{
      Status: "400",
      Message: "Please specify a format to convert to",
    }
    response, err := json.Marshal(resBody)
    if err != nil {
      http.Error(res, err.Error(), http.StatusInternalServerError)
      return
    }
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusBadRequest)
    res.Write(response)
  }
}
