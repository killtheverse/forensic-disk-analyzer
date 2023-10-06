package main

import (
	"flag"
	"fmt"
)

func main() {

	filepath := flag.String("f", "", "")
	flag.Parse()

	//fmt.Println("filepath: ", *filepath)
	
	err := StoreHashes(*filepath)
	if err != nil {
		fmt.Println(err)
	}

	err = AnalyzeImage(*filepath)
	if err != nil {
		fmt.Println(err)
	}
}
