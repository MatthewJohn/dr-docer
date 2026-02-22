package discovery

import (
	"fmt"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/internal/domains/metadata"
)

type EntitySource interface {
	GetEntities(existingCollection *EntityCollection) error
	// Determine the priority when the entity can be run
	GetPriority() int
}

type EntityStore interface {
	GetEntityByNameAndType(name string, entityType metadata.EntityType)
}

type EntityFactory struct {
	entitySources []EntitySource
}

func NewEntityFactory() (*EntityFactory, error) {
	return &EntityFactory{}, nil
}

func (m *EntityFactory) RegisterEntitySource(entitySource *EntitySource) error {
	if entitySource == nil {
		return fmt.Errorf("RegisterEntitySource: Cannot register nil entitySource")
	}
	m.entitySources = append(m.entitySources, *entitySource)
	return nil
}

func MergeEntities([]metadata.Entity, []metadata.Entity) {
	// Do some magic here
}

func (m *EntityFactory) LoadEntities() (*EntityCollection, error) {
	// Create empty collection
	entityCollection, err := NewEntityCollection()
	if err != nil {
		return nil, err
	}
	if entityCollection == nil {
		return nil, fmt.Errorf("LoadEntities: Unable to create entityCollection")
	}

	for _, entitySource := range m.entitySources {
		err := entitySource.GetEntities(entityCollection)
		if err != nil {
			// @TODO: Probably return error and skip the source
			return nil, err
		}
	}
	return entityCollection, nil
}

type EntityCollection struct {
	entities []metadata.Entity
}

func NewEntityCollection() (*EntityCollection, error) {
	return &EntityCollection{}, nil
}

func (e *EntityCollection) GetEntityByNameAndType(name string, entityType metadata.EntityType) *metadata.Entity {
	for _, entity := range e.entities {
		if entity.Name == name && entity.Type == entityType {
			return &entity
		}
	}
	return nil
}

func (e *EntityCollection) AddEntity(entity metadata.Entity) {
	// if
}
