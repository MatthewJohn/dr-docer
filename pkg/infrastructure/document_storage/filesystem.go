package documentstorage

import (
	"fmt"
	"os"
	"path"

	documentgenerator "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/document_generator"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
)

type DocumentStorageFileConfig struct {
	OutputDirectory string
}

type DocumentStorageFile struct {
	config *DocumentStorageFileConfig
}

func NewDocumentStorageFile(config *DocumentStorageFileConfig) (*DocumentStorageFile, error) {
	if config == nil {
		return nil, fmt.Errorf("NewDocumentStorageFile: config is nil")
	}

	if stat, err := os.Stat(config.OutputDirectory); err == nil && !stat.IsDir() {
		return nil, fmt.Errorf("NewDocumentStorageFile: OutputDirectory is not a directory")
	}

	return &DocumentStorageFile{
		config: config,
	}, nil
}

func (d DocumentStorageFile) StoreDocument(entityName metadata.EntityName, entityType metadata.EntityType, document []byte) error {
	typeDir := path.Join(d.config.OutputDirectory, string(entityType))
	if err := os.MkdirAll(typeDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path.Join(typeDir, fmt.Sprintf("%s.md", string(entityName))), document, 0o644)
}

var _ documentgenerator.DocumentStorage = DocumentStorageFile{}
