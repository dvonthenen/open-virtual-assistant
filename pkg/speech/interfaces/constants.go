// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
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
