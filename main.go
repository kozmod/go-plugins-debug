package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
)

func main() {

	plugNum := "english"
	if len(os.Args) == 2 {
		plugNum = os.Args[1]
	}
	var mod string
	var f string
	switch plugNum {
	case "1":
		mod = "plug1.so"
		f = "GreetUniverse"
	case "2":
		mod = "plug2.so"
		f = "GreetWorld"
	default:
		log.Println("plugin did not chose")
		os.Exit(1)
	}

	log.Println("try open plugin:", plugNum, f)

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	path := filepath.Join(cwd, mod)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println(mod, err)
		os.Exit(0)
	}

	plug, err := plugin.Open(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	sym, err := plug.Lookup(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	greeter, ok := sym.(func() string)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		res := greeter()
		_, err := writer.Write([]byte(res))
		if err != nil {
			log.Fatal(err)
		}
	})

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
