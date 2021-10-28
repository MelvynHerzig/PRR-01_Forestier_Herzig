package config

import "strings"

// TODO remove mocking of config
type Server struct {
	Ip   string
	Port string
}

var Servers = [3]Server{ {Ip: "localhost", Port:"10000"},
						 {Ip: "localhost", Port:"10001"},
						 {Ip: "localhost", Port:"10002"}}

func GetRoomsCount() uint {
	return 10
}

func GetDaysCount() uint {
	return 10
}

func IsServerIP(address string) bool {

	var ip = strings.Split(address, ":")[0]
	for _, server := range Servers {
		if server.Ip == ip || (server.Ip == "localhost" && ip == "127.0.0.1") {
			return true
		}
	}

	return false
}
