package main

import (
	"fmt"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/config"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/discovery"
	discoveryInfra "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/infrastructure/discovery"
)

func main() {
	config, err := config.NewConfigFromFile("./config/main.yaml")
	if err != nil {
		panic(err)
	}

	// entityStore :=

	entityFactory, err := discovery.NewEntityFactory()
	if err != nil {
		panic(err)
	}

	filesystemDiscoveryConfig := discoveryInfra.FilesystemDiscoveryConfig{
		BaseDirectory:          config.DataDirectory,
		DirectoryToTypeMapping: config.DataDirectoryMappers,
	}
	filesystemDiscovery, err := discoveryInfra.NewFilesystemDiscovery(&filesystemDiscoveryConfig)
	if err != nil {
		panic(err)
	}

	entityFactory.RegisterEntitySource(filesystemDiscovery)

	entities, err := entityFactory.LoadEntities()
	if err != nil {
		panic(err)
	}

	for _, entity := range entities.GetEntities() {
		fmt.Println(fmt.Sprintf("Entity found: %s : %s", entity.GetName(), entity.GetType()))
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
