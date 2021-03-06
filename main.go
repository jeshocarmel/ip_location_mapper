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

func main() {
	setupRoutes()
	connectToRedis()
	log.Info("Go Web App Started on Port 8080")
	http.ListenAndServe(":8080", nil)
}

func setupRoutes() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/getlocation", getLocation)
	http.HandleFunc("/getmylocation", getMyLocation)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
}

func connectToRedis() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err == nil {
		log.Info(pong, err)
	} else {
		log.Error(err)
	}
}

func createFuncMap() template.FuncMap {

	// funcmap for go templates
	funcMap := template.FuncMap{

		/*

			"getTimeAgo": func(t time.Time) string {
				return timeago.English.Format(t)
			},
			"getTimeAgoForMillis": func(tUnix int64) string {
				return getTimeAgoForMillis(tUnix)
			},

		*/
	}

	return funcMap
}

/*

func getTimeAgoForMillis(tUnixNano int64) string {
	t := time.Unix(0, tUnixNano)
	return timeago.English.Format(t)
}

*/

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

		apiResponse, err := makeAPIRequest(ipaddress)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		tmp, _ := json.Marshal(apiResponse)
		w.Write(tmp)

	case "GET":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getMyLocation(w http.ResponseWriter, r *http.Request) {

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

		ipAddr := getIP(r)

		apiResponse, err := makeAPIRequest(ipAddr)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		tmp, _ := json.Marshal(apiResponse)
		w.Write(tmp)

	case "GET":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getIP(r *http.Request) string {

	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return IPAddress
}

func makeAPIRequest(ipaddress string) (*IPStackResponseSuccess, error) {

	//check if ipaddress available in cache before calling the API
	cacheResp, err := getFromCache(ipaddress)
	if err == nil && cacheResp != nil {
		log.Infof("response obtained from cache for ipaddress: %s", ipaddress)
		return cacheResp, nil
	}

	apiURL := fmt.Sprintf("%s%s?access_key=%s", defaultIPStackURL, ipaddress, os.Getenv("IPSTACK_API_KEY"))
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response obtained from ipstack api")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipStackFail IPStackResponseError
	err = json.Unmarshal(body, &ipStackFail)
	if err != nil {
		return nil, err
	}

	//check if ipStackFail struct empty i.e. making sure the response from ipstack is not a failure response
	if ipStackFail == (IPStackResponseError{}) {
		var ipStackResponseSuccess IPStackResponseSuccess
		err = json.Unmarshal(body, &ipStackResponseSuccess)
		if err != nil {
			return nil, err
		}

		//store in redis
		err = rdb.Set(ctx, ipaddress, &ipStackResponseSuccess, time.Hour*24).Err()
		if err != nil {
			log.Error(err)
		}
		return &ipStackResponseSuccess, nil

	}
	return nil, errors.New(ipStackFail.Error.Info)
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

//MarshalBinary ...
func (obj *IPStackResponseSuccess) MarshalBinary() ([]byte, error) {
	return json.Marshal(obj)
}

// UnmarshalBinary -
func (obj *IPStackResponseSuccess) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	return nil
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

func validateIP(ip string) error {

	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	return nil
}

func getFromCache(ipaddress string) (*IPStackResponseSuccess, error) {

	cacheData, err := rdb.Get(ctx, ipaddress).Result()
	if err != nil {
		return nil, err
	}

	var obj IPStackResponseSuccess
	err = obj.UnmarshalBinary([]byte(cacheData))
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
