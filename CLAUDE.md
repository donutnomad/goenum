# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `goenum`, a Go enum type struct generator that transforms user-defined enum types into enhanced structs with rich functionality. The tool generates type-safe enums with serialization support, state machine capabilities, and tagging systems.

## Core Architecture

### Main Components

- **Generator (`main.go`)**: Main code generation engine that parses Go source files and generates enum structs
- **Pluralization (`plurals.go`)**: English pluralization logic for generating container names
- **Enums Package (`enums/`)**: Core enum interface definitions and serialization utilities
- **Template System**: Embedded Go template for generating enum code

### Key Data Structures

- `EnumInfo`: Represents a complete enum definition with metadata
- `EnumValue`: Individual enum constant with names, tags, and state transitions
- `EnumOptions`: Configuration flags for features like SQL, JSON, YAML serialization

## Common Development Commands

### Build and Run
```bash
make build          # Build the goenum binary (REQUIRED before first use)
make install        # Install to $GOPATH/bin
make generate       # Run generator on test files
./goenum <file.go>  # Generate enums for a specific file
```

**重要提醒**: 在首次使用或修改代码后，必须先运行 `make build` 来构建最新的goenum二进制文件，然后才能使用 `./goenum` 命令生成枚举代码。

### Testing and Quality
```bash
make test           # Run all tests
make fmt            # Format code
make lint           # Run golangci-lint
```

### Dependencies
```bash
make deps           # Download and tidy dependencies
```

## Code Generation Process

1. **Parse**: AST parsing of Go source files to find enum type definitions with `goenums:` comments
2. **Extract**: Parse enum options, constants, and comment annotations
3. **Generate**: Use Go templates to generate enhanced enum structs
4. **Output**: Write to `*_enums.go` files in same directory

## Enum Definition Syntax

Enum types are defined with special comments:
```go
// goenums: -sql -json -serde/value -genName -statemachine
type MyEnum int

const (
    // invalid              // Marks value as invalid
    // name1,name2          // Multiple names for enum value
    // tag: group1,group2   // Tags for grouping
    // state: -> Target1, Target2  // State transitions
    // state: [final]       // Terminal state
    MyValue MyEnum = 1
)
```

## Generated Code Features

### Core Interface
All generated enums implement `enums.Enum[R, Self]` interface providing:
- `Val()` - Access underlying value
- `All()` - Iterator over all valid values
- `IsValid()` - Validity checking
- `FromName()/FromValue()` - Construction from name/value
- `Name()/Names()` - Name access
- `String()` - String representation

### Serialization Support
- **SQL**: `database/sql` Scanner/Valuer interfaces
- **JSON**: `encoding/json` marshaling
- **YAML**: `gopkg.in/yaml.v3` marshaling
- **Text/Binary**: `encoding` interfaces

### State Machine
- `CanTransitionTo()` - Check valid transitions
- `ValidTransitions()` - Get all valid next states
- `IsTerminalState()` - Check if final state
- `TerminalStateSlice()` - Get all terminal states

### Tagging System
- `Is{Tag}()` - Check if enum has specific tag
- `{Tag}Slice()` - Get all enums with specific tag

## File Structure

- Root contains main generator and pluralization logic
- `enums/` package provides core interfaces and serialization utilities
- `test*.go` files contain example enum definitions
- Generated `*_enums.go` files contain the enhanced enum structs

## Testing Strategy

The codebase uses multiple test files (`test_*.go`) as examples and test cases. The `enums/` package has dedicated unit tests for serialization functionality.

## Dependencies

- Core Go standard library
- Optional: `gopkg.in/yaml.v3` for YAML support
- Build tools: `golangci-lint` for linting

## 注意事项
使用git mv取代mv
