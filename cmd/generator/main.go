package main

import (
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/config"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/metadata"
)

func main() {
	config, err := config.NewConfigFromFile("./config/main.yaml")
	if err != nil {
		panic(err)
	}

	_, err = metadata.NewMetadataLoader(config)
	if err != nil {
		panic(err)
	}

	// decoder := yaml.NewDecoder(bytes.NewReader(data))

	// var docIndex int
	// for {
	// 	var raw interface{}
	// 	err := decoder.Decode(&raw)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		break
	// 	}
	// 	docIndex++
	// 	fmt.Printf("--- Document %d ---\n", docIndex)
	// 	fmt.Printf("%#v\n", raw)
	// }
}
