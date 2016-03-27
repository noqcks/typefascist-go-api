package main

import (
  "fmt"
  "log"
  "strings"
  "net/http"
  "encoding/json"
  "path/filepath"
  "os/exec"
  "reflect"
  "os"
  "io"
  "io/ioutil"

  "github.com/gorilla/mux"
  // "github.com/dchest/uniuri"
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

type Test struct {
  Start string
}

func (t Test) OTF_TO_WOFF2() ([]byte, error) {
  cmd := exec.Command("./sfnt_to_woff2", t.Start)
  return cmd.Output()
}

func CallMethod(i interface{}, methodName string) interface{} {
  var ptr reflect.Value
  var value reflect.Value
  var finalMethod reflect.Value

  value = reflect.ValueOf(i)

  // if we start with a pointer, we need to get value pointed to
  // if we start with a value, we need to get a pointer to that value
  if value.Type().Kind() == reflect.Ptr {
    ptr = value
    value = ptr.Elem()
  } else {
    ptr = reflect.New(reflect.TypeOf(i))
    temp := ptr.Elem()
    temp.Set(value)
  }

  // check for method on value
  method := value.MethodByName(methodName)
  if method.IsValid() {
    finalMethod = method
  }
  // check for method on pointer
  method = ptr.MethodByName(methodName)
  if method.IsValid() {
    finalMethod = method
  }

  if (finalMethod.IsValid()) {
    return finalMethod.Call([]reflect.Value{})[0].Interface()
  }

  // return or panic, method not found of either type
  return ""
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
    return
  }


  /////////////
  /// GET FILE ////
  /////////////
  req.ParseMultipartForm(32 << 20)
  file, handler, err := req.FormFile("file")
  if err != nil {
    if err.Error() == "http: no such file" {
      resBody := &ErrorBody{
        Status: "400",
        Message: "Please upload a file",
      }
      response, err := json.Marshal(resBody)
      if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
      }
      res.Header().Set("Content-Type", "application/json")
      res.WriteHeader(http.StatusMethodNotAllowed)
      res.Write(response)
      return
    } else {
      fmt.Println(err.Error())
      http.Error(res, err.Error(), http.StatusInternalServerError)
      return
    }
  }
  defer file.Close()


  //////////////
  /// FORMAT ///
  //////////////
  convert_from := strings.Trim(filepath.Ext(handler.Filename), ".")
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
    return
  }


  /////////////////
  /// SAVE FILE ///
  /////////////////
  // fileName := uniuri.NewLen(12)
  s := []string{"/tmp/", handler.Filename}
  filePath := strings.Join(s, "")
  out, err := os.Create(filePath)
  if err != nil {
    fmt.Fprintf(res, "Unable to create the file for writing. Check your write access privilege")
    return
  }
  defer out.Close()
  _, err = io.Copy(out, file)

  s = []string{strings.ToUpper(convert_from), "_", "TO", "_",strings.ToUpper(convert_to)}
  font := strings.Join(s, "")

  i := Test{Start: handler.Filename}
  data := CallMethod(i, font)
  if data == "" {
    s = []string{"The conversion from ", convert_from, " to ", convert_to, " is not yet supported."}
    errMsg := strings.Join(s, "")
    resBody := &ErrorBody{
      Status: "400",
      Message: errMsg,
    }
    response, err := json.Marshal(resBody)
    if err != nil {
      http.Error(res, err.Error(), http.StatusInternalServerError)
      return
    }
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(http.StatusMethodNotAllowed)
    res.Write(response)
    return
  }
  final := strings.TrimSuffix(handler.Filename, filepath.Ext(handler.Filename))
  s = []string{"/tmp/", final, ".", convert_to}
  afterPath := strings.Join(s, "")
  dat, err := ioutil.ReadFile(afterPath)
  if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }
  res.Write(dat)
  return


  // fmt.Fprintf(w, data)
  // fmt.Println(data)
}
