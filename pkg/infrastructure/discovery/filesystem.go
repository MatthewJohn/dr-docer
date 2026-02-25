package discovery

import (
	"bytes"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/attribute"
	commontypes "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/common_types"
	discoveryDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/discovery"
	metadataDomain "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
	"go.yaml.in/yaml/v3"
)

// type StorageMetadata struct {
// }

type FilesystemEntityMetadata struct {
	Type      metadataDomain.EntityType `yaml:"type"`
	Name      string                    `yaml:"name"`
	IpAddress string                    `yaml:"ip_address"`
	Url       string                    `yaml:"url"`
	// Criticality  string                    `yaml:"criticality"`
	// Host         string                    `yaml:"host"`
	// Storage      StorageMetadata           `yaml:"storage"`
	// Dependencies []string                  `yaml:"dependencies"`
	// Terraform    []string                  `yaml:"terraform"`
}

type FilesystemDiscoveryConfig struct {
	BaseDirectory          string
	DirectoryToTypeMapping map[string]string
	FileExtensions         []string
}

type FilesystemDiscovery struct {
	config *FilesystemDiscoveryConfig
}

func (m *FilesystemDiscovery) GetEntityTypes() []metadataDomain.EntityType {
	return []metadataDomain.EntityType{
		commontypes.EntityServer,
		commontypes.EntityService,
	}
}

func (m *FilesystemDiscovery) GetAttributes() []attribute.Attribute {
	return []attribute.Attribute{
		commontypes.AttributeIpAddress,
		commontypes.AttributeUrl,
	}
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		// Path does not exist or cannot be accessed
		return false, err
	}
	return info.IsDir(), nil
}

func findFiles(rootDir string, extensions []string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(rootDir, func(name string, dirEntry fs.DirEntry, err error) error {
		fmt.Printf("Processing file %s\n", name)
		if err != nil {
			return err
		}
		for _, extension := range extensions {
			if filepath.Ext(dirEntry.Name()) == extension {
				files = append(files, name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func NewFilesystemDiscovery(config *FilesystemDiscoveryConfig) (*FilesystemDiscovery, error) {
	if config == nil {
		return nil, fmt.Errorf("NewFilesystemDiscovery passed with nil config")
	}

	if isDir, err := isDirectory(config.BaseDirectory); err != nil || !isDir {
		if err != nil {
			return nil, fmt.Errorf("NewFilesystemDiscovery: Failed to check if data directory is valid: %s", err)
		}
		return nil, fmt.Errorf("NewFilesystemDiscovery: BaseDirectory does not exist or is not a directory")
	}

	return &FilesystemDiscovery{
		config: config,
	}, nil
}

func (m *FilesystemDiscovery) convertFilepathToType(filePath string) metadataDomain.EntityType {
	for dirMatch := range maps.Keys(m.config.DirectoryToTypeMapping) {
		for _, pathPart := range strings.Split(filePath, string(os.PathSeparator)) {
			if pathPart == dirMatch {
				return metadataDomain.EntityType(m.config.DirectoryToTypeMapping[dirMatch])
			}
		}
	}
	return ""
}

func (m *FilesystemDiscovery) processRawFilesystemMetadata(raw *FilesystemEntityMetadata, collection *discoveryDomain.EntityCollection, filePath string) error {
	// Attempt to extract type from path
	if raw.Type == "" {
		if pathType := m.convertFilepathToType(filePath); pathType != "" {
			raw.Type = pathType
		}
	}
	if raw.Type == "" {
		return fmt.Errorf("Empty type found")
	}

	if raw.Name == "" {
		return fmt.Errorf("Empty enitty name")
	}

	entity, err := metadataDomain.NewEntity(metadataDomain.EntityName(raw.Name), raw.Type, m.GetPriority())
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", raw.Type)
	switch entity.Type {
	case commontypes.EntityServer:
		fmt.Printf("Processing Server entity\n")
		entity.SetAttribute(&commontypes.AttributeIpAddress, raw.IpAddress)

	case commontypes.EntityService:
		fmt.Printf("Processing Service entity\n")
		entity.SetAttribute(&commontypes.AttributeUrl, raw.Url)

	default:
		return fmt.Errorf("Unknown entity type: %s\n", raw.Type)
	}
	fmt.Printf("Entity: %#v\n", entity)
	if entity != nil {
		err := collection.AddEntity(entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *FilesystemDiscovery) processFile(collection *discoveryDomain.EntityCollection, filePath string) error {
	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Convert to entity
	decoder := yaml.NewDecoder(bytes.NewReader(fileData))

	for {
		var raw FilesystemEntityMetadata
		err := decoder.Decode(&raw)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println(err)
			continue
		}
		fmt.Printf("--- Document ---\n")
		fmt.Printf("%#v\n", raw)
		err = m.processRawFilesystemMetadata(&raw, collection, filePath)
		if err != nil {
			fmt.Printf("processFile: Error processing file fragment: %s\n", err)
		}
	}
	return nil
}

func (m *FilesystemDiscovery) GetEntities(collection *discoveryDomain.EntityCollection) error {
	filePaths, err := findFiles(m.config.BaseDirectory, m.config.FileExtensions)
	if err != nil {
		return err
	}

	for _, filePath := range filePaths {
		m.processFile(collection, filePath)
	}
	return nil
}
func (m *FilesystemDiscovery) GetPriority() int {
	return 50
}

var _ discoveryDomain.EntitySource = &FilesystemDiscovery{}
