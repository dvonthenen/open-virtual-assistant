// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	cfginterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/client/interfaces"

	assistantimpl "github.com/dvonthenen/open-virtual-assistant/cmd/assistant/pkg/assistant"
	assistant "github.com/dvonthenen/open-virtual-assistant/pkg/assistant"
	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
)

func main() {
	/*
		Init
	*/
	ctx := context.Background()

	assistant.Init(assistant.AssistantInit{
		LogLevel: assistant.LogLevelStandard,
	})

	/*
		Assistant
	*/
	callback := assistantimpl.NewInsightHandler(&speech.SpeechOpts{
		VoiceType: speech.SpeechVoiceFemale,
	})

	config := &cfginterfaces.StreamingConfig{
		InsightTypes: []string{"topic", "question", "action_item", "follow_up"},
		Config: cfginterfaces.Config{
			MeetingTitle:        "my-meeting",
			ConfidenceThreshold: 0.7,
			SpeechRecognition: cfginterfaces.SpeechRecognition{
				Encoding:        "LINEAR16",
				SampleRateHertz: 16000,
			},
		},
		Speaker: cfginterfaces.Speaker{
			Name:   "Jane Doe",
			UserID: "user@email.com",
		},
	}

	myAssistant, err := assistant.NewWithConfig(ctx, config, callback)
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

	// clean up
	myAssistant.Stop()
}
