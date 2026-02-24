# dr-docer


## Decision

### Config files

#### Markdown with markup

e.g.:
```
[[service:my-service]]
# This is my document
```

Pros:
 * Single document
 * Annotations are directly for the markdown itself

Cons:
 * Would probably need to be custom made
 * How would the markdown fit into templates? Rather than having markdown for specific sections (e.g. markdown for complete service loss), it's just "here's some markdown"
 * Would need custom interpretter (most likely?) - would mean difficult learning curve for users (obviously just me atm, but maybe open source)

#### YAML with markdown document

e.g.
```yaml
---
service: my-service
...
---

# This is my markdown content

here
```

Pros:
 * Single file
 * Markdown appears to not ender the first document
 * Rendering as YAML:
   * Second document is just a string and some variable corruption
   * Adding further YAML AFTER the markdown does not appear to work
      * Could in theory put YAML/YAML/MD - but how would we determine which entity the markdown relates to

#### YAML with markdown variables

e.g.:
```yaml
---
service: my-service
---
markdown: |
  mf
  asdad

---
```

pros:
 * Allows for multiple document and ordering of documents would work, so any markdown after entity - the markdown is for the entity.
 * We could include direct annotations for the markdown e.g. what is it

cons:
 * All markdown needs to be indented

#### Multiple files

e.g.:
```yaml
#my-service/metadata.yml
name: blah


#my-service/somemd.md
# My markdown
```

Pros:
 * No limitations

Cons:
 * Multiple files are harder to manage, if you imagine hundreds servers/services


## Ideas

 * Templates contain sections marked as:
   * MUST - throw error if they don't
   * SHOULD - shows a markdown error (red?) that this isn't provided
   * OPTIONAL - renders empty or with default

Template ideas:
```
# {{ .name }}

### Basic Info

| Ip Address | {{ .ip_address }}
| etc. | |

## Useful links

{{ .links }}

## Procedures

{{ template "procedures" . }}



{{ define "procedures" }}

### Setup

{{ template "procedures_setup" }}

{{ end }}
```



EntitySource ---> EntityFactory ---> EntityStore
     |               |                  |
      --(Creates)-->  --> Stores -------

Entity ---> Relationships
