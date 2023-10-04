// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"context"
	"fmt"
	"time"

	klog "k8s.io/klog/v2"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
)

// handles question...
type InsightHandlerFunc func(options *speech.SpeechOptions, text string) error

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
		TriggerUnladenSwallow1,
		TriggerUnladenSwallow2,
		TriggerUnladenSwallow3,
		TriggerUnladenSwallow4,
	},
		WhatIsUnladenSwallow,
	},
}

func HowAreYou(options *speech.SpeechOptions, text string) error {
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

func WhatTimeIsIt(options *speech.SpeechOptions, text string) error {
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

func WhatIsYourName(options *speech.SpeechOptions, text string) error {
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

func WhatIsYourQuest(options *speech.SpeechOptions, text string) error {
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

func WhatIsUnladenSwallow(options *speech.SpeechOptions, text string) error {
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
