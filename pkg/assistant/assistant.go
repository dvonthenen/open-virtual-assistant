// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package assistant

import (
	"context"
	"os"

	klog "k8s.io/klog/v2"

	ainterfaces "github.com/dvonthenen/open-virtual-assistant/pkg/assistant/interfaces"
	sinterfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
	tinterfaces "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/interfaces"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
	config "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
	dgtranscriber "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/deepgram"
	gtranscriber "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/google"
)

func New(assistantImpl *ainterfaces.AssistantImpl, opts *AssistantOptions) (*Assistant, error) {
	ctx := context.Background()

	if opts == nil {
		opts = &AssistantOptions{}
	}

	// transcriber callback
	var callback tinterfaces.ResponseCallback
	callback = *assistantImpl

	// assistant
	assistant := &Assistant{
		speechOptions: &speech.SpeechOptions{
			VoiceType:    opts.VoiceType,
			LanguageCode: opts.LanguageCode,
		},
		transcriberOptions: &config.TranscribeOptions{
			InputChannels: opts.InputChannels,
			SamplingRate:  opts.SamplingRate,
			Callback:      &callback,
		},
		assistantImpl: assistantImpl,
	}

	// text-to-speech client
	speech, err := speech.New(ctx, assistant.speechOptions)
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		return nil, err
	}

	// which transcriber?
	var transcriberStr string
	if v := os.Getenv("ASSISTANT_TRANSCRIBER"); v != "" {
		klog.V(2).Infof("ASSISTANT_TRANSCRIBER found")
		transcriberStr = v
	}

	// get the transcriber
	var transcriber Transcriber

	switch transcriberStr {
	case ainterfaces.DEEPGRAM_TRANSCRIBER:
		transcribe, errTranscribe := dgtranscriber.New(ctx, assistant.transcriberOptions)
		if errTranscribe != nil {
			klog.V(1).Infof("dgtranscriber.New failed. Err: %v\n", errTranscribe)
			return nil, errTranscribe
		}
		transcriber = transcribe
	case ainterfaces.GOOGLE_TRANSCRIBER:
		fallthrough
	default:
		transcribe, errTranscribe := gtranscriber.New(ctx, assistant.transcriberOptions)
		if errTranscribe != nil {
			klog.V(1).Infof("gtranscriber.New failed. Err: %v\n", errTranscribe)
			return nil, errTranscribe
		}
		transcriber = transcribe
	}

	// housekeeping
	var playback sinterfaces.Speech
	playback = speech

	assistant.speech = speech
	assistant.transcriber = &transcriber
	(*assistantImpl).SetSpeech(&playback)
	assistant.assistantImpl = assistantImpl

	return assistant, nil
}

func (a *Assistant) Start() error {
	err := (*a.transcriber).Start()
	if err != nil {
		klog.V(1).Infof("transcriber.Start failed. Err: %v\n", err)
	}
	return err
}

func (a *Assistant) Speak(text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, a.speechOptions)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, text)
	klog.V(2).Infof("Response:\n%s\n\n", text)
	speechClient.Close()

	return nil
}

func (a *Assistant) Stop() error {
	err := (*a.transcriber).Stop()
	if err != nil {
		klog.V(1).Infof("transcriber.Stop failed. Err: %v\n", err)
	}
	return err
}
