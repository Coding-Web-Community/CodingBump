package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"

	"log"
)

const (
  URL             = "localhost"
	PORT            = ":8080"
	BUMP_INTERVAL   = 60 // 1 minute in seconds
	STORE_FILE_NAME = "store.json"
)

var (
	TempTestInterval = 0    // used to set lower interval during testing
	Logging          = true // used to disable logging during testing
	gs               GuildStore
)

type Guild struct {
	GuildId   int   `json:"guildId"`
	Timestamp int64 `json:"timestamp"`
}

type GuildStore struct {
	Guilds []Guild `json:"guilds"`
	mutex  sync.Mutex
}

type BumpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload Guild  `json:"payload"`
}

type FetchResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Payload []Guild `json:"paypload"`
}

func init() {
	var err error
	gs.Guilds, err = LoadStore()
	if err != nil {
		log.Print(err)
	} else {
		log.Print("'Database' successfully restored!")
		log.Print(fmt.Sprintf("%v guilds in 'database'", len(gs.Guilds)))
	}
}

func middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if Logging {
			log.Print(fmt.Sprintf("%s%s - %s", r.Host, r.URL.Path, r.Method))
		}
		f(w, r)
	}
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/V1/bump", middleware(BumpGuild)).Methods("POST")
	router.HandleFunc("/V1/fetch", middleware(FetchGuilds)).Methods("GET")

	log.Print(fmt.Sprintf("Now serving: %s%s",URL ,PORT))
	err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Print(err)
	}
}

// makes BumpResponse object and writes it to ResponseWriter
func WriteBumpResponse(w http.ResponseWriter, code int, message string, payload Guild) {
	br := BumpResponse{
		Code:    code,
		Message: message,
		Payload: payload,
	}

	payloadByte, _ := json.Marshal(br)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payloadByte)
}

func WriteFetchResponse(w http.ResponseWriter, code int, message string, payload []Guild) {
	fr := FetchResponse{
		Code:    code,
		Message: message,
		Payload: payload,
	}

	payloadByte, _ := json.Marshal(fr)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payloadByte)
}

func main() {
	HandleRequests()
}
