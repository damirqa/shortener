package config

import "flag"

var Config = struct {
	Address string
	Port    string
}{}

func Init() {
	flag.StringVar(&Config.Address, "b", "localhost", "address to run server")
	flag.StringVar(&Config.Port, "a", "8080", "port to run server")
	flag.Parse()
}
