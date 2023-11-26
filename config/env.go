package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Type Vars
//
// Environment variables
type Vars struct {
	SECRET    string
	MONGO_URI string
	PORT      string
}

func (v Vars) Validate() {
	switch {
	case v.SECRET == "":
		log.Fatal("missing SECRET env")
	case v.MONGO_URI == "":
		log.Fatal("missing MONGO_URI env")
	}
}

// Env() returns Vars struct of environment variables
func Env() Vars {
	var v Vars
	// Load if not a test. This isn't required during testing.
	if flag.Lookup("test.v") == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		v = Vars{
			SECRET:    os.Getenv("SECRET"),
			MONGO_URI: os.Getenv("MONGO_URI"),
		}

		switch {
		case os.Getenv("PORT") == "":
			v.PORT = ":3000"
		default:
			v.PORT = ":" + os.Getenv("PORT")
		}
		v.Validate()
	}

	return v
}
