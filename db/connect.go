package db

import (
	"log"
)

var DB = NewClient()

func Connect() {
	if err := DB.Prisma.Connect(); err != nil {
		log.Fatalln("Error connecting")
	}
}

func Disconnect() {
	DB.Prisma.Disconnect()
}
