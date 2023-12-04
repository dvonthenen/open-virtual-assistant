// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package speech

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
)

type SpeechOptions struct {
	VoiceType    texttospeechpb.SsmlVoiceGender
	LanguageCode string
}

type Client struct {
	options           *SpeechOptions
	client            *texttospeech.Client
	googleCredentials string
}

func New(ctx context.Context, opts *SpeechOptions) (*Client, error) {
	klog.V(6).Infof("speech.New ENTER\n")

	if opts.LanguageCode == "" {
		opts.LanguageCode = interfaces.DefaultLanguageCode
	}
	if opts.VoiceType == 0 {
		opts.VoiceType = interfaces.SpeechVoiceNeutral
	}

	var googleCredentials string
	if v := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); v != "" {
		klog.V(4).Info("GOOGLE_APPLICATION_CREDENTIALS found")
		googleCredentials = v
	} else {
		klog.Error("GOOGLE_APPLICATION_CREDENTIALS not found")
		klog.V(6).Infof("speech.New LEAVE\n")
		return nil, ErrInvalidInput
	}

	googleClient, err := texttospeech.NewClient(ctx)
	if err != nil {
		klog.V(1).Infof("texttospeech.NewClient failed. Err: %v\n", err)
		klog.V(6).Infof("speech.New LEAVE\n")
		return nil, err
	}

	client := &Client{
		options:           opts,
		client:            googleClient,
		googleCredentials: googleCredentials,
	}

	klog.V(3).Infof("speech.New Succeeded\n")
	klog.V(6).Infof("speech.New LEAVE\n")

	return client, nil
}

func (sc *Client) TextToSpeech(ctx context.Context, text string) ([]byte, error) {
	klog.V(6).Infof("Client.TextToSpeech ENTER\n")
	klog.V(4).Infof("text: %s\n", text)

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: sc.options.LanguageCode,
			SsmlGender:   sc.options.VoiceType,
		},
		// Select the type of audio file you want returned.
		// TODO: hardcoded since we only support MP3 currently
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := sc.client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		klog.V(1).Infof("client.SynthesizeSpeech Failed. Err: %v\n", err)
		klog.V(6).Infof("Client.TextToSpeech LEAVE\n")
		return []byte{}, err
	}

	klog.V(3).Infof("Client.TextToSpeech Succeeded\n")
	klog.V(6).Infof("Client.TextToSpeech LEAVE\n")
	return resp.AudioContent, nil
}

func (sc *Client) Write(stream []byte) (int, error) {
	size := len(stream)
	err := sc.PlayAudio(stream)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (sc *Client) PlayAudio(stream []byte) error {
	klog.V(6).Infof("Client.PlayAudio ENTER\n")

	stringReader := bytes.NewReader(stream)
	stringReadCloser := io.NopCloser(stringReader)

	streamer, format, err := mp3.Decode(stringReadCloser)
	if err != nil {
		klog.V(1).Infof("mp3.Decode Failed. Err: %v\n", err)
		klog.V(6).Infof("Client.PlayAudio LEAVE\n")
		return err
	}
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/60))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	speechData := buffer.Streamer(0, buffer.Len())

	done := make(chan bool)
	speaker.Play(beep.Seq(speechData, beep.Callback(func() {
		done <- true
	})))

	// wait until done... blocking!
	<-done

	klog.V(3).Infof("PlayAudio Succeeded\n")
	klog.V(6).Infof("Client.PlayAudio LEAVE\n")

	return nil
}

func (sc *Client) Play(ctx context.Context, text string) error {
	klog.V(6).Infof("Client.Play ENTER\n")

	stream, err := sc.TextToSpeech(ctx, text)
	if err != nil {
		klog.V(1).Infof("TextToSpeech Failed. Err: %v\n", err)
		klog.V(6).Infof("Client.Play LEAVE\n")
		return err
	}

	err = sc.PlayAudio(stream)
	if err != nil {
		klog.V(1).Infof("PlayAudio Failed. Err: %v\n", err)
		klog.V(6).Infof("Client.Play LEAVE\n")
		return err
	}

	klog.V(3).Infof("Play Succeeded\n")
	klog.V(6).Infof("Client.Play LEAVE\n")
	return nil
}

func (sc *Client) Close() {
	sc.client.Close()
}
