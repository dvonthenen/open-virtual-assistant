// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import "errors"

const (
	QuestionHow   string = "How"
	QuestionWho   string = "Who"
	QuestionWhat  string = "What"
	QuestionWhen  string = "When"
	QuestionWhere string = "Where"
	QuestionWhy   string = "Why"

	AssistantName string = "Kitt"
)

// triggers
const (
	TriggerHowAreYouDoing1  string = "How are you doing today?"
	TriggerHowAreYouDoing2  string = "Hello, How are you doing today?"
	TriggerHowAreYouDoing3  string = "Hi, How are you doing today?"
	TriggerHowAreYouDoing4  string = "How are you doing?"
	TriggerHowAreYouDoing5  string = "Hello, How are you doing?"
	TriggerHowAreYouDoing6  string = "Hi, How are you doing?"
	TriggerHowAreYouDoing7  string = "How are you today?"
	TriggerHowAreYouDoing8  string = "Hello, How are you today?"
	TriggerHowAreYouDoing9  string = "Hi, How are you today?"
	TriggerHowAreYouDoing10 string = "How are you?"
	TriggerHowAreYouDoing11 string = "Hi, How are you?"
	TriggerHowAreYouDoing12 string = "Hello, How are you?"

	TriggerWhatTimeIsIt1 string = "What time is it?"
	TriggerWhatTimeIsIt2 string = "Do you have the time?"

	TriggerWhatIsYourName  string = "What is your name?"
	TriggerWhatIsYourQuest string = "What is your quest?"
	TriggerUnladenSwallow1 string = "What is the air-speed velocity of an unladen swallow?"
	TriggerUnladenSwallow2 string = "What is the air speed velocity of an un latent swallow?"
	TriggerUnladenSwallow3 string = "What is the airs speed velocity of an un laden swallow?"
	TriggerUnladenSwallow4 string = "What is the airspeed velocity of an unladen swallow?"
)

// response
const (
	ResponseWhatIsYourQuest string = "To seek the Holy Grail."
	ResponseUnladenSwallow  string = "What do you mean? An African or a European swallow?"
)

var (
	// ErrTextToSpeectInvalid text-to-speech interface is nil
	ErrTextToSpeectInvalid = errors.New("text-to-speech interface is nil")
)
