package main

// goenums: -json -sql -serde/value -statemachine
// SwapStatus
type swapStatus int

//
//func (s SwapStatus) IsCanceled() bool {
//	return s == SwapStatuses.Step1Canceled || s == SwapStatuses.Step1CanceledByBank
//}
//
//func (s SwapStatus) IsPending() bool {
//	return slices.Contains([]SwapStatus{
//		SwapStatuses.Step1Pending,
//		SwapStatuses.Step1BankMarkAmount,
//		SwapStatuses.Step1BankMarkCanceled,
//	}, s)
//}
//
//func (s SwapStatus) FromName2(name string) (SwapStatus, bool) {
//	if name == "wait_counterparty_sign" {
//		return SwapStatuses.WaitCounterpartySign, true
//	} else if name == "wait_sender_accept_cancel" {
//		return SwapStatuses.WaitSenderAcceptOrCancel, true
//	} else {
//		return new(SwapStatus).FromName(name)
//	}
//}

const (
	_invalidSwapStatus swapStatus = 0 // invalid

	//////////////////////////////////////// STEP1-START ///////////////////////////////////////////

	// pending
	// 刚创建swap请求 / 还没有选择counterparty， 或者目标是银行，等待银行的同意
	step1Pending swapStatus = 1000

	/////// 售卖给银行内部流程-START

	// canceled
	// 用户自己取消了
	// state: [final]
	step1Canceled swapStatus = 1030

	// step1_bank_mark_amount
	// 银行标记了金额，正在等待checker同意; 如果意见不一致，回到pending
	// state: -> waitSenderAcceptOrCancel, step1Pending
	step1BankMarkAmount swapStatus = 1020
	// step1_bank_mark_canceled
	// 银行标记为取消
	// state: -> step1CanceledByBank, step1Pending
	step1BankMarkCanceled swapStatus = 1021
	// canceled_by_bank
	// 被银行取消了
	// state: [final]
	step1CanceledByBank swapStatus = 1031

	/////// 售卖给银行内部流程-END

	//////////////////////////////////////// STEP1-END ///////////////////////////////////////////

	//////////////////////////////////////// STEP2-START ///////////////////////////////////////////

	// wait_counterparty_sign
	// 等待counterparty签名，等同于[waitSenderAcceptOrCancel]
	// state: -> waitSenderBroadcast
	waitCounterpartySign swapStatus = 2000

	/////// 售卖给银行内部流程-START

	// wait_sender_accept_cancel
	// 等待发送者接受或者取消; 接受后，系统会自动使用钱包进行签名。
	// state: -> waitSenderBroadcast
	waitSenderAcceptOrCancel swapStatus = 2000

	/////// 售卖给银行内部流程-END

	//////////////////////////////////////// STEP2-END ///////////////////////////////////////////

	// wait_sender_broadcast
	// 等待sender执行swap交易
	// state: -> waitTxConfirm
	waitSenderBroadcast swapStatus = 3000

	// wait_tx_confirm
	// 等待tx确认
	// state: -> success, expired, failed
	waitTxConfirm swapStatus = 3001

	// success
	// 成功
	// state: [final]
	success swapStatus = 4000

	// expired
	// 过期(超过了swap的deadline)
	// state: [final]
	expired swapStatus = 5000

	// failed
	// 失败(sender广播交易失败)
	// state: [final]
	failed swapStatus = 6000
)
