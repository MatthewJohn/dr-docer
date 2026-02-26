# Dr-Docer

A framework for automatically generating Disaster Recovery (DR) documentation by discovering infrastructure entities from multiple sources, merging their attributes, and rendering Go templates to produce structured documentation.

> **Note:** This README documentation was written by AI. The code itself was not written by AI.

## Overview

Dr-Docer solves the problem of keeping DR documentation up-to-date by automatically discovering infrastructure entities (servers, services, etc.) from various sources, merging data from multiple providers with priority-based conflict resolution, and generating documentation from Go templates.

The framework is designed to be extended with custom entity sources and templates, allowing organizations to integrate with their existing infrastructure and documentation workflows.

## Architecture

```
EntitySource(s)     EntityFactory      EntityCollection
      |                   |                    |
      |--(provides)------>|----(stores)------->|
                           |
                           v
                    [Merged Entities]
                           |
                           v
                    DocumentGenerator
                           |
                    [Template Rendering]
                           |
                           v
                    [DocumentStorage] --> [DR Documentation]
```

### Core Concepts

**Entities** represent infrastructure components such as servers or services. Each entity has:
- A name (unique identifier)
- A type (e.g., server, service)
- A map of dynamic attributes with type-safe values
- Priority-based attribute merging from multiple sources

**Attributes** are typed, named values that can be attached to entities:
- Each attribute has a name, type, and optional default value
- Attribute instances have a priority for merge conflict resolution
- Type-safe value assignment with `SetValue()` ensures type correctness
- Common attributes (like IP address, URL) are shared across sources

**Entity Sources** are pluggable data providers that discover and provide entities. Sources can query external APIs, parse configuration files, read from filesystems, clone Git repositories, or connect to any other data source.

**The Entity Factory** orchestrates the discovery process:
1. Accepts registrations from multiple entity sources
2. Calls each source to populate the entity collection
3. Handles automatic merging of entities with the same name/type
4. Returns a unified collection ready for documentation generation

## Entity Types

### Common Entity Types

The framework provides common entity types in the `common_types` package:
- `EntityServer` - Represents physical or virtual servers
- `EntityService` - Represents applications or services

### Common Attributes

Common attributes are shared across entity sources:
- `AttributeIpAddress` - IP address of a server
- `AttributeUrl` - URL of a service

### Extending with Custom Attributes

To add custom attributes for your entity sources, define them in your package:

```go
var AttributeParentStorage attribute.Attribute = attribute.Attribute{
    Name:         "parent_storage",
    Type:         reflect.TypeOf(""),
    DefaultValue: "",
}
```

Then return them from your `GetAttributes()` method:

```go
func (m *MyDiscovery) GetAttributes() []attribute.Attribute {
    return []attribute.Attribute{
        commontypes.AttributeIpAddress,  // Common attribute
        AttributeParentStorage,           // Custom attribute
    }
}
```

### Extending with Custom Entity Types

To define custom entity types, simply define them as constants:

```go
var (
    EntityDatabase  metadata.EntityType = "database"
    EntityLoadBalancer metadata.EntityType = "load_balancer"
)
```

Then return them from your `GetEntityTypes()` method:

```go
func (m *MyDiscovery) GetEntityTypes() []metadata.EntityType {
    return []metadata.EntityType{
        EntityDatabase,
        EntityLoadBalancer,
    }
}
```

Attributes support:
- Type-safe values using `reflect.Type`
- `SetValue()` validates types before assignment
- Default values for when attributes are unset
- Priority-based merging (higher priority overrides lower)

## Extending the Framework

### Creating a Custom Entity Source

To add a new data source, implement the `EntitySource` interface:

```go
type EntitySource interface {
    GetEntities(collection *EntityCollection) error
    GetPriority() int
    GetEntityTypes() []metadata.EntityType
    GetAttributes() []attribute.Attribute
}
```

The `GetEntities` method receives the entity collection, allowing it to:
- Add new entities
- Merge additional attributes into existing entities

The `GetPriority` method determines the order sources are called (lower values run first).

The `GetEntityTypes` method returns the entity types that this source provides.

The `GetAttributes` method returns the attributes that this source uses.

### Example: Infrastructure-as-Code Discovery

The following example demonstrates discovering servers from infrastructure-as-code files in a Git repository:

```go
type IaCServerDiscovery struct {
    parser      *IaCParser
    gitService  *GitService
    repoUrl     string
}

func (d *IaCServerDiscovery) GetEntityTypes() []metadata.EntityType {
    return []metadata.EntityType{commontypes.EntityServer}
}

func (d *IaCServerDiscovery) GetAttributes() []attribute.Attribute {
    return []attribute.Attribute{
        commontypes.AttributeIpAddress,
    }
}

func (d *IaCServerDiscovery) GetEntities(collection *EntityCollection) error {
    // Clone the repository
    repo, err := d.gitService.CloneRepository(d.repoUrl)
    if err != nil {
        return err
    }

    // Parse infrastructure files
    resources, err := d.parser.Parse(repo)
    if err != nil {
        return err
    }

    // Extract server information from resources
    for _, resource := range resources {
        if resource.Type == "server" {
            entity, err := metadata.NewEntity(
                metadata.EntityName(resource.Name),
                commontypes.EntityServer,
                d.GetPriority(),
            )
            if err != nil {
                return err
            }
            entity.SetAttribute(&commontypes.AttributeIpAddress, resource.Attributes["ip_address"])
            collection.AddEntity(entity)
        }
    }

    return nil
}

func (d *IaCServerDiscovery) GetPriority() int {
    return 10
}
```

### Example: Terraform Module Discovery

For Terraform-based infrastructure, you can parse Terraform modules to extract entities:

```go
type TerraformDiscovery struct {
    parser      *terraform.TerraformParser
    gitService  *git.GitService
    repoUrl     string
}

func (d *TerraformDiscovery) GetEntities(collection *discovery.EntityCollection) error {
    // Clone git repo
    repo, err := d.gitService.CloneRepository(d.repoUrl)
    if err != nil {
        return err
    }

    // Parse Terraform files
    tfEntities, err := d.parser.ParseTerraform(repo, "")
    if err != nil {
        return err
    }

    // Extract entities from Terraform modules
    for _, tfModule := range tfEntities.Modules {
        name, ok := tfModule.Inputs["name"]
        if !ok || name.AsString() == "" {
            continue
        }

        entity, err := metadata.NewEntity(
            metadata.EntityName(name.AsString()),
            commontypes.EntityServer,
            d.GetPriority(),
        )
        if err != nil {
            return err
        }

        // Set attributes from Terraform inputs
        if ipAddress, ok := tfModule.Inputs["ip_address"]; ok {
            entity.SetAttribute(&commontypes.AttributeIpAddress, ipAddress.AsString())
        }

        collection.AddEntity(entity)
    }

    return nil
}
```

### Example: API-Based Discovery

```go
type APIDiscovery struct {
    client  *http.Client
    baseUrl string
}

func (d *APIDiscovery) GetEntityTypes() []metadata.EntityType {
    return []metadata.EntityType{commontypes.EntityServer}
}

func (d *APIDiscovery) GetAttributes() []attribute.Attribute {
    return []attribute.Attribute{commontypes.AttributeIpAddress}
}

func (d *APIDiscovery) GetEntities(collection *EntityCollection) error {
    // Query API for infrastructure data
    resp, err := d.client.Get(d.baseUrl + "/api/servers/")
    if err != nil {
        return err
    }

    // Parse response and add/update entities
    for _, item := range resp.Items {
        existing := collection.GetEntityByNameAndType(
            metadata.EntityName(item.Name),
            commontypes.EntityServer,
        )
        if existing != nil {
            // Merge additional attributes into existing entity
            existing.SetAttribute(&commontypes.AttributeIpAddress, item.IP)
        } else {
            // Add new entity
            entity, _ := metadata.NewEntity(
                metadata.EntityName(item.Name),
                commontypes.EntityServer,
                d.GetPriority(),
            )
            entity.SetAttribute(&commontypes.AttributeIpAddress, item.IP)
            collection.AddEntity(entity)
        }
    }

    return nil
}

func (d *APIDiscovery) GetPriority() int {
    return 5  // Higher priority than IaC source
}
```

### Example: Filesystem-Based Discovery

```go
type FilesystemDiscovery struct {
    baseDirectory string
    typeMappings  map[string]metadata.EntityType
}

func (d *FilesystemDiscovery) GetEntityTypes() []metadata.EntityType {
    return []metadata.EntityType{
        commontypes.EntityServer,
        commontypes.EntityService,
    }
}

func (d *FilesystemDiscovery) GetAttributes() []attribute.Attribute {
    return []attribute.Attribute{
        commontypes.AttributeIpAddress,
        commontypes.AttributeUrl,
    }
}

func (d *FilesystemDiscovery) GetEntities(collection *EntityCollection) error {
    // Walk directory structure
    filepath.Walk(d.baseDirectory, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }

        // Parse YAML files
        var entityConfig struct {
            Name      string
            Type      metadata.EntityType
            IpAddress string
            Url       string
        }
        file, _ := os.ReadFile(path)
        yaml.Unmarshal(file, &entityConfig)

        // Create and add entity
        entity, _ := metadata.NewEntity(
            metadata.EntityName(entityConfig.Name),
            entityConfig.Type,
            d.GetPriority(),
        )
        if entityConfig.IpAddress != "" {
            entity.SetAttribute(&commontypes.AttributeIpAddress, entityConfig.IpAddress)
        }
        if entityConfig.Url != "" {
            entity.SetAttribute(&commontypes.AttributeUrl, entityConfig.Url)
        }
        collection.AddEntity(entity)

        return nil
    })

    return nil
}
```

## Template System

The document generation system uses Go templates to generate documentation. Templates are markdown files with YAML front matter defining the entity type they apply to.

### Template Structure

```yaml
---
entity_type: service
---

# {{.Name}}

{{.Get "url"}}

## Overview

Service overview goes here.

## Restore Procedure

1. Step one
2. Step two
```

The YAML front matter must contain:
- `entity_type` - The entity type this template applies to

### Template Variables

Templates receive a `TemplateEntityShim` object with the following methods:

- `.Name` - The entity's name
- `.Get(attributeName string)` - Get an attribute value by name (returns empty string if not found)

### Type-Safe Attribute Values

Attributes are type-safe - the `SetValue()` method ensures values match the attribute's defined type:

```go
// This will succeed
entity.SetAttribute(&commontypes.AttributeIpAddress, "192.168.1.1")

// This will fail with an error (wrong type)
entity.SetAttribute(&commontypes.AttributeIpAddress, 12345)
```

### Document Storage Interface

Generated documents are written through the `DocumentStorage` interface:

```go
type DocumentStorage interface {
    StoreDocument(name metadata.EntityName, entityType metadata.EntityType, body []byte) error
}
```

The framework includes two built-in storage implementations in `pkg/infrastructure/document_storage`:

#### Stdout Storage

Writes documents to standard output for testing or debugging:

```go
import documentstorage "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/infrastructure/document_storage"

storage, _ := documentstorage.NewDocumentStorageStdout()
```

#### Filesystem Storage

Writes documents to a directory, organized by entity type:

```go
import documentstorage "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/infrastructure/document_storage"

storage, _ := documentstorage.NewDocumentStorageFile(&documentstorage.DocumentStorageFileConfig{
    OutputDirectory: "./output",
})
// Documents are written as: ./output/{entityType}/{entityName}.md
```

You can also implement custom storage backends to write documents to:
- Git repositories
- Wikis or Confluence
- S3 or other cloud storage

### Example: Generating Documents

```go
import (
    documentgenerator "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/document_generator"
    documentstorage "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/infrastructure/document_storage"
)

// Create document storage using built-in filesystem implementation
documentStorage, _ := documentstorage.NewDocumentStorageFile(&documentstorage.DocumentStorageFileConfig{
    OutputDirectory: "./output",
})

// Create document generator with template directory
docGen, _ := documentgenerator.NewDocumentGenerator(documentStorage, "./config/templates")

// Generate documents for all entities
for _, entity := range entities.GetEntities() {
    err := docGen.GenerateDocumentForEntity(entity)
    if err != nil {
        log.Printf("Error generating document for %s: %v", entity.GetName(), err)
    }
}
```

## Project Structure

```
dr-docer/
├── pkg/
│   ├── domains/
│   │   ├── discovery/         # Entity factory and source interfaces
│   │   ├── git/               # Git repository management
│   │   ├── metadata/          # Entity models and types
│   │   ├── attribute/         # Dynamic attribute system with type-safe SetValue
│   │   ├── common_types/      # Shared entity types and attributes
│   │   ├── document_generator/# Template rendering and document generation
│   │   └── terraform/         # Infrastructure-as-code parsing
│   └── infrastructure/
│       ├── gitlab/            # GitLab provider implementation
│       ├── discovery/         # Built-in filesystem discovery
│       └── document_storage/  # Built-in stdout and filesystem storage
├── templates/                 # Documentation templates
├── dr-docer-custom/           # Example custom implementation
│   ├── cmd/generator/         # Main application entry point
│   ├── config/
│   │   ├── main.yaml          # Example configuration
│   │   └── templates/         # Example templates
│   └── pkg/
│       ├── config/            # Configuration models
│       ├── discovery/         # Custom entity source implementations
│       └── infrastructure/    # Custom storage backends
```

The `dr-docer-custom` directory demonstrates how the framework can be extended with custom implementations while reusing the core discovery and entity management functionality.

## Usage Pattern

1. **Define Configuration**: Create a configuration file with your source credentials
2. **Implement Entity Sources**: Create sources for your data providers
3. **Register Sources**: Register each source with the entity factory
4. **Load Entities**: Execute discovery to populate and merge all infrastructure data
5. **Apply Templates**: Generate documentation using the discovered entities

### Main Application Pattern

```go
import (
    documentgenerator "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/document_generator"
    documentstorage "gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/infrastructure/document_storage"
)

func main() {
    // Load configuration
    config, _ := config.NewConfigFromFile("./config/main.yaml")

    // Create entity factory
    factory, _ := discovery.NewEntityFactory()

    // Initialize shared services
    gitService, _ := git.NewGitService()
    httpClient := http.Client{Timeout: time.Second * 10}

    // Register entity sources
    iacDiscovery := NewIaCDiscovery(config.InfrastructureRepo, gitService)
    factory.RegisterEntitySource(iacDiscovery)

    apiDiscovery := NewAPIDiscovery(config.InfraAPIUrl, &httpClient)
    factory.RegisterEntitySource(apiDiscovery)

    fsDiscovery := NewFilesystemDiscovery(config.DataDirectory)
    factory.RegisterEntitySource(fsDiscovery)

    // Discover and merge entities
    entities, _ := factory.LoadEntities()

    // Generate documentation using built-in filesystem storage
    documentStorage, _ := documentstorage.NewDocumentStorageFile(&documentstorage.DocumentStorageFileConfig{
        OutputDirectory: "./output",
    })
    docGen, _ := documentgenerator.NewDocumentGenerator(documentStorage, "./config/templates")

    for _, entity := range entities.GetEntities() {
        fmt.Printf("Entity: %s (%s)\n", entity.GetName(), entity.GetType())
        docGen.GenerateDocumentForEntity(entity)
    }
}
```

## Key Features

- **Pluggable Architecture**: Add new data sources by implementing a simple interface
- **Dynamic Attributes**: Type-safe attribute system with priority-based merging
- **Type-Safe Value Assignment**: `SetValue()` ensures values match attribute types
- **Entity Merging**: Combine data from multiple sources with automatic conflict resolution
- **Git Integration**: Built-in support for Git repositories with extensible provider pattern
- **Infrastructure-as-Code Parsing**: Extract infrastructure definitions from configuration files
- **Template-Based Generation**: Go template rendering with YAML front matter
- **Document Storage Interface**: Custom backends for writing generated documentation
- **Priority-Based Discovery**: Control the order and precedence of data sources

## Extensibility

The framework can be extended with:

- **New Entity Types**: Add custom entity types in the metadata package
- **New Data Sources**: Implement `EntitySource` for any API, file format, or service
- **New Git Providers**: Implement `GitRepoProvider` for GitHub, GitLab, Bitbucket, etc.
- **Custom Document Storage**: Implement `DocumentStorage` to write to any destination
- **New Templates**: Create Go templates with YAML front matter for any entity type
- **Custom Merge Logic**: Override entity merging behavior for specific use cases
