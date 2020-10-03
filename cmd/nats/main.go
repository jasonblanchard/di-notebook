package main

import (
	_ "github.com/lib/pq"
)

func main() {
	service, err := NewService()
	if err != nil {
		panic(err)
	}

	err = service.Run()

	if err != nil {
		panic(err)
	}
}
