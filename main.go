package main

import (
	"fmt"
	"log"
	"os"

	grt "github.com/noris-network/kustomize-generalreplacementstransformer/internal"
)

var build = "dev"

func main() {

	// check args
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println("GeneralReplacementsTransformer", build)
		os.Exit(0)
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: GeneralReplacementsTransformer <configfile>|version")
		os.Exit(1)
	}

	// setup logging
	log.SetPrefix("# GeneralReplacementsTransformer: ")
	log.SetFlags(0)

	if logDest := os.Getenv("GRT_LOG"); logDest != "" {
		file, err := os.Create(logDest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		log.SetOutput(file)
		defer file.Close()
	}

	// new transformer
	tx, err := grt.New(grt.WithConfigFile(os.Args[1]))
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := tx.ReadStream(os.Stdin); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := tx.ScanForValues(); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := tx.WriteStream(os.Stdout); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

}
