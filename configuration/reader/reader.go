package reader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var reader *configReader = nil

type Server struct {
	Ip       string `json:"ip"`
	Port     uint   `json:"port"`
}

type configReader struct{
	NbRooms uint `json:"nbRooms"`
	NbNights uint `json:"nbNights"`
	Servers []Server `json:"servers"`
}

func Init(path string) error {
	rand.Seed(time.Now().UnixNano())
	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)


	json.Unmarshal(byteValue, &reader)

	return nil
}

func GetRoomsCount() (uint, error){
	if reader == nil{
		return 0, errors.New("reader not initialized")
	}

	return reader.NbRooms, nil
}

func GetNightsCount() (uint, error){
	if reader == nil{
		return 0, errors.New("reader not initialized")
	}

	return reader.NbNights, nil
}

func GetServerById(id uint) (*Server, error) {
	if reader == nil{
		return nil, errors.New("reader not initialized")
	}

	if id >= uint(len(reader.Servers)) {
		return nil, errors.New("invalid id")
	}
	return &reader.Servers[id], nil
}

func GetServerRandomly() (*Server, error)    {
	if reader == nil{
		return nil, errors.New("reader not initialized")
	}

	return GetServerById(uint(rand.Intn(len(reader.Servers))))
}


func GetServers() ([]Server, error) {
	if reader == nil{
		return nil, errors.New("reader not initialized")
	}

	return reader.Servers, nil
}

func IsServerIP(address string) (bool, error) {
	if reader == nil{
		return false, errors.New("reader not initialized")
	}

	var ip = strings.Split(address, ":")[0]
	for _, server := range reader.Servers {
		if server.Ip == ip || (isLocalhost(server.Ip) && isLocalhost(ip)) {
			return true, nil
		}
	}

	return false, nil
}

func isLocalhost(address string) bool {
	return address == "127.0.0.1" || address == "localhost"
}

