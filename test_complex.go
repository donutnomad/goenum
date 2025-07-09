package main

// goenums: -sql -json -serde/value -genName -statemachine
type tokenRequestStatus int

const (
	// invalid
	invalid tokenRequestStatus = 0

	// started,init,begin
	// state: -> step1Canceled, step1MarkAllowed
	step1Initialized tokenRequestStatus = 1000

	// allowed,approved
	// tag: mark,step1
	step1MarkAllowed tokenRequestStatus = 1001

	// denied,rejected
	// tag: mark,step1
	step1MarkDenied tokenRequestStatus = 1002

	// denied,rejected
	// state: [final]
	step1Denied tokenRequestStatus = 7011

	// failed,error
	// tag: failed
	step1Failed tokenRequestStatus = 8010

	// canceled,cancelled
	// state: [final]
	step1Canceled tokenRequestStatus = 9010

	// waiting,payment
	// tag: step2
	step2WaitingPayment tokenRequestStatus = 2000

	// confirm,waiting
	// tag: step2
	step2WaitingTxConfirm tokenRequestStatus = 2001

	// failed,error
	// tag: failed
	step2Failed tokenRequestStatus = 8020

	// init,step3
	// tag: step3
	step3Initialized tokenRequestStatus = 3000

	// allowed,approved
	// tag: mark,step3
	step3MarkAllowed tokenRequestStatus = 3001

	// failed,error
	// tag: failed
	step3Failed tokenRequestStatus = 8030

	// success,completed,done
	// state: [final]
	step4Success tokenRequestStatus = 4000
)
