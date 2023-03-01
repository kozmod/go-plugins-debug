package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
)

func main() {

	plugNum := "1"
	if len(os.Args) == 2 {
		plugNum = os.Args[1]
	}
	var mod string
	var f string
	switch plugNum {
	case "1":
		mod = "out/plug1.so"
		f = "GreetUniverse"
	case "2":
		mod = "out/plug2.so"
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

	cmd := exec.Command("ls", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal("cmd:", err)
	}
	log.Println(out.String())

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

	res := greeter()
	log.Println("greeter: ", res)
}
