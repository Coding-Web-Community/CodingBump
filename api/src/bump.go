package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BumpGuild(w http.ResponseWriter, r *http.Request) {
	var guild Guild

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		WriteBumpResponse(w, http.StatusInternalServerError, "Unable to read request body", guild)
		return
	}

	// Unmarshal body into guild object
	err = json.Unmarshal(body, &guild)
	if err != nil {
		log.Print(err)
		if strings.Contains(err.Error(), "invalid character") {
			WriteBumpResponse(w, http.StatusBadRequest, "Request body contains invalid character", guild)
			return
		}
		WriteBumpResponse(w, http.StatusBadRequest, "Unable to process request body", guild)
		return
	}

	// check if GuildId's length is 18
	if len(strconv.Itoa(guild.GuildId)) != 18 {
		errString := "GuildId does not conform to 18 character long integer requirement"
		log.Print(errString)
		WriteBumpResponse(w, http.StatusBadRequest, errString, guild)
		return
	}

	// Adding timestamp to guild object
	ts := time.Now().Unix()
	guild.Timestamp = ts

	if !gs.GuildInStore(guild) {
		// guild is not yet present in GuildStore
		gs.AddToStore(guild)
		WriteBumpResponse(w, http.StatusOK, "Guild added and bumped", guild)
	} else {
		// guild is present in GuildStore
		if gs.PastInterval(guild) {
			// guild.Timestamp has exceeded BUMP_INTERVAL, Timestamp has been updated (bumped)
			WriteBumpResponse(w, http.StatusOK, "Guild bumped", guild)
		} else {
			// guild.Timestamp has not exceeded BUMP_INTERVAL, not updated
			guild.Timestamp = gs.GetTimestamp(guild)
			WriteBumpResponse(w, http.StatusTooEarly, "Guild bumped too early", guild)
		}
	}

	err = gs.WriteStore()
	if err != nil {
		log.Panic(err)
		WriteBumpResponse(w, http.StatusInternalServerError, "Could not store changes", guild)
		return
	}
}
