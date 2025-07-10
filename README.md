# GoEnums: Enhanced Enum Generator for Go

GoEnums is a powerful Go tool that generates type-safe enums with a rich set of features, designed to eliminate boilerplate and enhance functionality. Simply define your enum type, add a `goenums:` directive, and let the generator do the rest.

## Features

- **Type-Safe Enums**: Generates structs and methods that make your enums robust and easy to use.
- **Serialization Support**: Out-of-the-box support for JSON, YAML, SQL, Text, and Binary formats.
- **State Machine Generation**: Define state transitions and terminal states directly in your enum comments.
- **Custom Naming & Grouping**: Assign multiple names to enum values and group them using tags.
- **Flexible Configuration**: Control serialization format (by name or value) and generated methods with simple flags.
- **Automatic Code Generation**: Integrates seamlessly with `go generate`.

## Quick Start

### 1. Define Your Enum
Create a standard Go enum definition. Add a `// goenums:` comment directive above the type definition to configure the generator.

```go
package mypackage

// goenums: -sql -json -serde/value -genName -statemachine
type orderStatus int

const (
    // invalid: this value is considered invalid
    unknown orderStatus = 0

    // state: -> processing, canceled
    // tag: pending
    pending orderStatus = 100

    // state: -> shipped, failed
    // tag: active
    processing orderStatus = 200

    // state: -> delivered
    // tag: active
    shipped orderStatus = 300

    // state: [final]
    // tag: terminal
    delivered orderStatus = 400

    // state: [final]
    // tag: terminal
    canceled orderStatus = 500

    // state: -> pending, canceled
    // tag: terminal
    failed orderStatus = 600
)
```

### 2. Run the Generator
Execute the generator to create the `_enums.go` file. It is recommended to add a `go:generate` directive to your source file for automation.

```sh
go run main.go your_file.go
```
Or, using `go:generate`:
```go
//go:generate go run github.com/your-repo/goenum main.go your_file.go
package mypackage
...
```

## Generator Directives

You control the generated code using flags in the `// goenums:` comment.

| Flag            | Description                                                                                             |
|-----------------|---------------------------------------------------------------------------------------------------------|
| `-sql`          | Generates `sql.Scanner` and `driver.Valuer` interfaces for database integration.                        |
| `-json`         | Generates `json.Marshaler` and `json.Unmarshaler` interfaces.                                           |
| `-yaml`         | Generates `yaml.Marshaler` and `yaml.Unmarshaler` interfaces.                                           |
| `-text`         | Generates `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces.                           |
| `-binary`       | Generates `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler` interfaces.                       |
| `-genName`      | Generates a `Name()` method that returns the string representation of the enum constant.                |
| `-serde/name`   | Sets the default serialization format to be the enum's name (string).                                   |
| `-serde/value`  | Sets the default serialization format to be the enum's underlying value (e.g., `int`).                  |
| `-statemachine` | Generates methods for state transitions (`CanTransitionTo`, `ValidTransitions`, `IsTerminalState`).     |


## Comment-Based Features

Enhance your enums with special comments inside the `const` block.

### State Machine
Define state transitions and final states for process validation.

- **Transition**: `// state: -> StateA, StateB` indicates that the current state can transition to `StateA` or `StateB`.
- **Final State**: `// state: [final]` marks the state as a terminal state with no further transitions.

This generates the following methods:
```go
func (o OrderStatus) CanTransitionTo(target OrderStatus) bool
func (o OrderStatus) ValidTransitions() []OrderStatus
func (o OrderStatus) IsTerminalState() bool
```

### Tagging
Group related enum values using tags.

- **Syntax**: `// tag: group1, group2`

This generates helper methods to check for group membership and retrieve all values within a group:
```go
func (o OrderStatus) IsPending() bool
func (o OrderStatus) IsActive() bool

func (c orderStatusesContainer) PendingSlice() []OrderStatus
func (c orderStatusesContainer) ActiveSlice() []OrderStatus
```

### Custom Naming
Provide alternative names for an enum value. The first name is the default for the `Name()` method.

- **Syntax**: `// name, altName, anotherName`

This generates a `NameWith(idx int)` method to access the alternative names.

### Invalid Value
Mark a specific value as invalid, which will be excluded from `All()` iterations and fail `IsValid()` checks.

- **Syntax**: `// invalid`

## Generated Code Example

The generator creates a new file (`<source>_enums.go`) containing the enum struct, a container for all values, and the methods you requested.

```go
// OrderStatus is a type that represents a single enum value.
type OrderStatus struct {
	orderStatus
}

// Verify that OrderStatus implements the Enum interface
var _ enums.Enum[int, OrderStatus] = OrderStatus{}

// OrderStatuses is a main entry point for accessing enum values.
var OrderStatuses = orderStatusesContainer{
	Pending:    OrderStatus{orderStatus: pending},
	Processing: OrderStatus{orderStatus: processing},
    // ... and so on
}

// Val returns the underlying enum value.
func (o OrderStatus) Val() int {
	return int(o.orderStatus)
}

// All returns an iterator over all valid enum values.
func (o OrderStatus) All() iter.Seq[OrderStatus] {
    // ...
}

// FromName finds an enum value by name.
func (o OrderStatus) FromName(name string) (OrderStatus, bool) {
    // ...
}

// Scan implements the database/sql.Scanner interface.
func (o *OrderStatus) Scan(value any) error {
    // ...
}

// Value implements the database/sql/driver.Valuer interface.
func (o OrderStatus) Value() (driver.Value, error) {
    // ...
}

// MarshalJSON implements the json.Marshaler interface.
func (o OrderStatus) MarshalJSON() ([]byte, error) {
    // ...
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License.