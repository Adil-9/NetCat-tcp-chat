package server

import (
	"log"
	"os"
)

func LoggerCreate(logger *log.Logger) os.File {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		logger.Fatal(err)
	}
	// defer file.Close()

	logger.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	logger.SetOutput(file)
	logger.SetPrefix("INFO: 	")

	// for input := range channel {
	// 	log.Println(input)
	// }
	return *file
}
