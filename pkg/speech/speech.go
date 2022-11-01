// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

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
)

type SpeechOpts struct {
	VoiceType    texttospeechpb.SsmlVoiceGender
	LanguageCode string
}

type SpeechClient struct {
	config SpeechOpts

	speechClient                 *texttospeech.Client
	googleApplicationCredentials string
}

func New(ctx context.Context, config SpeechOpts) (*SpeechClient, error) {
	klog.V(6).Infof("speech.New ENTER\n")

	if config.LanguageCode == "" {
		config.LanguageCode = DefaultLanguageCode
	}
	if config.VoiceType == 0 {
		config.VoiceType = SpeechVoiceNeutral
	}

	var googleApplicationCredentials string
	if v := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); v != "" {
		klog.V(4).Info("GOOGLE_APPLICATION_CREDENTIALS found")
		googleApplicationCredentials = v
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

	speechClient := &SpeechClient{
		config:                       config,
		speechClient:                 googleClient,
		googleApplicationCredentials: googleApplicationCredentials,
	}

	klog.V(3).Infof("speech.New Succeeded\n")
	klog.V(6).Infof("speech.New LEAVE\n")

	return speechClient, nil
}

func (sc *SpeechClient) TextToSpeech(ctx context.Context, text string) ([]byte, error) {
	klog.V(6).Infof("SpeechClient.TextToSpeech ENTER\n")
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
			LanguageCode: sc.config.LanguageCode,
			SsmlGender:   sc.config.VoiceType,
		},
		// Select the type of audio file you want returned.
		// TODO: hardcoded since we only support MP3 currently
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := sc.speechClient.SynthesizeSpeech(ctx, &req)
	if err != nil {
		klog.V(1).Infof("speechClient.SynthesizeSpeech Failed. Err: %v\n", err)
		klog.V(6).Infof("SpeechClient.TextToSpeech LEAVE\n")
		return []byte{}, err
	}

	klog.V(3).Infof("SpeechClient.TextToSpeech Succeeded\n")
	klog.V(6).Infof("SpeechClient.TextToSpeech LEAVE\n")
	return resp.AudioContent, nil
}

func (sc *SpeechClient) PlayAudio(stream []byte) error {
	klog.V(6).Infof("SpeechClient.PlayAudio ENTER\n")

	stringReader := bytes.NewReader(stream)
	stringReadCloser := io.NopCloser(stringReader)

	streamer, format, err := mp3.Decode(stringReadCloser)
	if err != nil {
		klog.V(1).Infof("mp3.Decode Failed. Err: %v\n", err)
		klog.V(6).Infof("SpeechClient.PlayAudio LEAVE\n")
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
	klog.V(6).Infof("SpeechClient.PlayAudio LEAVE\n")

	return nil
}

func (sc *SpeechClient) Play(ctx context.Context, text string) error {
	klog.V(6).Infof("SpeechClient.Play ENTER\n")

	stream, err := sc.TextToSpeech(ctx, text)
	if err != nil {
		klog.V(1).Infof("TextToSpeech Failed. Err: %v\n", err)
		klog.V(6).Infof("SpeechClient.Play LEAVE\n")
		return err
	}

	err = sc.PlayAudio(stream)
	if err != nil {
		klog.V(1).Infof("PlayAudio Failed. Err: %v\n", err)
		klog.V(6).Infof("SpeechClient.Play LEAVE\n")
		return err
	}

	klog.V(3).Infof("Play Succeeded\n")
	klog.V(6).Infof("SpeechClient.Play LEAVE\n")
	return nil
}

func (sc *SpeechClient) Close() {
	sc.speechClient.Close()
}
