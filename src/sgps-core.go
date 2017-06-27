package main

import (
  "encoding/json"
  "log"
  "net/http"
  "html"
  "regexp"
  "strings"
  "io/ioutil"
)

type test_struct struct {
  Test string
}

func get_loc_old(w http.ResponseWriter, req *http.Request) {
  var str = html.EscapeString(req.URL.Path)
  var strarr = strings.Split(str, "/")
  str = strarr[len(strarr)-1]
  // allow only MAC addr separated by ,
  var validID = regexp.MustCompile(`^([[:xdigit:]]{12},){0,}[[:xdigit:]]{12}$`)
  if validID.MatchString(str) {
    strarr = strings.Split(str, ",")
    for _, e := range strarr {
      log.Println(e)
    }
  }else {
    panic(str)
  }
}

func get_loc(rw http.ResponseWriter, req *http.Request) {
  body, err := ioutil.ReadAll(req.Body)
  if err != nil {
    panic(err)
  }
  //DEBUG: print get request
  log.Println(string(body))
  var t test_struct
  err = json.Unmarshal(body, &t)
  if err != nil {
    panic(err)
  }
  log.Println(t.Test)
}

func main() {
  // need to handle by config
  var path = "/test"
  var old_path = "/api/v1/bssids/"
  var port = "8082"

  http.HandleFunc(old_path, get_loc_old)

  http.HandleFunc(path, get_loc)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
