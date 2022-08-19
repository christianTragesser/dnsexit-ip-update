package dnsexit

import "os"

func ExampleCLIWorkflow_invalid() {
	var testEvent Event

	os.Args[1] = "test"

	CLIWorkflow(testEvent)

	// Output: Invalid argument: 'test'
	//  'dnsexit -h' to list valid options
}

func ExampleCLIWorkflow_help() {
	var testEvent Event

	os.Args[1] = "-h"

	CLIWorkflow(testEvent)

	// Output: dnsexit options:
	//  dnsexit update -h
}
