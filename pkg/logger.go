package pkg

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "message server - ", log.Lshortfile)
