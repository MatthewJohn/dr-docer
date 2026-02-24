package metadata

type EntityType string

const (
	EntityTypeSerer   EntityType = "server"
	EntityTypeService EntityType = "service"
)

type Entity interface {
	GetType() EntityType
	GetName() string
	MergeAttributes(new Entity)
}

type BaseEntity struct {
	Name string
	Type EntityType
}

func (b BaseEntity) GetName() string {
	return b.Name
}

func (b *BaseEntity) MergeAttributes(new Entity) {
	// No attributes to current merge
}

type EntityServer struct {
	BaseEntity
	IpAddress     string
	ParentStorage string
}

func (e *EntityServer) GetType() EntityType {
	return EntityTypeSerer
}

func (e *EntityServer) MergeAttributes(new Entity) {
	other, ok := new.(*EntityServer)
	if !ok {
		return
	}

	e.BaseEntity.MergeAttributes(other)

	if e.IpAddress == "" {
		e.IpAddress = other.IpAddress
	}
}

var _ Entity = &EntityServer{}

type EntityService struct {
	BaseEntity
	Url string
}

func (e *EntityService) GetType() EntityType {
	return EntityTypeService
}

func (e *EntityService) MergeAttributes(new Entity) {
	other, ok := new.(*EntityService)
	if !ok {
		return
	}

	e.BaseEntity.MergeAttributes(other)

	if e.Url == "" {
		e.Url = other.Url
	}
}

var _ Entity = &EntityService{}
