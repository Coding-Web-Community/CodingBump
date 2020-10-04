package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
  "sort"
  "time"
)

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

func (gs *GuildStore) WriteStore() (err error) {
  gs.mutex.Lock()
  defer gs.mutex.Unlock()

  gs.Guilds = sortGuilds(gs.Guilds)

  g, err := json.Marshal(gs.Guilds)
  if err != nil {
    return err
  }

  // Write marshalled guildStore.guilds to store.json
  err = ioutil.WriteFile(STORE_FILE_NAME, g, 0644)
  return err
}

func sortGuilds(guilds []Guild) []Guild {
  sort.Slice(guilds, func(i, j int) bool {
    return guilds[i].Timestamp > guilds[j].Timestamp
  })

  return guilds
}

func (gs *GuildStore) GuildInStore(guild Guild) (bool) {
  gs.mutex.Lock()
  defer gs.mutex.Unlock()

  for _, storedGuild := range gs.Guilds {
    if storedGuild.GuildId == guild.GuildId {
      return true
    }
  }

  // Guild not in gs.Guilds
  return false
}

func (gs *GuildStore) AddToStore(guild Guild) {
  gs.mutex.Lock()
  defer gs.mutex.Unlock()

  gs.Guilds = append(gs.Guilds, guild)

}

func (gs *GuildStore) PastInterval(guild Guild) bool {
  gs.mutex.Lock()
  defer gs.mutex.Unlock()

  for i, storedGuild := range gs.Guilds {
    if storedGuild.GuildId == guild.GuildId {
      ts := time.Now().Unix()
      if (ts - storedGuild.Timestamp) >= BUMP_INTERVAL {
        gs.Guilds[i].Timestamp = ts
        return true
      } else {
        return false
      }
    }
  }

  return false
}

func (gs *GuildStore) GetTimestamp(guild Guild) int64 {
  gs.mutex.Lock()
  defer gs.mutex.Unlock()

  for _, storedGuild := range gs.Guilds {
    if storedGuild.GuildId == guild.GuildId {
      return storedGuild.Timestamp
    }
  }

  return 0
}
