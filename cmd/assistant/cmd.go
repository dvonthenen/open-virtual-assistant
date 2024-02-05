// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"os"

	assistant "github.com/dvonthenen/open-virtual-assistant/pkg/assistant"
	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/assistant/interfaces"
	initlib "github.com/dvonthenen/open-virtual-assistant/pkg/init"

	assistantimpl "github.com/dvonthenen/open-virtual-assistant/cmd/assistant/impl"
)

func main() {
	/*
		Init
	*/
	initlib.Init(initlib.AssistantInit{
		LogLevel: initlib.LogLevelDefault, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	/*
		Assistant
	*/
	myAssistant := assistantimpl.New()

	var assistImpl interfaces.AssistantImpl
	assistImpl = myAssistant

	assist, err := assistant.New(&assistImpl, &assistant.AssistantOptions{})
	if err != nil {
		fmt.Printf("assistant.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nStarting the Open Virtual Assistant...\n\n")

	// blocking call
	err = assist.Start()
	if err != nil {
		fmt.Printf("myAssistant.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// clean up
	assist.Stop()
}
