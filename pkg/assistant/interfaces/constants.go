// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
)

// constants...
const (
	// transcriber options
	DEEPGRAM_TRANSCRIBER string = "deepgram"
	GOOGLE_TRANSCRIBER   string = "google"
)

const (
	SpeechVoiceNeutral = interfaces.SpeechVoiceNeutral
	SpeechVoiceFemale  = interfaces.SpeechVoiceFemale
	SpeechVoiceMale    = interfaces.SpeechVoiceMale
)
