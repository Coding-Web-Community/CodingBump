package main

import (
	"net/http"
)

func FetchGuilds(w http.
	ResponseWriter, r *http.Request) {
	guilds := gs.GetGuilds()

	if len(guilds) > 0 {
		WriteFetchResponse(w, 200, "Ok", guilds)
		return
	}

	WriteFetchResponse(w, 400, "BadRequest", guilds)
}
