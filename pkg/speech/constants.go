// Copyright 2022 contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package speech

import (
	"errors"

	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

const (
	SpeechVoiceNeutral = texttospeechpb.SsmlVoiceGender_NEUTRAL
	SpeechVoiceFemale  = texttospeechpb.SsmlVoiceGender_FEMALE
	SpeechVoiceMale    = texttospeechpb.SsmlVoiceGender_MALE
)

const (
	DefaultLanguageCode string = "en-US"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
