// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	assistant "github.com/dvonthenen/open-virtual-assistant/cmd/assistant/pkg/assistant"
	initlib "github.com/dvonthenen/open-virtual-assistant/pkg/init"
)

func main() {
	/*
		Init
	*/
	initlib.Init(initlib.AssistantInit{
		LogLevel: initlib.LogLevelStandard, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	/*
		Assistant Options
	*/
	oldDemo := false
	if v := os.Getenv("ASSISTANT_OLD_DEMO"); v != "" {
		oldDemo = strings.EqualFold(strings.ToLower(v), "true")
	}

	/*
		Assistant
	*/
	myAssistant, err := assistant.New(&assistant.AssistantOptions{
		OldDemo: oldDemo,
	})
	if err != nil {
		fmt.Printf("assistant.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nStarting the Open Virtual Assistant...\n\n")

	// blocking call
	err = myAssistant.Start()
	if err != nil {
		fmt.Printf("myAssistant.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// clean up
	myAssistant.Stop()
}
