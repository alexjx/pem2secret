package main

import (
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var name = flag.String("name", "", "the name of the secret")
var namespace = flag.String("namespace", "", "the namespace of the secret")

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("ERR: missing arguments")
		os.Exit(1)
	}
	if *name == "" {
		fmt.Println("ERR: missing pem files")
		os.Exit(1)
	}

	secretItems := make(map[string]string)
	for _, fname := range flag.Args() {
		if fp, err := os.Open(fname); err != nil {
			fmt.Printf("ERR: %s: %s", fname, err.Error())
		} else if content, err := ioutil.ReadAll(fp); err != nil {
			fmt.Printf("ERR: %s read fail: %s", fname, err.Error())
		} else if block, _ := pem.Decode(content); block == nil {
			fmt.Printf("ERR: %s is not a pem file\n", fname)
		} else {
			secretItems[fname] = string(content)
		}
	}

	if len(secretItems) > 0 {
		fmt.Println("--------------------")
		fmt.Printf("apiVersion: v1\n")
		fmt.Printf("kind: Secret\n")
		fmt.Printf("metadata:\n")
		if *namespace != "" {
			fmt.Printf("  namespace: %s\n", *namespace)
		}
		fmt.Printf("  name: %s\n", *name)
		fmt.Printf("data:\n")
		for fname, content := range secretItems {
			fmt.Printf("  %s: |\n", fname)
			for _, line := range strings.Split(content, "\n") {
				fmt.Printf("    %s\n", line)
			}
		}
	}
}
