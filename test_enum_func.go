package main

import (
	"fmt"
)

func testEnum() {
	// Test multiple names
	fmt.Println("=== Testing Multiple Names ===")
	fmt.Printf("Step1Initialized.Name(): %s\n", Tokenrequests.Step1Initialized.Name())
	fmt.Printf("Step1Initialized.NameWith(1): %s\n", Tokenrequests.Step1Initialized.NameWith(1))
	fmt.Printf("Step1Initialized.NameWith(2): %s\n", Tokenrequests.Step1Initialized.NameWith(2))
	fmt.Printf("Step1Initialized.NameWith(10): %s\n", Tokenrequests.Step1Initialized.NameWith(10))
	fmt.Printf("Step1Initialized.Names(): %v\n", Tokenrequests.Step1Initialized.Names())

	// Test FromName with different names
	fmt.Println("\n=== Testing FromName ===")
	if val, ok := Tokenrequests.FromName("started"); ok {
		fmt.Printf("FromName('started'): %s (%d)\n", val.Name(), val.Val())
	}
	if val, ok := Tokenrequests.FromName("init"); ok {
		fmt.Printf("FromName('init'): %s (%d)\n", val.Name(), val.Val())
	}
	if val, ok := Tokenrequests.FromName("begin"); ok {
		fmt.Printf("FromName('begin'): %s (%d)\n", val.Name(), val.Val())
	}

	// Test IsValid
	fmt.Println("\n=== Testing IsValid ===")
	fmt.Printf("Invalid.IsValid(): %v\n", Tokenrequests.Invalid.IsValid())
	fmt.Printf("Step1Initialized.IsValid(): %v\n", Tokenrequests.Step1Initialized.IsValid())

	// Test state machine
	fmt.Println("\n=== Testing State Machine ===")
	fmt.Printf("Step1Initialized.IsTerminalState(): %v\n", Tokenrequests.Step1Initialized.IsTerminalState())
	fmt.Printf("Step1Canceled.IsTerminalState(): %v\n", Tokenrequests.Step1Canceled.IsTerminalState())
	fmt.Printf("Step1Denied.IsTerminalState(): %v\n", Tokenrequests.Step1Denied.IsTerminalState())
	fmt.Printf("Step4Success.IsTerminalState(): %v\n", Tokenrequests.Step4Success.IsTerminalState())

	transitions := Tokenrequests.Step1Initialized.ValidTransitions()
	fmt.Printf("Step1Initialized.ValidTransitions(): ")
	for _, t := range transitions {
		fmt.Printf("%s ", t.Name())
	}
	fmt.Println()

	// Test comparable
	fmt.Println("\n=== Testing Comparable ===")
	fmt.Printf("Step1Initialized == Step1Initialized: %v\n", Tokenrequests.Step1Initialized == Tokenrequests.Step1Initialized)
	fmt.Printf("Step1Initialized == Step1Canceled: %v\n", Tokenrequests.Step1Initialized == Tokenrequests.Step1Canceled)
}
