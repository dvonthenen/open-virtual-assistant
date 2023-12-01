// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package deepgram

import (
	"context"
	"errors"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	live "github.com/deepgram/deepgram-go-sdk/pkg/client/live"

	microphone "github.com/dvonthenen/open-virtual-assistant/pkg/microphone"
	config "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
)

type Transcribe struct {
	options *config.TranscribeOptions

	client *live.Client
	mic    *microphone.Microphone
}

var micInitAlready = false

func New(ctx context.Context, opts *config.TranscribeOptions) (*Transcribe, error) {
	klog.V(6).Infof("transcribe.New ENTER\n")

	if opts.InputChannels == 0 {
		opts.InputChannels = 1
	}
	if opts.SamplingRate == 0 {
		opts.SamplingRate = 16000
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// mic stuf
	if !micInitAlready {
		klog.V(4).Infof("Calling microphone.Initialize...")
		microphone.Initialize()
		micInitAlready = true
	}

	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: opts.InputChannels,
		SamplingRate:  float32(opts.SamplingRate),
	})
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		return nil, err
	}

	// Deepgram init
	options := interfaces.LiveTranscriptionOptions{
		Language:   "en-US",
		Encoding:   "linear16",
		Channels:   opts.InputChannels,
		SampleRate: opts.SamplingRate,
		Punctuate:  true,
		// Keywords:   []string{"Hey Kitt:32", "Hey Kit:16", "Hey:16", "Hello:16", "Hey:16", "Kitt:16", "Kit:16"},
		// Endpointing: "500",
	}
	// klog.V(2).Infof("options: %v\n", options)

	handler := NewInsightHandler(&InsightOptions{
		TranscribeOptions: opts,
		Microphone:        mic,
	})

	// create a new client
	client, err := live.NewWithDefaults(ctx, options, handler)
	if err != nil {
		klog.V(1).Infof("NewDeepGramWSClientDefault failed", err)
		return nil, err
	}

	// crete client
	transcribe := &Transcribe{
		options: opts,
		client:  client,
		mic:     mic,
	}

	klog.V(3).Infof("transcribe.New Succeeded\n")
	klog.V(6).Infof("transcribe.New LEAVE\n")

	return transcribe, nil
}

func (a *Transcribe) Start() error {
	klog.V(6).Infof("transcribe.Start ENTER\n")

	// call connect!
	wsconn := a.client.Connect()
	if wsconn == nil {
		err := errors.New("client.Connect failed")
		klog.V(1).Infof("transcribe.Start failed", err)
		return err
	}
	klog.V(4).Infof("client.Connect succeeded")

	// start the mic
	err := a.mic.Start()
	if err != nil {
		klog.V(1).Infof("mic.Start failed. Err: %v\n", err)
		klog.V(6).Infof("transcribe.Start LEAVE\n")
		return err
	}
	klog.V(4).Infof("mic.Start succeeded")

	// this is a blocking call
	go func() {
		a.mic.Stream(a.client)
	}()

	klog.V(3).Infof("transcribe.Start Succeeded\n")
	klog.V(6).Infof("transcribe.Start LEAVE\n")
	return nil
}

func (a *Transcribe) Stop() error {
	klog.V(6).Infof("transcribe.Stop ENTER\n")

	// close client
	a.client.Stop()
	klog.V(4).Infof("client.Stop succeeded")

	// close mic stream
	err := a.mic.Stop()
	if err != nil {
		klog.V(1).Infof("mic.Stop failed. Err: %v\n", err)
	}
	klog.V(4).Infof("mic.Stop succeeded")

	// microphone teardown
	klog.V(4).Infof("Calling microphone.Teardown...")
	microphone.Teardown()

	klog.V(6).Infof("transcribe.Stop Succeeded\n")
	klog.V(6).Infof("transcribe.Stop LEAVE\n")
	return nil
}
