package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/donutnomad/goenum/enums"
	"iter"
)

// SwapStatus is a type that represents a single enum value.
// It combines the core information about the enum constant and its defined fields.
type SwapStatus struct {
	swapStatus
}

// Verify that SwapStatus implements the Enum interface
var _ enums.Enum[int, SwapStatus] = SwapStatus{}

// swapStatusContainer is the container for all enum values.
// It is private and should not be used directly use the public methods on the SwapStatus type.
type swapStatusContainer struct {
	// invalid (0)
	_invalidSwapStatus SwapStatus
	// pending (1000)
	// 刚创建swap请求 / 还没有选择counterparty， 或者目标是银行，等待银行的同意
	Step1Pending SwapStatus
	// canceled (1030)
	// 用户自己取消了
	// state: [final]
	Step1Canceled SwapStatus
	// step1_bank_mark_amount (1020)
	// 银行标记了金额，正在等待checker同意; 如果意见不一致，回到pending
	// state: -> waitSenderAcceptOrCancel, step1Pending
	Step1BankMarkAmount SwapStatus
	// step1_bank_mark_canceled (1021)
	// 银行标记为取消
	// state: -> step1CanceledByBank, step1Pending
	Step1BankMarkCanceled SwapStatus
	// canceled_by_bank (1031)
	// 被银行取消了
	// state: [final]
	Step1CanceledByBank SwapStatus
	// wait_counterparty_sign (2000)
	// 等待counterparty签名，等同于[waitSenderAcceptOrCancel]
	// state: -> waitSenderBroadcast
	WaitCounterpartySign SwapStatus
	// wait_sender_accept_cancel (2000)
	// 等待发送者接受或者取消; 接受后，系统会自动使用钱包进行签名。
	// state: -> waitSenderBroadcast
	WaitSenderAcceptOrCancel SwapStatus
	// wait_sender_broadcast (3000)
	// 等待sender执行swap交易
	// state: -> waitTxConfirm
	WaitSenderBroadcast SwapStatus
	// wait_tx_confirm (3001)
	// 等待tx确认
	// state: -> success, expired, failed
	WaitTxConfirm SwapStatus
	// success (4000)
	// 成功
	// state: [final]
	Success SwapStatus
	// expired (5000)
	// 过期(超过了swap的deadline)
	// state: [final]
	Expired SwapStatus
	// failed (6000)
	// 失败(sender广播交易失败)
	// state: [final]
	Failed SwapStatus
}

// Swaps is a main entry point using the SwapStatus type.
// It is a container for all enum values and provides a convenient way to access all enum values and perform
// operations, with convenience methods for common use cases.
var Swaps = swapStatusContainer{
	_invalidSwapStatus: SwapStatus{
		swapStatus: _invalidSwapStatus,
	},
	Step1Pending: SwapStatus{
		swapStatus: step1Pending,
	},
	Step1Canceled: SwapStatus{
		swapStatus: step1Canceled,
	},
	Step1BankMarkAmount: SwapStatus{
		swapStatus: step1BankMarkAmount,
	},
	Step1BankMarkCanceled: SwapStatus{
		swapStatus: step1BankMarkCanceled,
	},
	Step1CanceledByBank: SwapStatus{
		swapStatus: step1CanceledByBank,
	},
	WaitCounterpartySign: SwapStatus{
		swapStatus: waitCounterpartySign,
	},
	WaitSenderAcceptOrCancel: SwapStatus{
		swapStatus: waitSenderAcceptOrCancel,
	},
	WaitSenderBroadcast: SwapStatus{
		swapStatus: waitSenderBroadcast,
	},
	WaitTxConfirm: SwapStatus{
		swapStatus: waitTxConfirm,
	},
	Success: SwapStatus{
		swapStatus: success,
	},
	Expired: SwapStatus{
		swapStatus: expired,
	},
	Failed: SwapStatus{
		swapStatus: failed,
	},
}

// swapstatusNamesMap maps enum values to their names array
var swapstatusNamesMap = map[SwapStatus][]string{
	Swaps._invalidSwapStatus: {
		"_invalidSwapStatus",
	},
	Swaps.Step1Pending: {
		"pending",
	},
	Swaps.Step1Canceled: {
		"canceled",
	},
	Swaps.Step1BankMarkAmount: {
		"step1_bank_mark_amount",
	},
	Swaps.Step1BankMarkCanceled: {
		"step1_bank_mark_canceled",
	},
	Swaps.Step1CanceledByBank: {
		"canceled_by_bank",
	},
	Swaps.WaitCounterpartySign: {
		"wait_counterparty_sign",
	},
	Swaps.WaitSenderAcceptOrCancel: {
		"wait_sender_accept_cancel",
	},
	Swaps.WaitSenderBroadcast: {
		"wait_sender_broadcast",
	},
	Swaps.WaitTxConfirm: {
		"wait_tx_confirm",
	},
	Swaps.Success: {
		"success",
	},
	Swaps.Expired: {
		"expired",
	},
	Swaps.Failed: {
		"failed",
	},
}

// SwapsRaw is a type alias for the underlying enum type swapStatus.
// It provides direct access to the raw enum values for cases where you need
// to work with the underlying type directly.
type SwapsRaw = swapStatus

// allSlice returns a slice of all enum values.
func (t swapStatusContainer) allSlice() []SwapStatus {
	return []SwapStatus{
		Swaps._invalidSwapStatus,
		Swaps.Step1Pending,
		Swaps.Step1Canceled,
		Swaps.Step1BankMarkAmount,
		Swaps.Step1BankMarkCanceled,
		Swaps.Step1CanceledByBank,
		Swaps.WaitCounterpartySign,
		Swaps.WaitSenderAcceptOrCancel,
		Swaps.WaitSenderBroadcast,
		Swaps.WaitTxConfirm,
		Swaps.Success,
		Swaps.Expired,
		Swaps.Failed,
	}
}

// Val implements the Enum interface.
func (t SwapStatus) Val() int {
	return int(t.swapStatus)
}

// All implements the Enum interface.
func (t SwapStatus) All() iter.Seq[SwapStatus] {
	return func(yield func(SwapStatus) bool) {
		for _, v := range Swaps.allSlice() {
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
func (t SwapStatus) IsValid() bool {
	if t == Swaps._invalidSwapStatus {
		return false
	}
	return true
}

// Name implements the Enum interface.
// Returns the first name of the enum value.
func (t SwapStatus) Name() string {
	if names, ok := swapstatusNamesMap[t]; ok && len(names) > 0 {
		return names[0]
	}
	return ""
}

// NameWith returns the name at the specified index.
// If the index is out of bounds, returns the last name.
func (t SwapStatus) NameWith(idx int) string {
	names, ok := swapstatusNamesMap[t]
	if !ok || len(names) == 0 {
		return ""
	}
	if idx < 0 || idx >= len(names) {
		return names[len(names)-1]
	}
	return names[idx]
}

// Names returns all names of the enum value.
func (t SwapStatus) Names() []string {
	if names, ok := swapstatusNamesMap[t]; ok {
		return names
	}
	return []string{}
}

// String implements the Stringer interface.
func (t SwapStatus) String() string {
	if names, ok := swapstatusNamesMap[t]; ok && len(names) > 0 {
		return names[0]
	}
	return fmt.Sprintf("swapStatus(%v)", t.swapStatus)
}

// SerdeFormat implements the Enum interface.
func (t SwapStatus) SerdeFormat() enums.Format {
	return enums.FormatValue
}

// FromName implements the Enum interface.
func (t SwapStatus) FromName(name string) (SwapStatus, bool) {
	for enumValue, names := range swapstatusNamesMap {
		for _, n := range names {
			if n == name {
				return enumValue, enumValue.IsValid()
			}
		}
	}
	var zero SwapStatus
	return zero, false
}

// FromValue implements the Enum interface.
func (t SwapStatus) FromValue(value int) (SwapStatus, bool) {
	for v := range Swaps.All() {
		if v.Val() == value {
			return v, true
		}
	}
	var zero SwapStatus
	return zero, false
}

// All container methods for convenience
func (t swapStatusContainer) All() iter.Seq[SwapStatus] {
	return SwapStatus{}.All()
}

func (t swapStatusContainer) FromName(name string) (SwapStatus, bool) {
	return SwapStatus{}.FromName(name)
}

func (t swapStatusContainer) FromValue(value int) (SwapStatus, bool) {
	return SwapStatus{}.FromValue(value)
}

// Scan implements the database/sql.Scanner interface for SwapStatus.
func (t *SwapStatus) Scan(value any) error {
	result, err := enums.SQLScan(*t, value)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// Value implements the database/sql/driver.Valuer interface for SwapStatus.
func (t SwapStatus) Value() (driver.Value, error) {
	return enums.SQLValue(t)
}

// MarshalJSON implements the json.Marshaler interface for SwapStatus.
func (t SwapStatus) MarshalJSON() ([]byte, error) {
	return enums.MarshalJSON(t, t.swapStatus)
}

// UnmarshalJSON implements the json.Unmarshaler interface for SwapStatus.
func (t *SwapStatus) UnmarshalJSON(data []byte) error {
	result, err := enums.UnmarshalJSON(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// CanTransitionTo checks if the current state can transition to the target state.
func (t SwapStatus) CanTransitionTo(target SwapStatus) bool {
	transitions := t.ValidTransitions()
	for _, validTarget := range transitions {
		if validTarget == target {
			return true
		}
	}
	return false
}

// ValidTransitions returns all valid target states that this state can transition to.
func (t SwapStatus) ValidTransitions() []SwapStatus {
	if t == Swaps.Step1BankMarkAmount {
		return []SwapStatus{
			Swaps.WaitSenderAcceptOrCancel,
			Swaps.Step1Pending,
		}
	}
	if t == Swaps.Step1BankMarkCanceled {
		return []SwapStatus{
			Swaps.Step1CanceledByBank,
			Swaps.Step1Pending,
		}
	}
	if t == Swaps.WaitCounterpartySign {
		return []SwapStatus{
			Swaps.WaitSenderBroadcast,
		}
	}
	if t == Swaps.WaitSenderAcceptOrCancel {
		return []SwapStatus{
			Swaps.WaitSenderBroadcast,
		}
	}
	if t == Swaps.WaitSenderBroadcast {
		return []SwapStatus{
			Swaps.WaitTxConfirm,
		}
	}
	if t == Swaps.WaitTxConfirm {
		return []SwapStatus{
			Swaps.Success,
			Swaps.Expired,
			Swaps.Failed,
		}
	}
	return []SwapStatus{}
}

// IsTerminalState returns true if this state is a terminal (final) state.
func (t SwapStatus) IsTerminalState() bool {
	if t == Swaps.Step1Canceled {
		return true
	}
	if t == Swaps.Step1CanceledByBank {
		return true
	}
	if t == Swaps.Success {
		return true
	}
	if t == Swaps.Expired {
		return true
	}
	if t == Swaps.Failed {
		return true
	}
	return false
}

// TerminalStateSlice returns a slice of all terminal states.
func (t SwapStatus) TerminalStateSlice() []SwapStatus {
	return []SwapStatus{
		Swaps.Step1Canceled,
		Swaps.Step1CanceledByBank,
		Swaps.Success,
		Swaps.Expired,
		Swaps.Failed,
	}
}
