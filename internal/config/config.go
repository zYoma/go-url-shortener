package config

import (
	"flag"
)

var flagRunAddr string
var flagBaseShortURL string

func ParseFlags() (string, string) {

	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagBaseShortURL, "b", "http://localhost:8080", "base short url")
	flag.Parse()

	return flagRunAddr, flagBaseShortURL
}
