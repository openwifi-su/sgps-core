package main

import (
  "os"
  "strconv"
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
  "github.com/spf13/viper"
  "math"
  "bytes"
)

type MLS struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Accuracy float64 `json:"accuracy"`
}

func mls_request(apikey string, bssids []string) (ret MLS){
  url := "https://location.services.mozilla.com/v1/geolocate?key="+apikey
  var reg_str = ""
  for _, bssid := range bssids {
    reg_str = reg_str+`{ "macAddress": "`+bssid+`" },`
  }
  // rm last character if it comma
  if last := len(reg_str) -1; last >= 0 && reg_str[last] == ',' {
    reg_str = reg_str[:last]
  }
  var jsonStr = []byte(`{ "wifiAccessPoints": [`+reg_str+`], "fallbacks": {"lacf": false, "ipf": false }}`)
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil { panic(err) }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil { panic(err) }
  err = json.Unmarshal(body, &ret)
  if err != nil { panic(err) }
  return
}

// convert degre to radiate
func toRadian(deg float64) (rad float64){
  rad = deg * math.Pi / 180
  return
}

// convert radiate to degre
func toDegres(rad float64) (deg float64){
  deg = rad / math.Pi * 180
  return
}

// middle two positions
func mid_position(lat0 float64, lon0 float64, lat1 float64, lon1 float64) (lat2, lon2 float64){
  var dlon = toRadian(lon1 - lon0)

  lat0 = toRadian(lat0)
  lat1 = toRadian(lat1)
  lon0 = toRadian(lon0)

  var Bx = math.Cos(lat1) * math.Cos(dlon)
  var By = math.Cos(lat1) * math.Sin(dlon)
  lat2 = math.Atan2(math.Sin(lat0) + math.Sin(lat1), math.Sqrt((math.Cos(lat0) + Bx) * (math.Cos(lat0) + Bx) + By * By))
  lon2 = lon0 + math.Atan2(By, math.Cos(lat0) + Bx)
  lat2 = toDegres(lat2)
  lon2 = toDegres(lon2)
  return
}

func filter_unknown_bssid(arr [][]string, req []string) (ret []string){
  for _, req_elem := range req {
    var isAval = false
    for _, db_elem := range arr {
      if strings.ToUpper(req_elem) == strings.ToUpper(db_elem[0]) {
	isAval = true
      }
    }
    if isAval == false {
      ret = append(ret, strings.ToUpper(req_elem))
    }
  }
  return
}

// Old position request
func get_loc_old(w http.ResponseWriter, req *http.Request, config [4]string) {
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
    config[0], config[1], config[2])
    db, err := sql.Open("postgres", dbinfo)
    if err != nil { panic(err) }
    //TODO select quallity as well
    q := "SELECT BSSID, LAT, LON FROM bssid WHERE BSSID IN ("+str+")"
    rows, err := db.Query(q)
    if err != nil { panic(err) }
    defer db.Close()
    var retarr = [][]string{}
    for rows.Next() {
      var BSSID string
      var LAT string
      var LON string
      err = rows.Scan(&BSSID, &LAT, &LON)
      if err != nil { panic(err) }
      var tmp = []string{BSSID, LAT, LON}
      retarr = append(retarr, tmp)
    }
    //filter unknown bssids for request seperatly by MLS
    var unknown_bssid = filter_unknown_bssid(retarr, strarr)
    if len(unknown_bssid) > 1 {
      var mls = mls_request(config[3], unknown_bssid)
      if mls.Location.Lat != 0 && mls.Location.Lng != 0 {
	var tmp = []string{"", strconv.FormatFloat(mls.Location.Lat, 'f', 9, 64) , strconv.FormatFloat(mls.Location.Lng, 'f', 9, 64)}
	retarr = append(retarr, tmp)
      }
    }
    //TODO save reuested informations if there are unknown e.g. signal qually, MLS stuff...
    if len(retarr) > 1 {
      lat0, _ := strconv.ParseFloat(retarr[0][1], 64)
      lon0, _ := strconv.ParseFloat(retarr[0][2], 64)
      for _, elem := range retarr {
	tmpLat, _ := strconv.ParseFloat(elem[1], 64)
	tmpLon, _ := strconv.ParseFloat(elem[2], 64)
	lat0, lon0 = mid_position(lat0, lon0, tmpLat, tmpLon)
      }
      ret = strconv.FormatFloat(lat0,'f', 9,64)+","+strconv.FormatFloat(lon0,'f', 9,64)
    } else if len(retarr) > 0 {
      ret = retarr[0][1]+","+retarr[0][2]
    }
  }
  fmt.Fprintf(w, "%q\n", ret)
}

func main() {
  viper.SetConfigName("sgps")
  viper.AddConfigPath("config")
  err := viper.ReadInConfig()
  if err != nil {
    fmt.Println("Config file not found...")
    os.Exit(1)
  }
  var db [4]string
  db[0] = viper.GetString("database.db_user")
  db[1] = viper.GetString("database.db_password")
  db[2] = viper.GetString("database.db_name")
  db[3] = viper.GetString("MLS.apikey")
  //var path = viper.GetString("new_api.path")
  var old_path = viper.GetString("old_api.path")
  var port = strconv.Itoa(viper.GetInt("old_api.port"))

  //TODO Multicore able
  http.HandleFunc(old_path, func(w http.ResponseWriter, r *http.Request) {
    get_loc_old(w, r, db)
  })
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
