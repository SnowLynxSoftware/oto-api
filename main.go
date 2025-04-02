package main

import (
	"fmt"
	"log"

	"github.com/sonwlynxsoftware/oto-api/cmd"
)

func main() {
	handler := cmd.NewHandler()
	err := handler.ExecuteCommand()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
