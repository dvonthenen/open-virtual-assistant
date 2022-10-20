// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistant

import (
	"context"
	"os"
	"os/signal"

	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/assistant/interfaces"
	microphone "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone"
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	cfginterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/client/interfaces"
	klog "k8s.io/klog/v2"
)

type Assistant struct {
	symblClient *symbl.StreamClient
	mic         *microphone.Microphone
}

func New(ctx context.Context, callback interfaces.AssistantCapabilities) (*Assistant, error) {
	klog.V(6).Infof("assistant.New ENTER\n")

	defaultCfg := symbl.GetDefaultConfig()

	assistant, err := NewWithConfig(ctx, defaultCfg, callback)
	if err != nil {
		klog.V(1).Infof("Initialize failed. Err: %v\n", err)
		klog.V(6).Infof("assistant.New LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("assistant.New Succeeded\n")
	klog.V(6).Infof("assistant.New LEAVE\n")

	return assistant, nil
}

func NewWithConfig(ctx context.Context, config *cfginterfaces.StreamingConfig, callback interfaces.AssistantCapabilities) (*Assistant, error) {
	klog.V(6).Infof("assistant.NewWithConfig ENTER\n")

	if ctx == nil {
		ctx = context.Background()
	}

	// mic stuf
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	mic, err := microphone.Initialize(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  float32(config.Config.SpeechRecognition.SampleRateHertz),
	})
	if err != nil {
		klog.V(1).Infof("Initialize failed. Err: %v\n", err)
		klog.V(6).Infof("assistant.NewWithConfig LEAVE\n")
		return nil, err
	}

	callback.SetMicrophone(mic)

	// hook up callbacks to framework
	symblClient, err := symbl.NewStreamClient(ctx, config, callback)
	if err != nil {
		klog.V(1).Infof("NewStreamClientWithDefaults failed. Err: %v\n", err)
		klog.V(6).Infof("assistant.NewWithConfig LEAVE\n")
		return nil, err
	}

	// crete client
	assistant := &Assistant{
		symblClient: symblClient,
		mic:         mic,
	}

	klog.V(3).Infof("assistant.NewWithConfig Succeeded\n")
	klog.V(6).Infof("assistant.NewWithConfig LEAVE\n")

	return assistant, nil
}

func (a *Assistant) Start() error {
	klog.V(6).Infof("assistant.Start ENTER\n")

	// start the mic
	err := a.mic.Start()
	if err != nil {
		klog.V(1).Infof("mic.Start failed. Err: %v\n", err)
		klog.V(6).Infof("assistant.Start LEAVE\n")
		return err
	}

	// this is a blocking call
	a.mic.Stream(a.symblClient)

	klog.V(3).Infof("assistant.Start Succeeded\n")
	klog.V(6).Infof("assistant.Start LEAVE\n")
	return nil
}

func (a *Assistant) Stop() error {
	klog.V(6).Infof("assistant.Stop ENTER\n")

	// close stream
	err := a.mic.Stop()
	if err != nil {
		klog.V(1).Infof("mic.Stop failed. Err: %v\n", err)
		klog.V(6).Infof("assistant.Stop LEAVE\n")
		return err
	}
	microphone.Teardown()

	// close client
	a.symblClient.Stop()

	klog.V(6).Infof("assistant.Stop Succeeded\n")
	klog.V(6).Infof("assistant.Stop LEAVE\n")
	return nil
}
