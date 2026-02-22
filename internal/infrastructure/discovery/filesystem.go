package discovery

import (
	"fmt"

	discoveryDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/discovery"
	metadataDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/metadata"
)

type StorageMetadata struct {
}

type FilesystemEntityMetadata struct {
	Type metadataDomain.EntityType `yaml:"type"`
	Name string                    `yaml:"name"`
	// Criticality  string                    `yaml:"criticality"`
	// Host         string                    `yaml:"host"`
	// Storage      StorageMetadata           `yaml:"storage"`
	// Dependencies []string                  `yaml:"dependencies"`
	// Terraform    []string                  `yaml:"terraform"`
}

type FilesystemDiscoveryConfig struct {
	BaseDirectory          string
	DirectoryToTypeMapping map[string]string
}

type FilesystemDiscovery struct {
	config *FilesystemDiscoveryConfig
}

func NewFilesystemDiscovery(config *FilesystemDiscoveryConfig) (*FilesystemDiscovery, error) {
	if config == nil {
		return nil, fmt.Errorf("NewMetadataLoader passed with nil config")
	}
	return &FilesystemDiscovery{
		config: config,
	}, nil
}

func (m *FilesystemDiscovery) GetEntities(existingCollection *discoveryDomain.EntityCollection) error {
	return nil
}
func (m *FilesystemDiscovery) GetPriority() int {
	return 1
}

var _ discoveryDomain.EntitySource = &FilesystemDiscovery{}
