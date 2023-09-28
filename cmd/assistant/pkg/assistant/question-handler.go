// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"context"
	"fmt"
	"time"

	matchr "github.com/antzucaro/matchr"
	micinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone/interfaces"
	klog "k8s.io/klog/v2"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
)

type InsightHandlerFunc func(options speech.SpeechOpts, text string) error

type QuestionResponse struct {
	keys     []string
	callback InsightHandlerFunc
}

var QuestionResponses = []QuestionResponse{
	{[]string{
		TriggerHowAreYouDoing1,
		TriggerHowAreYouDoing2,
		TriggerHowAreYouDoing3,
		TriggerHowAreYouDoing4,
		TriggerHowAreYouDoing5,
		TriggerHowAreYouDoing6,
		TriggerHowAreYouDoing7,
		TriggerHowAreYouDoing8,
		TriggerHowAreYouDoing9,
		TriggerHowAreYouDoing10,
		TriggerHowAreYouDoing11,
		TriggerHowAreYouDoing12,
	},
		HowAreYou,
	},
	{[]string{
		TriggerWhatTimeIsIt1,
		TriggerWhatTimeIsIt2,
	},
		WhatTimeIsIt,
	},
	{[]string{
		TriggerWhatIsYourName,
	},
		WhatIsYourName,
	},
	{[]string{
		TriggerWhatIsYourQuest,
	},
		WhatIsYourQuest,
	},
	{[]string{
		TriggerWhatIsYourQuest,
	},
		WhatIsYourQuest,
	},
	{[]string{
		TriggerUnladenSwallow,
	},
		WhatIsUnladenSwallow,
	},
}

func HowAreYou(options speech.SpeechOpts, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, "I am good. Thanks for asking.")
	fmt.Printf("Response:\nI am good. Thanks for asking.\n\n")
	speechClient.Close()

	return nil
}

func WhatTimeIsIt(options speech.SpeechOpts, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	dt := time.Now()
	speechClient.Play(ctx, fmt.Sprintf("The time is currently %s pacific", dt.Format("3:04 PM")))
	fmt.Printf("Response:\nThe time is currently %s pacific.\n\n", dt.Format("3:04 PM"))
	speechClient.Close()

	return nil
}

func WhatIsYourName(options speech.SpeechOpts, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, fmt.Sprintf("My name is %s.", AssistantName))
	fmt.Printf("Response:\nMy name is %s.\n\n", AssistantName)
	speechClient.Close()

	return nil
}

func WhatIsYourQuest(options speech.SpeechOpts, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, ResponseWhatIsYourQuest)
	fmt.Printf("Response:\n%s\n\n", ResponseWhatIsYourQuest)
	speechClient.Close()

	return nil
}

func WhatIsUnladenSwallow(options speech.SpeechOpts, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, ResponseUnladenSwallow)
	fmt.Printf("Response:\n%s\n\n", ResponseUnladenSwallow)
	speechClient.Close()

	return nil
}

func ExecuteActionOnMatch(mic micinterfaces.Microphone, options speech.SpeechOpts, text string) {
	klog.V(4).Infof("Content: %s\n", text)
	for _, triggers := range QuestionResponses {
		for _, key := range triggers.keys {
			klog.V(5).Infof("Key: %s\n", key)

			if percent := matchr.JaroWinkler(key, text, false); percent > 0.95 {
				// Mute the mic so we don't "listen" to ourself!
				mic.Mute()

				klog.V(5).Infof("\n\n--------------------------------\n")
				klog.V(3).Infof("MATCH (%f): %s = %s", percent, key, text)
				klog.V(5).Infof("\n--------------------------------\n\n")

				fmt.Printf("Heard:\nMATCH (%f): %s = %s\n\n", percent, key, text)

				err := triggers.callback(options, text)
				if err != nil {
					klog.V(1).Infof("actionFunc failed. Err: %v\n", err)
				}

				// Unmute!
				mic.Unmute()

				// exit on first match!
				return
			}
		}
	}
}
