package documentgenerator

import "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"

type DocumentStorage interface {
	StoreDocument(name metadata.EntityName, entityType metadata.EntityType, body []byte) error
}
