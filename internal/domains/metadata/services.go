package metadata

import (
	"fmt"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/config"
)

type MetadataLoader struct {
	config *config.Config
}

func NewMetadataLoader(config *config.Config) (*MetadataLoader, error) {
	if config == nil {
		return nil, fmt.Errorf("NewMetadataLoader passed with nil config")
	}
	return &MetadataLoader{
		config: config,
	}, nil
}

func (m *MetadataLoader) SearchForMetadataInDirectory() {

}
