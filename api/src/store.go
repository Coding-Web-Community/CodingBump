package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// loads the guildstore from store.json
// creates a new file it it doesn't edist
func LoadStore() (guilds []Guild, err error) {
	// Open store.json file
	file, err := os.OpenFile(STORE_FILE_NAME, os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		return guilds, err
	}

	// Read file into contents
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return guilds, err
	}

	if len(contents) == 0 {
		// Store is empty, no need to unmarshal
		return guilds, err
	}

	err = json.Unmarshal(contents, &guilds)
	return guilds, err
}

// writes the guildstore to store.json
func (gs *GuildStore) WriteStore() (err error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.Guilds = sortGuilds(gs.Guilds)

	// marshalled guilds
	m, err := json.Marshal(gs.Guilds)
	if err != nil {
		return err
	}

	// Write marshalled guildStore.guilds to store.json
	err = ioutil.WriteFile(STORE_FILE_NAME, m, 0644)
	return err
}

// sorts the guilds based on their timestamp
func sortGuilds(guilds []Guild) []Guild {
	sort.Slice(guilds, func(i, j int) bool {
		return guilds[i].Timestamp > guilds[j].Timestamp
	})

	return guilds
}

// checks whether or not a guild is in the guildstore
func (gs *GuildStore) GuildInStore(guild Guild) bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for _, gsGuild := range gs.Guilds {
		if gsGuild.GuildId == guild.GuildId {
			// Guild in GuildStore
			return true
		}
	}

	// Guild not in GuildStore
	return false
}

// adds a guild to the guild store
func (gs *GuildStore) AddToStore(guild Guild) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.Guilds = append(gs.Guilds, guild)

}

// checks if a guild has passed the guild bumping interval
func (gs *GuildStore) PastInterval(guild Guild) bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for i, gsGuild := range gs.Guilds {
		if gsGuild.GuildId == guild.GuildId {

			ts := time.Now().Unix()
			var interval int64
			if TempTestInterval == 0 {
				interval = BUMP_INTERVAL
			} else {
				interval = int64(TempTestInterval)
			}

			if (ts - gsGuild.Timestamp) >= interval {
				gs.Guilds[i].Timestamp = ts

				return true
			} else {
				return false
			}
		}
	}

	return false
}

// returns the timestamp in the guildstore for the guild
func (gs *GuildStore) GetTimestamp(guild Guild) int64 {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for _, gsGuild := range gs.Guilds {
		if gsGuild.GuildId == guild.GuildId {
			// found guild, now returning Timestamp
			return gsGuild.Timestamp
		}
	}

	return 0
}
