package main

import (
  "encoding/json"
  "log"
  "net/http"
  "html"
  "regexp"
  "strings"
  "io/ioutil"
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
)

type test_struct struct {
  Test string
}

const (
  DB_USER     = "<user>"
  DB_PASSWORD = "<password>"
  DB_NAME     = "<name>"
)

// Old position request
func get_loc_old(w http.ResponseWriter, req *http.Request) {
  var ret = "Input not valid! "+req.URL.Path
  var str = html.EscapeString(req.URL.Path)
  // remove all up to the last splash
  var strarr = strings.Split(str, "/")
  str = strarr[len(strarr)-1]
  // allow only MAC addr without colon and separated by comma
  var validID = regexp.MustCompile(`^([[:xdigit:]]{12},){0,}[[:xdigit:]]{12}[,]?$`)
  if validID.MatchString(str) {
    strarr = strings.Split(str, ",")
    str = ""
    // Escape SQL string and set upper
    for _, e := range strarr {
      str = str+"'"+strings.ToUpper(e)+"',"
    }
    // rm last character if it comma
    if last := len(str) -1; last >= 0 && str[last] == ',' {
      str = str[:last]
    }
    // Open DB
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
    DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
      panic(err)
    }
    //TODO select quallity as well
    q := "SELECT BSSID, LAT, LON FROM bssid WHERE BSSID IN ("+str+")"
    rows, err := db.Query(q)
    if err != nil {
      panic(err)
    }
    defer db.Close()
    var retarr = [][]string{}
    for rows.Next() {
      var BSSID string
      var LAT string
      var LON string
      err = rows.Scan(&BSSID, &LAT, &LON)
      if err != nil {
        panic(err)
      }
      var tmp = []string{BSSID, LAT, LON}
      retarr = append(retarr, tmp)
    }
    //TODO filter unknown bssids and request seperatly by MLS
    //TODO save reuested informations if there are unknown e.g. signal qually, MLS stuff...
    //TODO calculate lat/lon
    fmt.Fprintf(w, "%q\n", retarr)
    return
  }
  fmt.Fprintf(w, "%q\n", ret)
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

//TODO read parameters from config
func main() {
  // need to handle by config
  var path = "/test"
  var old_path = "/api/v1/bssids/"
  var port = "8082"

  //TODO Multicore able
  http.HandleFunc(old_path, get_loc_old)

  http.HandleFunc(path, get_loc)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
