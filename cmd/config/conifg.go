package config

import "flag"

var FlagRunPort string
var FlagRunAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "b", "localhost", "address to run server")
	flag.StringVar(&FlagRunPort, "a", ":8080", "port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
