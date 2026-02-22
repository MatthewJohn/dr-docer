package main

import (
	"bytes"
	"fmt"
	"os"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/metadata"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("config/servers/inthetz.md")
	if err != nil {
		panic(err)
	}

	metadataScanner, err := metadata.NewMetadataLoader(config)
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	var docIndex int
	for {
		var raw interface{}
		err := decoder.Decode(&raw)
		if err != nil {
			fmt.Println(err)
			break
		}
		docIndex++
		fmt.Printf("--- Document %d ---\n", docIndex)
		fmt.Printf("%#v\n", raw)
	}
}
