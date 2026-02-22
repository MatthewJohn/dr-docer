package discovery

import (
	"fmt"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/config"
	discoveryDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/discovery"
	metadataDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/metadata"
)

type StorageMetadata struct {
}

type FilesystemEntityMetadata struct {
	Type         metadataDomain.EntityType `yaml:"type"`
	Name         string                    `yaml:"name"`
	Criticality  string                    `yaml:"criticality"`
	Host         string                    `yaml:"host"`
	Storage      StorageMetadata           `yaml:"storage"`
	Dependencies []string                  `yaml:"dependencies"`
	Terraform    []string                  `yaml:"terraform"`
}

type FilesystemDiscovery struct {
	config *config.Config
}

func NewFilesystemDiscovery(config *config.Config) (*FilesystemDiscovery, error) {
	if config == nil {
		return nil, fmt.Errorf("NewMetadataLoader passed with nil config")
	}
	return &FilesystemDiscovery{
		config: config,
	}, nil
}

func (m FilesystemDiscovery) GetEntities(existingCollection *discoveryDomain.EntityCollection) error {
	return nil
}
func (m FilesystemDiscovery) GetPriority() int {
	return 1
}

var _ discoveryDomain.EntitySource = FilesystemDiscovery{}
