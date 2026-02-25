package relationship

import (
	"fmt"
)

type RelationshipStore interface {
	GetEntityByName(name string) (*Entity, error)
	UpsertEntity(entity Entity) error
}

type RelationshipService struct {
	relationshipStore RelationshipStore
}

func NewRelationshipService(relationshipStore *RelationshipStore) (*RelationshipService, error) {
	if relationshipStore == nil {
		return nil, fmt.Errorf("NewRelationshipService: relationshipStore is nil")
	}
	return &RelationshipService{
		relationshipStore: *relationshipStore,
	}, nil
}

func (r *RelationshipService) getOrCreateEntity(name string) (*Entity, error) {
	entity, err := r.relationshipStore.GetEntityByName(name)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		entity = &Entity{
			Name:       name,
			Dependents: []Relationship{},
			DependsOn:  []Relationship{},
		}
	}
	return entity, nil
}

func (r *RelationshipService) AddEntityRelationship(name string, parentName string, relationshipType RelationshipType) error {
	entity, err := r.getOrCreateEntity(name)
	if err != nil {
		return err
	}

	parentEntity, err := r.getOrCreateEntity(parentName)
	if err != nil {
		return err
	}

	// Add depdency on entity
	entity.DependsOn = append(entity.Dependents, Relationship{
		Target: *parentEntity,
		Type:   relationshipType,
	})

	// Add relationship to parent
	parentEntity.Dependents = append(entity.Dependents, Relationship{
		Target: *entity,
		Type:   relationshipType,
	})

	// Store updated entities
	if err = r.relationshipStore.UpsertEntity(*entity); err != nil {
		return err
	}
	if err = r.relationshipStore.UpsertEntity(*parentEntity); err != nil {
		return err
	}
	return nil
}

func (r *RelationshipService) getEntityByName(name string) (*Entity, error) {
	entity, err := r.relationshipStore.GetEntityByName(name)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("Entity does not exist")
	}
	return entity, nil
}

// Get list of entities parent names
func (r *RelationshipService) GetEntityParents(name string, relationshipType RelationshipType) (*[]string, error) {
	entity, err := r.getEntityByName(name)
	if err != nil {
		return nil, err
	}
	var parents []string
	for _, relationship := range entity.DependsOn {
		if relationship.Type == relationshipType {
			parents = append(parents, relationship.Target.Name)
		}
	}
	return &parents, nil
}

// Get list of entities children names
func (r *RelationshipService) GetEntityChildren(name string, relationshipType RelationshipType) (*[]string, error) {
	entity, err := r.getEntityByName(name)
	if err != nil {
		return nil, err
	}
	var children []string
	for _, relationship := range entity.Dependents {
		if relationship.Type == relationshipType {
			children = append(children, relationship.Target.Name)
		}
	}
	return &children, nil
}
