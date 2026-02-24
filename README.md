# Dr-Docer

A framework for automatically generating Disaster Recovery (DR) documentation by discovering infrastructure entities from multiple sources and applying templates to produce structured documentation.

> **Note:** This README documentation was written by AI. The code itself was not written by AI.

## Overview

Dr-Docer solves the problem of keeping DR documentation up-to-date by automatically discovering infrastructure entities (servers, services, etc.) from various sources, merging data from multiple providers, and generating documentation from templates.

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
                    [Templates] --------> [DR Documentation]
```

### Core Concepts

**Entities** represent infrastructure components such as servers or services. Each entity has:
- A name (unique identifier)
- A type (e.g., server, service)
- Attributes specific to the type (IP address, URL, storage, dependencies, etc.)
- The ability to merge attributes from other entities of the same name/type

**Entity Sources** are pluggable data providers that discover and provide entities. Sources can query external APIs, parse configuration files, read from filesystems, clone Git repositories, or connect to any other data source.

**The Entity Factory** orchestrates the discovery process:
1. Accepts registrations from multiple entity sources
2. Calls each source to populate the entity collection
3. Handles automatic merging of entities with the same name/type
4. Returns a unified collection ready for documentation generation

## Entity Types

### Server
- Represents physical or virtual servers
- Attributes: IP address, parent storage, host
- Example: A database server, load balancer, or application host

### Service
- Represents applications or services
- Attributes: URL, dependencies, criticality
- Example: A web service, API endpoint, or internal tool

## Extending the Framework

### Creating a Custom Entity Source

To add a new data source, implement the `EntitySource` interface:

```go
type EntitySource interface {
    GetEntities(existingCollection *EntityCollection) error
    GetPriority() int
}
```

The `GetEntities` method receives the existing entity collection, allowing it to:
- Add new entities
- Merge additional attributes into existing entities

The `GetPriority` method determines the order sources are called (lower values run first).

### Example: Infrastructure-as-Code Discovery

The following example demonstrates discovering servers from infrastructure-as-code files in a Git repository:

```go
type IaCServerDiscovery struct {
    parser      *IaCParser
    gitService  *GitService
    repoUrl     string
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
            collection.AddEntity(&metadata.EntityServer{
                BaseEntity: metadata.BaseEntity{
                    Name: resource.Name,
                    Type: metadata.EntityTypeServer,
                },
                IpAddress: resource.Attributes["ip_address"],
            })
        }
    }

    return nil
}

func (d *IaCServerDiscovery) GetPriority() int {
    return 10
}
```

### Example: API-Based Discovery

```go
type APIDiscovery struct {
    client  *http.Client
    baseUrl string
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
            item.Name,
            metadata.EntityTypeServer,
        )
        if existing != nil {
            // Merge additional attributes into existing entity
            server := existing.(*metadata.EntityServer)
            server.IpAddress = item.IP
        } else {
            // Add new entity
            collection.AddEntity(&metadata.EntityServer{
                BaseEntity: metadata.BaseEntity{
                    Name: item.Name,
                    Type: metadata.EntityTypeServer,
                },
                IpAddress: item.IP,
            })
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

func (d *FilesystemDiscovery) GetEntities(collection *EntityCollection) error {
    // Walk directory structure
    filepath.Walk(d.baseDirectory, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }

        // Parse YAML files
        var entityConfig struct {
            Name string
            Type string
            // ... other fields
        }
        file, _ := os.ReadFile(path)
        yaml.Unmarshal(file, &entityConfig)

        // Create and add entity
        entity := createEntityFromConfig(entityConfig)
        collection.AddEntity(entity)

        return nil
    })

    return nil
}
```

## Template System

Templates define the structure of generated documentation. The framework uses markdown templates with YAML front matter:

Example server template:
```yaml
---
name: postgres
host: docker-host-01
criticality: high
storage:
  volume: data
  participates_in_host_backup: true
dependencies:
  - dns
  - network
---

# Overview

Postgres database for internal services.

# Restore

1. Restore volume snapshot
2. Start container
3. Validate replication
```

Templates can define:
- **MUST** fields - Required, throw error if missing
- **SHOULD** fields - Recommended, show warning if missing
- **OPTIONAL** fields - Optional, render empty or default value

## Project Structure

```
dr-docer/
├── pkg/
│   ├── domains/
│   │   ├── discovery/       # Entity factory and source interfaces
│   │   ├── git/             # Git repository management
│   │   ├── metadata/        # Entity models and types
│   │   └── terraform/       # Infrastructure-as-code parsing
│   └── infrastructure/      # External service clients
├── templates/               # Documentation templates
├── dr-docer-custom/         # Example custom implementation
│   ├── cmd/generator/       # Main application entry point
│   └── pkg/
│       ├── config/          # Configuration models
│       └── discovery/       # Custom entity source implementations
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
func main() {
    // Load configuration
    config, _ := config.NewConfigFromFile("./config/main.yaml")

    // Create entity factory
    factory, _ := discovery.NewEntityFactory()

    // Initialize services
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

    // Generate documentation
    for _, entity := range entities.GetEntities() {
        fmt.Printf("Entity: %s (%s)\n", entity.GetName(), entity.GetType())
    }
}
```

## Key Features

- **Pluggable Architecture**: Add new data sources by implementing a simple interface
- **Entity Merging**: Combine data from multiple sources with automatic conflict resolution
- **Git Integration**: Built-in support for Git repositories with extensible provider pattern
- **Infrastructure-as-Code Parsing**: Extract infrastructure definitions from configuration files
- **Type-Safe Entities**: Strongly-typed entity models with compile-time safety
- **Template-Based Generation**: Structured documentation from configurable templates
- **Priority-Based Discovery**: Control the order and precedence of data sources

## Extensibility

The framework can be extended with:

- **New Entity Types**: Add custom entity types in the metadata package
- **New Data Sources**: Implement `EntitySource` for any API, file format, or service
- **New Git Providers**: Implement `GitRepoProvider` for GitHub, GitLab, Bitbucket, etc.
- **New Template Types**: Extend the template system for different output formats
- **Custom Merge Logic**: Override entity merging behavior for specific use cases
