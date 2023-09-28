// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"encoding/json"
	"fmt"
	"strings"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	micinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone/interfaces"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
	klog "k8s.io/klog/v2"
)

type Insights struct {
	speechOpts *speech.SpeechOpts
	mic        micinterfaces.Microphone
}

func NewInsightHandler(options *speech.SpeechOpts) *Insights {
	return &Insights{
		speechOpts: options,
	}
}

func (i *Insights) SetMicrophone(mic micinterfaces.Microphone) {
	i.mic = mic
}

func (i *Insights) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	if !rr.Message.IsFinal {
		klog.V(6).Infof("Only interested in IsFinal messages\n")
		return nil
	}

	data, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("RecognitionResult json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("RecognitionResultMessage Object DUMP:\n%v\n\n", string(data))
	klog.V(3).Infof("\nMessage:\n%v\n\n", rr.Message.Punctuated.Transcript)
	klog.V(3).Infof("-------------------------------\n\n")

	fmt.Printf("\n\nMessage:\n%v\n\n", rr.Message.Punctuated.Transcript)

	sentences := strings.Split(rr.Message.Punctuated.Transcript, ".!?")
	for _, sentence := range sentences {
		if len(sentence) == 0 {
			klog.V(5).Infof("Skip last empty segment")
			continue
		}
		klog.V(5).Infof("Sentence: %s\n", sentence)
		ExecuteActionOnMatch(i.mic, *i.speechOpts, sentence)
	}

	return nil
}

func (i *Insights) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	data, err := json.Marshal(mr)
	if err != nil {
		fmt.Printf("MessageResponse json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("MessageResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")
	return nil
}

func (i *Insights) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	for _, insight := range ir.Insights {
		switch insight.Type {
		case sdkinterfaces.InsightTypeQuestion:
			return i.HandleQuestion(&insight)
		case sdkinterfaces.InsightTypeFollowUp:
			return i.HandleFollowUp(&insight)
		case sdkinterfaces.InsightTypeActionItem:
			return i.HandleActionItem(&insight)
		default:
			data, err := json.Marshal(ir)
			if err != nil {
				klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
				return err
			}

			klog.V(3).Infof("\n\n-------------------------------\n")
			klog.V(3).Infof("TopicResponseMessage Object DUMP:\n%v\n", string(data))
			klog.V(3).Infof("-------------------------------\n\n")
			return nil
		}
	}

	return nil
}

func (i *Insights) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		fmt.Printf("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("TopicResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")
	return nil
}
func (i *Insights) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		fmt.Printf("TrackerResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("TrackerResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")
	return nil
}

func (i *Insights) UnhandledMessage(byMsg []byte) error {
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("UnhandledMessage Object DUMP:\n%v\n", string(byMsg))
	klog.V(3).Infof("-------------------------------\n\n")
	return nil
}

func (i *Insights) HandleQuestion(insight *sdkinterfaces.Insight) error {
	data, err := json.Marshal(insight)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n")
	klog.V(3).Infof("HandleQuestions: %s\n", data)
	klog.V(3).Infof("\n\n")

	return nil
}

func (i *Insights) HandleActionItem(insight *sdkinterfaces.Insight) error {
	data, err := json.Marshal(insight)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n")
	klog.V(3).Infof("HandleActionItem: %s\n", data)
	klog.V(3).Infof("\n\n")

	return nil
}

func (i *Insights) HandleFollowUp(insight *sdkinterfaces.Insight) error {
	data, err := json.Marshal(insight)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("\n\n")
	klog.V(3).Infof("HandleFollowUp: %s\n", data)
	klog.V(3).Infof("\n\n")

	return nil
}
