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

	fmt.Printf("Config: %#v", config)
	filesystemDiscoveryConfig := discoveryInfra.FilesystemDiscoveryConfig{
		BaseDirectory:          config.DataDirectory,
		DirectoryToTypeMapping: config.DataDirectoryMappers,
		FileExtensions:         config.DataFileExtensions,
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
		fmt.Printf("Entity found: %s : %s : %#v\n", entity.GetName(), entity.GetType(), entity)
	}
}
