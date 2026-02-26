package documentstorage

import (
	"fmt"

	documentgenerator "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/document_generator"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
)

type DocumentStorageStdout struct{}

func NewDocumentStorageStdout() (*DocumentStorageStdout, error) {
	return &DocumentStorageStdout{}, nil
}

func (d DocumentStorageStdout) StoreDocument(entityName metadata.EntityName, entityType metadata.EntityType, document []byte) error {
	fmt.Println(string(document))
	return nil
}

var _ documentgenerator.DocumentStorage = DocumentStorageStdout{}
