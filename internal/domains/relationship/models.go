package relationship

type RelationshipType string

const (
	// Standard relationship based on dependencies
	RelationshipTypeNormal RelationshipType = "Normal"
	// This is a Host->Service relationship
	// This is either implied by a service within a server
	// or via hosting_platform attribute
	RelationshipTypeHost RelationshipType = "Host"
)

type Relationship struct {
	Type   RelationshipType
	Target Entity
}

type Entity struct {
	Name       string
	Dependents []Relationship
	DependsOn  []Relationship
}
