package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func startServer() {
	path := "store.json"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}

	go HandleRequests()
}

func sendInt(guildId int) (br BumpResponse) {
	reqBody, _ := json.Marshal(map[string]int{
		"guildId": guildId,
	})

	resp, _ := http.Post("http://localhost:8080/V1/bump", "application/json", bytes.NewBuffer(reqBody))

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &br)

	return br
}

func sendString(guildId string) (br BumpResponse) {
	reqBody, _ := json.Marshal(map[string]string{
		"guildId": guildId,
	})

	resp, _ := http.Post("http://localhost:8080/V1/bump", "application/json", bytes.NewBuffer(reqBody))

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &br)

	return br

}

func Hash(guilds []Guild) [16]byte {
	guildBytes := []byte{}
	for _, item := range guilds {
		jsonBytes, _ := json.Marshal(item)
		guildBytes = append(guildBytes, jsonBytes...)
	}
	return md5.Sum(guildBytes)
}

func TestServerStart(t *testing.T) {
	go startServer()

	TempTestInterval = 1
	Logging = false

	time.Sleep(time.Millisecond * 200)
}

func TestBumpNormal(t *testing.T) {
	var guildId int = 636145886279237611

	var expected = BumpResponse{
		Code: 200, // OK
		Payload: Guild{
			GuildId: guildId,
		},
	}

	// Normal send, returns 200
	resp := sendInt(guildId)

	if resp.Code != expected.Code {
		t.Errorf("Status codes were not equal: %v != %v", resp.Code, expected.Code)
	}
	if resp.Payload.GuildId != expected.Payload.GuildId {
		t.Errorf("GuilId's were not equal: %v != %v", resp.Payload.GuildId, expected.Payload.GuildId)
	}

}

func TestBumpEarly(t *testing.T) {
	var guildId int = 636145886279237699

	_ = sendInt(guildId) // send first bump request

	var expected = BumpResponse{
		Code: 425, // Too Early
		Payload: Guild{
			GuildId: guildId,
		},
	}

	resp := sendInt(guildId) // send second (early) bump request

	if resp.Code != expected.Code {
		t.Errorf("Status codes were not equal: %v != %v", resp.Code, expected.Code)
	}

	time.Sleep(time.Second * 2)

	expected = BumpResponse{
		Code: 200, // OK
		Payload: Guild{
			GuildId: guildId,
		},
	}

	resp = sendInt(guildId) // send third (late) bump request!

	if resp.Code != expected.Code {
		t.Errorf("Status codes were not equal: %v != %v", resp.Code, expected.Code)
	}
	if resp.Payload.GuildId != expected.Payload.GuildId {
		t.Errorf("GuilId's were not equal: %v != %v", resp.Payload.GuildId, expected.Payload.GuildId)
	}

}

func TestBumpString(t *testing.T) {
	var expected = BumpResponse{
		Code:    400, // Bad Request
		Message: "Unable to process request body",
		Payload: Guild{
			GuildId: 0,
		},
	}

	resp := sendString("636145886279237652")

	if resp.Code != expected.Code {
		t.Errorf("Status codes were not equal: %v != %v", resp.Code, expected.Code)
	}
	if resp.Payload.GuildId != expected.Payload.GuildId {
		t.Errorf("GuilId's were not equal: %v != %v", resp.Payload.GuildId, expected.Payload.GuildId)
	}
	if resp.Message != expected.Message {
		t.Errorf("Messages were not equal: %v != %v", resp.Message, expected.Message)
	}
}

func TestBumpTooFewChars(t *testing.T) {
	var guildId int = 636145
	var expected = BumpResponse{
		Code:    400, // Bad Request
		Message: "GuildId does not conform to 18 character long integer requirement",
		Payload: Guild{
			GuildId: guildId,
		},
	}

	resp := sendInt(guildId)

	if resp.Code != expected.Code {
		t.Errorf("Status codes were not equal: %v != %v", resp.Code, expected.Code)
	}
	if resp.Payload.GuildId != expected.Payload.GuildId {
		t.Errorf("GuilId's were not equal: %v != %v", resp.Payload.GuildId, expected.Payload.GuildId)
	}
	if resp.Message != expected.Message {
		t.Errorf("Messages were not equal: %v != %v", resp.Message, expected.Message)
	}
}

func TestFetch(t *testing.T) {
	var fr FetchResponse
	guilds, _ := LoadStore()

	resp, _ := http.Get("http://localhost:8080/V1/fetch")

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &fr)

	if len(fr.Payload) != len(guilds) {
		t.Errorf("Length of stored guilds not equal to length of /V1/fetch result: %v != %v", len(fr.Payload), len(guilds))
	}

	if Hash(fr.Payload) != Hash(guilds) {
		t.Errorf("Guilds hash from GuildStore are not equal to /V1/fetch hash: %v != %v", Hash(fr.Payload), Hash(guilds))
	}

}
