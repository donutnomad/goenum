

# 功能说明
这是一个golang枚举类型结构体生成器，将用户定义的枚举类型生成结构体，增强枚举功能

例如
```go
package test

// goenums: -sql -json -serde/value -genName
type tokenRequestStatus int

const (
	invalid          tokenRequestStatus = 0    
	Step1Initialized tokenRequestStatus = 1000 
	Step1Canceled    tokenRequestStatus = 9010 
	Step1MarkAllowed tokenRequestStatus = 1001 
	Step1MarkDenied  tokenRequestStatus = 1002 
	Step1Failed      tokenRequestStatus = 8010 
	Step1Denied      tokenRequestStatus = 7011 

	Step2WaitingPayment   tokenRequestStatus = 2000 
	Step2WaitingTxConfirm tokenRequestStatus = 2001 
	Step2Failed           tokenRequestStatus = 8020 

	Step3Initialized tokenRequestStatus = 3000 
	Step3MarkAllowed tokenRequestStatus = 3001 
	Step3Failed      tokenRequestStatus = 8030 

	Step4Success tokenRequestStatus = 4000 
)


```
在类型定义上，支持注释格式为: goenums: -sql -json -yaml -text -binary -serde/name -serde/value -genName -statemachine

** 需要导入本目录下的enums包，里面封装了序列号和反序列化的方法
** 需要继承接口(enums包下的Enum)
```go
package test

// Enum interface definition
type Enum[R comparable, Self comparable] interface {
	Val() R
	All() iter.Seq[Self]
	IsValid() bool
	FromName(name string) (Self, bool) // Return complete enum instance
	FromValue(value R) (Self, bool)    // Return complete enum instance
	SerdeFormat() Format
	Name() string // Enum name, required value
	String() string
}
```

# 功能解析:
-sql 将会生成与SQL有关的序列化和反序列化方法
```go
// Scan implements the database/sql.Scanner interface for TokenRequestStatus.
// It parses the database value and stores it in the enum.
// It returns an error if the value cannot be parsed.
func (t *TokenRequestStatus) Scan(value any) error {
	result, err := enums.SQLScan(*t, value)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// Value implements the database/sql/driver.Valuer interface for TokenRequestStatus.
// It returns the database representation of the enum value.
func (t TokenRequestStatus) Value() (driver.Value, error) {
	return enums.SQLValue(t)
}

```
-json 将会生成JSON相关的方法
```go
// MarshalJSON implements the json.Marshaler interface for TokenRequestStatus.
// It returns the JSON representation of the enum value as a byte slice.
func (t TokenRequestStatus) MarshalJSON() ([]byte, error) {
	return enums.MarshalJSON(t, t.tokenRequestStatus)
}

// UnmarshalJSON implements the json.Unmarshaler interface for TokenRequestStatus.
// It parses the JSON representation of the enum value from the byte slice.
// It returns an error if the input is not a valid JSON representation.
func (t *TokenRequestStatus) UnmarshalJSON(data []byte) error {
	result, err := enums.UnmarshalJSON(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
```
-yaml
```go
// MarshalYAML implements the yaml.Marshaler interface for TokenRequestStatus.
// It returns the YAML representation of the enum value.
func (t TokenRequestStatus) MarshalYAML() (any, error) {
	return enums.MarshalYAML(t, t.tokenRequestStatus)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for TokenRequestStatus.
// It parses the YAML representation of the enum value.
// It returns an error if the YAML does not contain a valid enum value.
func (t *TokenRequestStatus) UnmarshalYAML(node *yaml.Node) error {
	result, err := enums.UnmarshalYAML(*t, node)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
```
-text
```go
// MarshalText implements the encoding.TextMarshaler interface for TokenRequestStatus.
// It returns the text representation of the enum value as a byte slice.
func (t TokenRequestStatus) MarshalText() ([]byte, error) {
	return enums.MarshalText(t, t.tokenRequestStatus)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for TokenRequestStatus.
// It parses the text representation of the enum value from the byte slice.
// It returns an error if the byte slice does not contain a valid enum value.
func (t *TokenRequestStatus) UnmarshalText(data []byte) error {
	result, err := enums.UnmarshalText(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
```
-binary
```go

// MarshalBinary implements the encoding.BinaryMarshaler interface for TokenRequestStatus.
// It returns the binary representation of the enum value as a byte slice.
func (t TokenRequestStatus) MarshalBinary() ([]byte, error) {
	return enums.MarshalBinary(t, t.tokenRequestStatus)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for TokenRequestStatus.
// It parses the binary representation of the enum value from the byte slice.
// It returns an error if the byte slice does not contain a valid enum value.
func (t *TokenRequestStatus) UnmarshalBinary(data []byte) error {
	result, err := enums.UnmarshalBinary(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
```

-serde/name
控制接口返回值
```go
// SerdeFormat implements the Enum interface.
// It returns the format used for serialization.
func (t TokenRequestStatus) SerdeFormat() enums.Format {
	return enums.FormatName
}
```
-serde/value
```go
// SerdeFormat implements the Enum interface.
// It returns the format used for serialization.
func (t TokenRequestStatus) SerdeFormat() enums.Format {
	return enums.FormatValue
}
```
-statemachine
将会生成状态机相关方法
```

// CanTransitionTo checks if the current state can transition to the target state.
// Returns true if the transition is allowed, false otherwise.
func (o OrderStatus) CanTransitionTo(target OrderStatus) bool {
	transitions := o.ValidTransitions()
	for _, validTarget := range transitions {
		if validTarget == target {
			return true
		}
	}
	return false
}

// ValidTransitions returns all valid target states that this state can transition to.
// Returns an empty slice if this is a terminal state or has no defined transitions.
func (o OrderStatus) ValidTransitions() []OrderStatus {
	if o == OrderStatuses.OrderPending {
		return []OrderStatus{
			OrderStatuses.XProcessing,
			OrderStatuses.OrderCancelled,
		}
	}
	if o == OrderStatuses.XProcessing {
		return []OrderStatus{
			OrderStatuses.OrderShipped,
			OrderStatuses.OrderFailed,
		}
	}
	if o == OrderStatuses.OrderShipped {
		return []OrderStatus{
			OrderStatuses.OrderDelivered,
		}
	}
	if o == OrderStatuses.OrderFailed {
		return []OrderStatus{
			OrderStatuses.OrderPending,
			OrderStatuses.OrderCancelled,
		}
	}
	return []OrderStatus{}
}

// IsTerminalState returns true if this state is a terminal (final) state.
// Terminal states cannot transition to any other state.
func (o OrderStatus) IsTerminalState() bool {
	if o == OrderStatuses.OrderDelivered {
		return true
	}
	if o == OrderStatuses.OrderCancelled {
		return true
	}
	return false
}

func (o OrderStatus) TerminalStateSlice() []OrderStatus {
	return []OrderStatus {OrderStatuses.OrderDelivered, OrderStatuses.OrderCancelled}
}

```


## 枚举定义的注释功能
```go

// goenums: -sql -json -serde/value -genName
type tokenRequestStatus int

const (
	invalid          tokenRequestStatus = 0    
    // name,name2   [支持多个name，默认Name()string 方法返回第一个name，除此之外会增加一个NameWith(idx int)方法来获取其他的Name，该方法没有错误，超过索引返回最后一个name值]
    // invalid      [这是一个关键字，invalid不能作为name，当它独立一行出现时，表示这个值是Invalid的，会在继承的方法IsValid() bool体现]
    // 用户自己的注释
    // state: -> Step1Initialized, Step1Canceled [这是状态机支持，表示这个状态会转变为哪些状态)
	// state: [final] [这表示最终状态, TerminalState]
	// tag: init, band1 [这表示分组，tag会生成方法，InitSlice() [], Band1Slice() [], 每个元素会生成IsInit() bool, IsBand1()]
	Step1Initialized tokenRequestStatus = 1000 
	Step1Canceled    tokenRequestStatus = 9010 
	Step1MarkAllowed tokenRequestStatus = 1001 
	Step1MarkDenied  tokenRequestStatus = 1002 
	Step1Failed      tokenRequestStatus = 8010 
	Step1Denied      tokenRequestStatus = 7011 

	Step2WaitingPayment   tokenRequestStatus = 2000 
	Step2WaitingTxConfirm tokenRequestStatus = 2001 
	Step2Failed           tokenRequestStatus = 8020 

	Step3Initialized tokenRequestStatus = 3000 
	Step3MarkAllowed tokenRequestStatus = 3001 
	Step3Failed      tokenRequestStatus = 8030 

	Step4Success tokenRequestStatus = 4000 
)
```

## 生成的枚举参考
```

// TokenRequestStatus is a type that represents a single enum value.
// It combines the core information about the enum constant and it's defined fields.
type TokenRequestStatus struct {
	tokenRequestStatus
}

// Verify that TokenRequestStatus implements the Enum interface
var _ enums.Enum[int, TokenRequestStatus] = TokenRequestStatus{}

// tokenRequestStatusesContainer is the container for all enum values.
// It is private and should not be used directly use the public methods on the TokenRequestStatus type.
type tokenRequestStatusesContainer struct {
	Invalid               TokenRequestStatus
	Step1Initialized      TokenRequestStatus // 1000, Step1 process started (PENDING)
	Step1Canceled         TokenRequestStatus // 9010, User manually canceled, process ended (CANCELED)
	Step1MarkAllowed      TokenRequestStatus // 1001, Marked as approved (PENDING)
	Step1MarkDenied       TokenRequestStatus // 1002, Marked as denied (PENDING)
	Step1Failed           TokenRequestStatus // 8010, Mark failed (others disagreed), waiting for user manual handling, re-enter process 1000, (FAILED)
	Step1Denied           TokenRequestStatus // 7011, Already denied, process ended (REJECTED)
	Step2WaitingPayment   TokenRequestStatus // 2000, Step1 passed, Step2 process started, user needs to start transferring Token to specified account (PENDING)
	Step2WaitingTxConfirm TokenRequestStatus // 2001, User has transferred, waiting for confirmation (at this time get user's TxHash) (PENDING)
	Step2Failed           TokenRequestStatus // 8020, Transfer inconsistent or transaction failed on chain, need to automatically return to process 2000 (automatically handled by program)
	Step3Initialized      TokenRequestStatus // 3000, Step2 received, Step3 process started (fill in bank TransactionID) (PENDING)
	Step3MarkAllowed      TokenRequestStatus // 3001, Marked as approved (PENDING)
	Step3Failed           TokenRequestStatus // 8030, Mark failed (others disagreed), waiting for user manual handling, re-enter process 3000 (FAILED)
	Step4Success          TokenRequestStatus // 4000
}

// TokenRequestStatusRaw is a type alias for the underlying enum type tokenRequestStatus.
// It provides direct access to the raw enum values for cases where you need
// to work with the underlying type directly.
type TokenRequestStatusRaw = tokenRequestStatus

// TokenRequestStatuses is a main entry point using the TokenRequestStatus type.
// It it a container for all enum values and provides a convenient way to access all enum values and perform
// operations, with convenience methods for common use cases.
var TokenRequestStatuses = tokenRequestStatusesContainer{
	Invalid: TokenRequestStatus{
		tokenRequestStatus: invalid,
	},
	Step1Initialized: TokenRequestStatus{
		tokenRequestStatus: Step1Initialized,
	},
	Step1Canceled: TokenRequestStatus{
		tokenRequestStatus: Step1Canceled,
	},
	Step1MarkAllowed: TokenRequestStatus{
		tokenRequestStatus: Step1MarkAllowed,
	},
	Step1MarkDenied: TokenRequestStatus{
		tokenRequestStatus: Step1MarkDenied,
	},
	Step1Failed: TokenRequestStatus{
		tokenRequestStatus: Step1Failed,
	},
	Step1Denied: TokenRequestStatus{
		tokenRequestStatus: Step1Denied,
	},
	Step2WaitingPayment: TokenRequestStatus{
		tokenRequestStatus: Step2WaitingPayment,
	},
	Step2WaitingTxConfirm: TokenRequestStatus{
		tokenRequestStatus: Step2WaitingTxConfirm,
	},
	Step2Failed: TokenRequestStatus{
		tokenRequestStatus: Step2Failed,
	},
	Step3Initialized: TokenRequestStatus{
		tokenRequestStatus: Step3Initialized,
	},
	Step3MarkAllowed: TokenRequestStatus{
		tokenRequestStatus: Step3MarkAllowed,
	},
	Step3Failed: TokenRequestStatus{
		tokenRequestStatus: Step3Failed,
	},
	Step4Success: TokenRequestStatus{
		tokenRequestStatus: Step4Success,
	},
}

// allSlice returns a slice of all enum values.
// This method is useful for iterating over all enum values in a loop.
func (t tokenRequestStatusesContainer) allSlice() []TokenRequestStatus {
	return []TokenRequestStatus{
		TokenRequestStatuses.Invalid,
		TokenRequestStatuses.Step1Initialized,
		TokenRequestStatuses.Step1Canceled,
		TokenRequestStatuses.Step1MarkAllowed,
		TokenRequestStatuses.Step1MarkDenied,
		TokenRequestStatuses.Step1Failed,
		TokenRequestStatuses.Step1Denied,
		TokenRequestStatuses.Step2WaitingPayment,
		TokenRequestStatuses.Step2WaitingTxConfirm,
		TokenRequestStatuses.Step2Failed,
		TokenRequestStatuses.Step3Initialized,
		TokenRequestStatuses.Step3MarkAllowed,
		TokenRequestStatuses.Step3Failed,
		TokenRequestStatuses.Step4Success,
	}
}

// Val implements the Enum interface.
// It returns the underlying enum value.
func (t TokenRequestStatus) Val() int {
	return int(t.tokenRequestStatus)
}

// All implements the Enum interface.
// It returns an iterator over all enum values.
func (t TokenRequestStatus) All() iter.Seq[TokenRequestStatus] {
	return func(yield func(TokenRequestStatus) bool) {
		for _, v := range TokenRequestStatuses.allSlice() {
			if !v.IsValid() {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// FromName implements the Enum interface.
// It finds an enum value by name and returns the enum instance and a boolean indicating if found.
func (t TokenRequestStatus) FromName(name string) (TokenRequestStatus, bool) {
	for enum, enumName := range tokenrequeststatusNamesMap {
		if enumName == name {
			return enum, true
		}
	}
	var zero TokenRequestStatus
	return zero, false
}

// FromValue implements the Enum interface.
// It finds an enum instance by its underlying value and returns the enum instance and a boolean indicating if found.
func (t TokenRequestStatus) FromValue(value int) (TokenRequestStatus, bool) {
	for v := range t.All() {
		if v.Val() == value {
			return v, true
		}
	}
	var zero TokenRequestStatus
	return zero, false
}

// 我们为container也会生成一些Slice有关的方法

// All returns an iterator over all enum values.
// This is a convenience method that delegates to the zero value enum instance.
func (t tokenRequestStatusesContainer) All() iter.Seq[TokenRequestStatus] {
	return TokenRequestStatus{}.All()
}

// FromName finds an enum value by name and returns the enum instance and a boolean indicating if found.
// This is a convenience method that delegates to the zero value enum instance.
func (t tokenRequestStatusesContainer) FromName(name string) (TokenRequestStatus, bool) {
	return TokenRequestStatus{}.FromName(name)
}

// FromValue finds an enum instance by its underlying value and returns the enum instance and a boolean indicating if found.
// This is a convenience method that delegates to the zero value enum instance.
func (t tokenRequestStatusesContainer) FromValue(value int) (TokenRequestStatus, bool) {
	return TokenRequestStatus{}.FromValue(value)
}


```


## 其他功能
支持一个文件中定义多个枚举，生成后的结果也在一个文件中，生成的文件名为 原文件名_enums.go