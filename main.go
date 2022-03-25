package main

import (
	"fmt"
	"os"

	grt "github.com/noris-network/kustomize-generalreplacementstransformer/internal"
)

var build = "dev"

func main() {

	// check args
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Printf("GeneralReplacementsTransformer %v", build)
		os.Exit(0)
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: GeneralReplacementsTransformer <configfile>|version")
		os.Exit(1)
	}

	// new transformer
	tx, err := grt.New(grt.WithConfigFile(os.Args[1]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tx.ReadStream(os.Stdin); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tx.ScanForValues(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tx.WriteStream(os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
