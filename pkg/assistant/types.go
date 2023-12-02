// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package assistant

import (
	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/assistant/interfaces"
	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
	config "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// Transcriber interface
type Transcriber interface {
	Start() error
	Stop() error
}

// assistant implementation
type AssistantOptions struct {
	InputChannels int
	SamplingRate  int

	VoiceType    texttospeechpb.SsmlVoiceGender
	LanguageCode string
}

type Assistant struct {
	transcriberOptions *config.TranscribeOptions
	speechOptions      *speech.SpeechOptions

	transcriber   *Transcriber
	speech        *speech.Client
	assistantImpl *interfaces.AssistantImpl
}
