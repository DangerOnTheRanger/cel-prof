package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	argsFile = "cel-args.json"
)

type args struct {
	List   []string `json:"l"`
	String string   `json:"s"`
}

func main() {
	numBytesStr := os.Args[1]
	numBytes, err := strconv.Atoi(numBytesStr)
	if err != nil {
		log.Fatal(err)
	}
	l := make([]string, numBytes)
	for i, _ := range l {
		l[i] = "0"
	}
	argsStr := os.Args[2]
	jsonArgs := args{l, argsStr}

	bytes, err := json.Marshal(jsonArgs)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(argsFile, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
