package main

import (
	"flag"

	_ "github.com/lib/pq"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "Config file")
	flag.Parse()

	err := initConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	service, err := NewServiceFromEnv()
	if err != nil {
		panic(err)
	}

	err = service.Run()

	if err != nil {
		panic(err)
	}
}
