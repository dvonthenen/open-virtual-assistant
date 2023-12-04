// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package google

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	speechtotext "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	microphone "github.com/dvonthenen/open-virtual-assistant/pkg/microphone"
	"github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
)

const (
	DefaultLanguage = "en-US"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)

type Transcribe struct {
	options *config.TranscribeOptions

	googleClient      *speechtotext.Client
	client            speechpb.Speech_StreamingRecognizeClient
	googleCredentials string

	ctx       context.Context
	ctxCancel context.CancelFunc

	mic *microphone.Microphone
}

var micInitAlready = false

func New(ctx context.Context, opts *config.TranscribeOptions) (*Transcribe, error) {
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

	// google speech to text
	var googleCredentials string
	if v := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); v != "" {
		klog.V(4).Info("GOOGLE_APPLICATION_CREDENTIALS found")
		googleCredentials = v
	} else {
		klog.Error("GOOGLE_APPLICATION_CREDENTIALS not found")
		klog.V(1).Infof("speech.New LEAVE\n")
		return nil, ErrInvalidInput
	}

	googleClient, err := speechtotext.NewClient(ctx)
	if err != nil {
		klog.V(1).Infof("speechtotext.NewClient failed. Err: %v\n", err)
		return nil, err
	}

	client, err := googleClient.StreamingRecognize(ctx)
	if err != nil {
		klog.V(1).Infof("client.StreamingRecognize failed. Err: %v\n", err)
		return nil, err
	}

	t := &Transcribe{
		options:           opts,
		ctx:               ctx,
		googleClient:      googleClient,
		client:            client,
		googleCredentials: googleCredentials,
		mic:               mic,
	}
	t.ctx, t.ctxCancel = context.WithCancel(ctx)

	return t, nil
}

func (t *Transcribe) Start() error {
	klog.V(6).Infof("transcribe.Start ENTER\n")

	// call connect!
	err := t.connect()
	if err != nil {
		err := errors.New("client.Connect failed")
		klog.V(1).Infof("transcribe.Start failed", err)
		return err
	}
	klog.V(4).Infof("client.Connect succeeded")

	// start the mic
	err = t.mic.Start()
	if err != nil {
		klog.V(1).Infof("mic.Start failed. Err: %v\n", err)
		klog.V(6).Infof("transcribe.Start LEAVE\n")
		return err
	}
	klog.V(4).Infof("mic.Start succeeded")

	// this is a blocking call
	go func() {
		t.mic.Stream(t)
	}()

	klog.V(3).Infof("transcribe.Start Succeeded\n")
	klog.V(6).Infof("transcribe.Start LEAVE\n")
	return nil
}

func (t *Transcribe) connect() error {
	klog.V(5).Infof("calling Transcribe.connect")

	config := &speechpb.RecognitionConfig{
		Model: "command_and_search",
		Adaptation: &speechpb.SpeechAdaptation{
			PhraseSets: []*speechpb.PhraseSet{
				{
					Phrases: []*speechpb.PhraseSet_Phrase{
						{Value: "${hello} ${gpt}"},
						{Value: "${gpt}"},
						{Value: "Hey ${gpt}"},
						{Value: "Kitt"},
						{Value: "Kit-t"},
						{Value: "Kit"},
						{Value: "${action} a task ${named}"},
						{Value: "${action} the task ${named}"},
						{Value: "${action} a job ${named}"},
						{Value: "${action} the job ${named}"},
						{Value: "task"},
						{Value: "job"},
						{Value: "${action}"},
						{Value: "${named}"},
					},
					Boost: 16,
				},
			},
			CustomClasses: []*speechpb.CustomClass{
				{
					CustomClassId: "hello",
					Items: []*speechpb.CustomClass_ClassItem{
						{Value: "Hi"},
						{Value: "Hello"},
						{Value: "Hey"},
					},
				},
				{
					CustomClassId: "gpt",
					Items: []*speechpb.CustomClass_ClassItem{
						{Value: "Kit"},
						{Value: "KITT"},
						{Value: "GPT"},
					},
				},
				{
					CustomClassId: "action",
					Items: []*speechpb.CustomClass_ClassItem{
						{Value: "create"},
						{Value: "activate"},
						{Value: "resume"},
					},
				},
				{
					CustomClassId: "named",
					Items: []*speechpb.CustomClass_ClassItem{
						{Value: "name"},
						{Value: "named"},
						{Value: "called"},
					},
				},
			},
		},
		UseEnhanced:       true,
		Encoding:          speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz:   int32(t.options.SamplingRate),
		AudioChannelCount: int32(t.options.InputChannels),
		LanguageCode:      DefaultLanguage,
	}

	if err := t.client.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				InterimResults: true,
				Config:         config,
			},
		},
	}); err != nil {
		klog.V(1).Infof("client.Send failed. Err: %v\n", err)
		return err
	}

	// kick off threads
	go t.listen()

	klog.V(6).Infof("new speech stream created successfully")
	return nil
}

func (t *Transcribe) listen() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			for {
				resp, err := t.client.Recv()
				if err != nil {
					klog.V(1).Infof("client.Recv failed. Err: %v\n", err)

					if status, ok := status.FromError(err); ok {
						if status.Code() == codes.OutOfRange {
							klog.V(1).Infof("client recognize out of range")
							break
						} else if status.Code() == codes.Canceled {
							klog.V(1).Infof("client recognize canceled")
							break
						}
					}
				}

				if resp.Error != nil {
					klog.V(1).Infof("client.Recv failed. resp.Error: ", resp.Error)
					break
				}

				// Read the whole transcription and put inside one string
				// We don't need to process each part individually (atm?)
				var sb strings.Builder
				for _, result := range resp.Results {
					alt := result.Alternatives[0]
					text := alt.Transcript

					// apepend to string builder
					sb.WriteString(text)

					if !result.IsFinal {
						// klog.V(4).Infof("isFinal = FALSE")
						continue
					}

					// Debug... what is being said word for word
					sentence := sb.String()
					klog.V(2).Infof("google transcription: text=%s final=%t\n", sentence, result.IsFinal)

					if t.options.Callback != nil {
						t.mic.Mute()
						(*t.options.Callback).Response(sentence)
						t.mic.Unmute()
					} else {
						klog.V(2).Infof("stream.Recv() text=%s final=%t\n", sentence, result.IsFinal)
					}
					sb.Reset()
				}
			}
		}
	}
}

// Write performs the lower level write operation
func (t *Transcribe) Write(buf []byte) (int, error) {
	if err := t.client.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			AudioContent: buf,
		},
	}); err != nil {
		klog.V(1).Infof("stream.Send failed. Err: %v\n", err)
		return 0, err
	}

	return len(buf), nil
}

func (t *Transcribe) Stop() error {
	klog.V(5).Infof("calling Transcribe.Stop")

	// google client
	t.googleClient.Close()

	// close mic stream
	err := t.mic.Stop()
	if err != nil {
		klog.V(1).Infof("mic.Stop failed. Err: %v\n", err)
	}
	klog.V(4).Infof("mic.Stop succeeded")

	// microphone teardown
	klog.V(4).Infof("Calling microphone.Teardown...")
	microphone.Teardown()

	return nil
}
