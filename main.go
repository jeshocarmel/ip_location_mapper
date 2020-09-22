package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"text/template"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/xeonx/timeago"
)

var templates *template.Template
var rdb *redis.Client
var ctx = context.Background()

const (
	defaultIPStackURL = "http://api.ipstack.com/"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	funcMap := createFuncMap()
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.gohtml"))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.gohtml", nil)
}

func getLocation(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		// to handle panics
		defer func() {
			msg := recover()
			if msg != nil { //catch
				errMsg := fmt.Sprintf(msg.(error).Error())
				myMap := make(map[string]interface{})
				myMap["error"] = errMsg
				tmp, _ := json.Marshal(myMap)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(tmp)
			}
		}()

		r.ParseForm()
		ipaddress := r.FormValue("ipaddress")
		if err := validateIP(ipaddress); err != nil {
			panic(err)
		}

		// apiURL := fmt.Sprintf("%s%s?access_key=%s", defaultIPStackURL, ipaddress, "f31c198f0daa774d2b2604243676cc34")
		apiURL := fmt.Sprintf("%s%s?access_key=%s", defaultIPStackURL, ipaddress, os.Getenv("IPSTACK_API_KEY"))
		resp, err := http.Get(apiURL)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			panic(errors.New("invalid response obtained from ipstack api"))
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var ipStackFail IPStackResponseError
		err = json.Unmarshal(body, &ipStackFail)
		if err != nil {
			panic(err)
		}

		//check if ipStackFail struct empty i.e. making sure the response from ipstack is not a failure response
		if ipStackFail == (IPStackResponseError{}) {
			var ipStackResponseSuccess IPStackResponseSuccess
			err = json.Unmarshal(body, &ipStackResponseSuccess)
			if err != nil {
				panic(err)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				tmp, _ := json.Marshal(ipStackResponseSuccess)
				w.Write(tmp)
			}
		} else {
			panic(errors.New(ipStackFail.Error.Info))
		}

	case "GET":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func validateIP(ip string) error {

	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	return nil
}

func setupRoutes() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/getlocation", getLocation)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
}

func main() {
	log.Info("Go Web App Started on Port 8080")
	setupRoutes()
	connectToRedis()
	http.ListenAndServe(":8080", nil)
}

func connectToRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("redis-password"), // no password set
		DB:       0,                           // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err == nil {
		log.Info(pong, err)
	} else {
		log.Error(err)
	}
	// Output: PONG <nil>
}

func createFuncMap() template.FuncMap {

	// funcmap for go templates
	funcMap := template.FuncMap{

		"getTimeAgo": func(t time.Time) string {
			return timeago.English.Format(t)
		},
		"getTimeAgoForMillis": func(tUnix int64) string {
			return getTimeAgoForMillis(tUnix)
		},
	}

	return funcMap
}

func getTimeAgoForMillis(tUnixNano int64) string {
	t := time.Unix(0, tUnixNano)
	return timeago.English.Format(t)
}

//IPStackResponseSuccess is the response object from ipstack
type IPStackResponseSuccess struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
}

//IPStackResponseError ...
type IPStackResponseError struct {
	Success bool `json:"success"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}
