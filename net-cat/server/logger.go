package server

import (
	"log"
	"os"
)

func LoggerCreate() {
	channel := LogChannel
	file, err := os.Create("Logs.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(file)

	for input := range channel {
		log.Println(input)
	}
}
