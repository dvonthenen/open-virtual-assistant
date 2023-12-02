// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	matchr "github.com/antzucaro/matchr"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
)

// handles question...
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

func HowAreYou() string {
	return "I am good. Thanks for asking."
}

func WhatTimeIsIt() string {
	dt := time.Now()
	return fmt.Sprintf("The time is currently %s pacific", dt.Format("3:04 PM"))
}

func WhatIsYourName() string {
	return fmt.Sprintf("My name is %s.", AssistantName)
}

func WhatIsYourQuest() string {
	return ResponseWhatIsYourQuest
}

func WhatIsUnladenSwallow() string {
	return ResponseUnladenSwallow
}

// My Assistant
func (a *MyAssistant) SetSpeech(s *interfaces.Speech) {
	a.speech = s
}

func (a *MyAssistant) Response(text string) error {
	text = strings.ToLower(text)
	klog.V(5).Infof("text: %s\n", text)

	/*
		Old Simple Demo
	*/
	for _, triggers := range QuestionResponses {
		for _, key := range triggers.keys {
			klog.V(5).Infof("Key: %s\n", key)

			if percent := matchr.JaroWinkler(key, text, false); percent > 0.95 {
				klog.V(5).Infof("\n\n--------------------------------\n")
				klog.V(3).Infof("MATCH (%f): %s = %s\n", percent, key, text)
				klog.V(5).Infof("\n--------------------------------\n\n")

				klog.V(2).Infof("Heard:\nMATCH (%f): %s = %s\n\n", percent, key, text)

				if a.speech == nil {
					klog.V(2).Infof("Unable to play reply audio: a.speech is nil\n")
					return ErrTextToSpeectInvalid
				}

				err := (*a.speech).Play(context.Background(), triggers.callback())
				if err != nil {
					klog.V(1).Infof("speech.Play failed. Err: %v\n", err)
				}

				// exit on first match!
				return nil
			}
		}
	}

	return nil
}
