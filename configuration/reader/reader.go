package reader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var reader *configReader = nil

type Server struct {
	Ip   string `json:"ip"`
	Port uint   `json:"port"`
}

type configReader struct {
	Debug    bool     `json:"debug"`
	NbRooms  uint     `json:"nbRooms"`
	NbNights uint     `json:"nbNights"`
	Servers  []Server `json:"servers"`
}

func Init(path string) {
	rand.Seed(time.Now().UnixNano())
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &reader)
}

func IsDebug() bool {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	return reader.Debug
}

func GetRoomsCount() uint {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	return reader.NbRooms
}

func GetNightsCount() uint {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	return reader.NbNights
}

func GetServerById(id uint) *Server {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	if id >= uint(len(reader.Servers)) {
		return nil
	}

	return &reader.Servers[id]
}

func GetServerRandomly() *Server {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	return GetServerById(uint(rand.Intn(len(reader.Servers))))
}

func GetServers() []Server {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	return reader.Servers
}

func IsServerIP(address string) bool {
	if reader == nil {
		log.Fatal("reader not initialized")
	}

	var ip = strings.Split(address, ":")[0]
	for _, server := range reader.Servers {
		if server.Ip == ip || (isLocalhost(server.Ip) && isLocalhost(ip)) {
			return true
		}
	}

	return false
}

func isLocalhost(address string) bool {
	return address == "127.0.0.1" || address == "localhost"
}
