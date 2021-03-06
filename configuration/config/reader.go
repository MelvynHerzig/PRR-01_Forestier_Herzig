package config

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

var localServerNumber uint

type Server struct {
	Ip     string `json:"ip"`
	Port   uint   `json:"port"`
	Parent uint   `json:"parent"`
}

type configReader struct {
	Debug    bool     `json:"debug"`
	NbRooms  uint     `json:"nbRooms"`
	NbNights uint     `json:"nbNights"`
	Servers  []Server `json:"servers"`
}

// Init Load the config file as a list of Server
// path is used to find the config file
// serverNumber is used on the server side, to remember who we are
func Init(path string, serverNumber uint) {
	rand.Seed(time.Now().UnixNano())
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &reader)

	localServerNumber = serverNumber
}

// InitSimple For the client, the second parameter is useless
func InitSimple(path string) {
	Init(path, 0)
}

// IsDebug Return the value of debug in config file
func IsDebug() bool {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.Debug
}

// GetRoomsCount Return the number of rooms specified in the config file
func GetRoomsCount() uint {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.NbRooms
}

// GetNightsCount Return the number of nights specified in the config file
func GetNightsCount() uint {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.NbNights
}

// GetServerById Return the server corresponding to the specified id
func GetServerById(id uint) *Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	if id >= uint(len(reader.Servers)) {
		return nil
	}

	return &reader.Servers[id]
}

// GetServerRandomly Return a server from the servers list
func GetServerRandomly() *Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return GetServerById(uint(rand.Intn(len(reader.Servers))))
}

// GetServers Return all servers from the list
func GetServers() []Server {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	return reader.Servers
}

// IsServerIP Check if the ip sent correspond to one of the server in the config file
func IsServerIP(address string) bool {
	if reader == nil {
		log.Fatal("config not initialized")
	}

	var ip = strings.Split(address, ":")[0]
	for _, server := range reader.Servers {
		if server.Ip == ip || (isLocalhost(server.Ip) && isLocalhost(ip)) {
			return true
		}
	}

	return false
}

// GetInitialChildrenIds Compute and return the children of specified Server id
func GetInitialChildrenIds(id uint) []uint {
	if reader == nil {
		log.Fatal("config not initialized")
	}
	var children []uint
	for index, child := range reader.Servers {
		if child.Parent == id && uint(index) != id{
			children = append(children, uint(index))
		}
	}

	return children
}

// isLocalhost Check if an adress is localhost
func isLocalhost(address string) bool {
	return address == "127.0.0.1" || address == "localhost"
}

// GetLocalServerNumber Get the id specified in the Init function
func GetLocalServerNumber() uint {
	return localServerNumber
}
