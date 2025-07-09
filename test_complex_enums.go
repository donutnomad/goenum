package main

import (
	"database/sql/driver"
	"github.com/donutnomad/goenum/enums"
	"iter"
)

// TokenRequestStatus is a type that represents a single enum value.
// It combines the core information about the enum constant and its defined fields.
type TokenRequestStatus struct {
	tokenRequestStatus
}

// Verify that TokenRequestStatus implements the Enum interface
var _ enums.Enum[int, TokenRequestStatus] = TokenRequestStatus{}

// tokenRequestStatusContainer is the container for all enum values.
// It is private and should not be used directly use the public methods on the TokenRequestStatus type.
type tokenRequestStatusContainer struct {
	Invalid               TokenRequestStatus
	Step1Initialized      TokenRequestStatus
	Step1MarkAllowed      TokenRequestStatus
	Step1MarkDenied       TokenRequestStatus
	Step1Denied           TokenRequestStatus
	Step1Failed           TokenRequestStatus
	Step1Canceled         TokenRequestStatus
	Step2WaitingPayment   TokenRequestStatus
	Step2WaitingTxConfirm TokenRequestStatus
	Step2Failed           TokenRequestStatus
	Step3Initialized      TokenRequestStatus
	Step3MarkAllowed      TokenRequestStatus
	Step3Failed           TokenRequestStatus
	Step4Success          TokenRequestStatus
}

// Tokenrequests is a main entry point using the TokenRequestStatus type.
// It is a container for all enum values and provides a convenient way to access all enum values and perform
// operations, with convenience methods for common use cases.
var Tokenrequests = tokenRequestStatusContainer{
	Invalid: TokenRequestStatus{
		tokenRequestStatus: invalid,
	},
	Step1Initialized: TokenRequestStatus{
		tokenRequestStatus: step1Initialized,
	},
	Step1MarkAllowed: TokenRequestStatus{
		tokenRequestStatus: step1MarkAllowed,
	},
	Step1MarkDenied: TokenRequestStatus{
		tokenRequestStatus: step1MarkDenied,
	},
	Step1Denied: TokenRequestStatus{
		tokenRequestStatus: step1Denied,
	},
	Step1Failed: TokenRequestStatus{
		tokenRequestStatus: step1Failed,
	},
	Step1Canceled: TokenRequestStatus{
		tokenRequestStatus: step1Canceled,
	},
	Step2WaitingPayment: TokenRequestStatus{
		tokenRequestStatus: step2WaitingPayment,
	},
	Step2WaitingTxConfirm: TokenRequestStatus{
		tokenRequestStatus: step2WaitingTxConfirm,
	},
	Step2Failed: TokenRequestStatus{
		tokenRequestStatus: step2Failed,
	},
	Step3Initialized: TokenRequestStatus{
		tokenRequestStatus: step3Initialized,
	},
	Step3MarkAllowed: TokenRequestStatus{
		tokenRequestStatus: step3MarkAllowed,
	},
	Step3Failed: TokenRequestStatus{
		tokenRequestStatus: step3Failed,
	},
	Step4Success: TokenRequestStatus{
		tokenRequestStatus: step4Success,
	},
}

// tokenrequeststatusNamesMap maps enum values to their names array
var tokenrequeststatusNamesMap = map[TokenRequestStatus][]string{
	Tokenrequests.Invalid:               {},
	Tokenrequests.Step1Initialized:      {},
	Tokenrequests.Step1MarkAllowed:      {},
	Tokenrequests.Step1MarkDenied:       {},
	Tokenrequests.Step1Denied:           {},
	Tokenrequests.Step1Failed:           {},
	Tokenrequests.Step1Canceled:         {},
	Tokenrequests.Step2WaitingPayment:   {},
	Tokenrequests.Step2WaitingTxConfirm: {},
	Tokenrequests.Step2Failed:           {},
	Tokenrequests.Step3Initialized:      {},
	Tokenrequests.Step3MarkAllowed:      {},
	Tokenrequests.Step3Failed:           {},
	Tokenrequests.Step4Success:          {},
}

// TokenrequestsRaw is a type alias for the underlying enum type tokenRequestStatus.
// It provides direct access to the raw enum values for cases where you need
// to work with the underlying type directly.
type TokenrequestsRaw = tokenRequestStatus

// allSlice returns a slice of all enum values.
func (t tokenRequestStatusContainer) allSlice() []TokenRequestStatus {
	return []TokenRequestStatus{
		Tokenrequests.Invalid,
		Tokenrequests.Step1Initialized,
		Tokenrequests.Step1MarkAllowed,
		Tokenrequests.Step1MarkDenied,
		Tokenrequests.Step1Denied,
		Tokenrequests.Step1Failed,
		Tokenrequests.Step1Canceled,
		Tokenrequests.Step2WaitingPayment,
		Tokenrequests.Step2WaitingTxConfirm,
		Tokenrequests.Step2Failed,
		Tokenrequests.Step3Initialized,
		Tokenrequests.Step3MarkAllowed,
		Tokenrequests.Step3Failed,
		Tokenrequests.Step4Success,
	}
}

// Val implements the Enum interface.
func (t TokenRequestStatus) Val() int {
	return int(t.tokenRequestStatus)
}

// All implements the Enum interface.
func (t TokenRequestStatus) All() iter.Seq[TokenRequestStatus] {
	return func(yield func(TokenRequestStatus) bool) {
		for _, v := range Tokenrequests.allSlice() {
			if !v.IsValid() {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// IsValid implements the Enum interface.
func (t TokenRequestStatus) IsValid() bool {
	return true
}

// Name implements the Enum interface.
// Returns the first name of the enum value.
func (t TokenRequestStatus) Name() string {
	if names, ok := tokenrequeststatusNamesMap[t]; ok && len(names) > 0 {
		return names[0]
	}
	return ""
}

// NameWith returns the name at the specified index.
// If the index is out of bounds, returns the last name.
func (t TokenRequestStatus) NameWith(idx int) string {
	names, ok := tokenrequeststatusNamesMap[t]
	if !ok || len(names) == 0 {
		return ""
	}
	if idx < 0 || idx >= len(names) {
		return names[len(names)-1]
	}
	return names[idx]
}

// Names returns all names of the enum value.
func (t TokenRequestStatus) Names() []string {
	if names, ok := tokenrequeststatusNamesMap[t]; ok {
		return names
	}
	return []string{}
}

// String implements the Stringer interface.
func (t TokenRequestStatus) String() string {
	return t.Name()
}

// SerdeFormat implements the Enum interface.
func (t TokenRequestStatus) SerdeFormat() enums.Format {
	return enums.FormatValue
}

// FromName implements the Enum interface.
func (t TokenRequestStatus) FromName(name string) (TokenRequestStatus, bool) {
	for enumValue, names := range tokenrequeststatusNamesMap {
		for _, n := range names {
			if n == name {
				return enumValue, true
			}
		}
	}
	var zero TokenRequestStatus
	return zero, false
}

// FromValue implements the Enum interface.
func (t TokenRequestStatus) FromValue(value int) (TokenRequestStatus, bool) {
	for _, v := range Tokenrequests.allSlice() {
		if v.Val() == value {
			return v, true
		}
	}
	var zero TokenRequestStatus
	return zero, false
}

// All container methods for convenience
func (t tokenRequestStatusContainer) All() iter.Seq[TokenRequestStatus] {
	return Tokenrequests.allSlice()[0].All()
}

func (t tokenRequestStatusContainer) FromName(name string) (TokenRequestStatus, bool) {
	return Tokenrequests.allSlice()[0].FromName(name)
}

func (t tokenRequestStatusContainer) FromValue(value int) (TokenRequestStatus, bool) {
	return Tokenrequests.allSlice()[0].FromValue(value)
}

// Scan implements the database/sql.Scanner interface for TokenRequestStatus.
func (t *TokenRequestStatus) Scan(value any) error {
	result, err := enums.SQLScan(*t, value)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// Value implements the database/sql/driver.Valuer interface for TokenRequestStatus.
func (t TokenRequestStatus) Value() (driver.Value, error) {
	return enums.SQLValue(t)
}

// MarshalJSON implements the json.Marshaler interface for TokenRequestStatus.
func (t TokenRequestStatus) MarshalJSON() ([]byte, error) {
	return enums.MarshalJSON(t, t.tokenRequestStatus)
}

// UnmarshalJSON implements the json.Unmarshaler interface for TokenRequestStatus.
func (t *TokenRequestStatus) UnmarshalJSON(data []byte) error {
	result, err := enums.UnmarshalJSON(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// CanTransitionTo checks if the current state can transition to the target state.
func (t TokenRequestStatus) CanTransitionTo(target TokenRequestStatus) bool {
	transitions := t.ValidTransitions()
	for _, validTarget := range transitions {
		if validTarget.tokenRequestStatus == target.tokenRequestStatus {
			return true
		}
	}
	return false
}

// ValidTransitions returns all vazlid target states that this state can transition to.
func (t TokenRequestStatus) ValidTransitions() []TokenRequestStatus {
	return []TokenRequestStatus{}
}

// IsTerminalState returns true if this state is a terminal (final) state.
func (t TokenRequestStatus) IsTerminalState() bool {
	return false
}

// TerminalStateSlice returns a slice of all terminal states.
func (t TokenRequestStatus) TerminalStateSlice() []TokenRequestStatus {
	return []TokenRequestStatus{}
}
