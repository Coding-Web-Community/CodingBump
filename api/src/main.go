package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"log"
)

const (
	PORT = ":8080"
  BUMP_INTERVAL = 60 // 1 minute in seconds
	STORE_FILE_NAME = "store.json"
)

var gs GuildStore

type BumpRaw struct {
	GuildId string `json:"guildId"`
}

type Guild struct {
	GuildId   int   `json:"guildId"`
	Timestamp int64 `json:"timestamp"`
}

type GuildStore struct {
	Guilds []Guild `json:"guilds"`
	mutex sync.Mutex
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
		log.Print(fmt.Sprintf("%s%s - %s", r.Host, r.URL.Path, r.Method))
		f(w, r)
	}
}


func BumpGuild(w http.ResponseWriter, r *http.Request) {
  // Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - InternalServerError"))
		return
	}

	// Unmarshal body into BumpRaw
	var br BumpRaw
	err = json.Unmarshal(body, &br)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - InternalServerError"))
		return
	}

	// Convert BumpRaw.Guild_id into int64 from string
	guildId, err := strconv.Atoi(br.GuildId)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - BadRequest"))
		return
	}

	// Constructing actual Bump object
	ts := time.Now().Unix()
	var guild = Guild{
		GuildId:   guildId,
		Timestamp: ts,
	}

  m, err := json.Marshal(guild)
  if err != nil {
    log.Print(err)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 - InternalServerError"))
    return
  }

  if !gs.GuildInStore(guild) {
    // guild is not yet present in GuildStore
    gs.AddToStore(guild)
    w.WriteHeader(http.StatusOK)
    w.Write(m)
  } else {

    if gs.PastInterval(guild) {
      // guild.Timestamp has exceeded BUMP_INTERVAL, Timestamp has been updated (bumped)
      w.WriteHeader(http.StatusOK)
      w.Write(m)
    } else {
      // guild.Timestamp has not exceeded BUMP_INTERVAL, not updated
      guild.Timestamp = gs.GetTimestamp(guild)

      // update m with new timestamp
      m, err := json.Marshal(guild)
      if err != nil {
        log.Print(err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500 - InternalServerError"))
        return
      }

      w.WriteHeader(http.StatusTooEarly)
      w.Write(m)
    }
  }

  err = gs.WriteStore()
  if err != nil {
    log.Panic(err)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 - InternalServerError"))
    return
  }
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/V1/bump", middleware(BumpGuild)).Methods("POST")
  log.Print(fmt.Sprintf("Now serving: localhost%s", PORT))
  err := http.ListenAndServe(PORT, router)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	HandleRequests()
}
