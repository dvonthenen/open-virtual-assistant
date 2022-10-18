// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
	klog "k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	flag.Set("v", "6")
	flag.Parse()

	// Instantiates a client.
	ctx := context.Background()

	client, err := speech.New(ctx, speech.SpeechInit{
		VoiceType:    speech.SpeechVoiceFemale,
		LanguageCode: speech.DefaultLanguageCode,
	})
	if err != nil {
		fmt.Errorf("New failed. Err: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	text := "How much wood could a woodchuck chuck? If a woodchuck could chuck wood? As much wood as a woodchuck could chuck, If a woodchuck could chuck wood."

	err = client.Play(ctx, text)
	if err != nil {
		fmt.Errorf("Play failed. Err: %v\n", err)
		os.Exit(1)
	}
}
