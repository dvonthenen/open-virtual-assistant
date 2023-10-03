// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	matchr "github.com/antzucaro/matchr"
	"golang.org/x/exp/slices"
	klog "k8s.io/klog/v2"

	personas "github.com/dvonthenen/chat-gpeasy/pkg/personas"
	gpeasyinterfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"

	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/interfaces"
	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
	config "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
	gtranscriber "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/google"
)

// constants...
const (
	// transcriber options
	GOOGLE_TRANSCRIBER string = "google"
)

var (
	// Naive trigger/activation implementation
	GreetingWords = []string{"hi", "hello", "hey", "hallo", "salut", "bonjour", "hola", "eh", "ey"}
	NameWords     = []string{"kit", "chatgpt", "gpt", "kitt", "kid", "kate", "kent", "kiss"}
)

// assistant implementation
type AssistantOptions struct {
	TranscribeOptions *config.TranscribeOptions
	SpeechOptions     *speech.SpeechOptions

	OldDemo bool
}

type Assistant struct {
	options *AssistantOptions

	speech      *speech.Client
	transcriber *interfaces.Transcriber
	tasks       map[string]*gpeasyinterfaces.AdvancedChatStream
	jobs        map[string]*gpeasyinterfaces.AdvancedChatStream
	activeTask  *gpeasyinterfaces.AdvancedChatStream
	activeJob   *gpeasyinterfaces.AdvancedChatStream
}

func New(opts *AssistantOptions) (*Assistant, error) {
	ctx := context.Background()

	if opts.SpeechOptions == nil {
		opts.SpeechOptions = &speech.SpeechOptions{}
	}
	if opts.TranscribeOptions == nil {
		opts.TranscribeOptions = &config.TranscribeOptions{}
	}

	// assistant
	assistant := &Assistant{
		options: opts,
		tasks:   make(map[string]*gpeasyinterfaces.AdvancedChatStream),
		jobs:    make(map[string]*gpeasyinterfaces.AdvancedChatStream),
	}

	// speech stuff
	speech, err := speech.New(ctx, opts.SpeechOptions)
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

	// set callback for transcriber
	var cb interfaces.ResponseCallback
	cb = assistant

	opts.TranscribeOptions.Callback = &cb

	// transcriber
	var transcriber interfaces.Transcriber

	switch transcriberStr {
	case GOOGLE_TRANSCRIBER:
		fallthrough
	default:
		transcribe, errTranscribe := gtranscriber.New(ctx, opts.TranscribeOptions)
		if errTranscribe != nil {
			klog.V(1).Infof("gtranscriber.New failed. Err: %v\n", errTranscribe)
			return nil, errTranscribe
		}
		transcriber = transcribe
	}

	// housekeeping
	assistant.speech = speech
	assistant.transcriber = &transcriber

	return assistant, nil
}

func (a *Assistant) Start() error {
	err := (*a.transcriber).Start()
	if err != nil {
		klog.V(1).Infof("transcriber.Start failed. Err: %v\n", err)
	}
	return err
}

func (a *Assistant) Response(text string) error {
	text = strings.ToLower(text)
	klog.V(5).Infof("text: %s\n", text)

	if a.options.OldDemo {
		/*
			Old Simple Demo
		*/
		for _, triggers := range QuestionResponses {
			for _, key := range triggers.keys {
				klog.V(5).Infof("Key: %s\n", key)

				if percent := matchr.JaroWinkler(key, text, false); percent > 0.95 {
					klog.V(5).Infof("\n\n--------------------------------\n")
					klog.V(3).Infof("MATCH (%f): %s = %s\n", percent, key, text)
					klog.V(5).Infof("\n--------------------------------\n\n")

					klog.V(2).Infof("Heard:\nMATCH (%f): %s = %s\n\n", percent, key, text)

					err := triggers.callback(a.options.SpeechOptions, text)
					if err != nil {
						klog.V(1).Infof("speech.Play failed. Err: %v\n", err)
					}

					// exit on first match!
					return nil
				}
			}
		}
	} else {
		/*
			Task based assistant demo
		*/

		// Check if text contains at least one Greet and Name Words founds
		words := strings.Split(text, " ")

		if len(words) < 3 {
			klog.V(4).Infof("Not enough words to process. Skipping...\n")
			return nil
		}

		foundGreet := ""
		for _, greet := range GreetingWords {
			klog.V(6).Infof("checking grett word = %s\n", greet)
			if greetIndex := slices.Index(words, greet); greetIndex != -1 && greetIndex < 5 {
				klog.V(4).Infof("greeting word FOUND = %s\n", greet)
				foundGreet = greet
				break
			}
		}

		foundName := ""
		for _, name := range NameWords {
			klog.V(6).Infof("checking greet word = %s\n", name)
			if nameIndex := slices.Index(words, name); nameIndex != -1 && nameIndex < 5 {
				klog.V(4).Infof("name word FOUND = %s\n", name)
				foundName = name
				break
			}
		}

		// If both found, activate kitt
		if foundGreet != "" && foundName != "" {
			klog.V(2).Infof("Greeting=%s and Name=%s found. Asking kitt.\n", foundGreet, foundName)

			// activate task?
			regTask, err := regexp.Compile("(activate|create|resume)\\s(a|the)??\\stask\\s(name|named|called)\\s{1}([a-z\\s]+)")
			if err != nil {
				klog.V(1).Infof("regexp.Compile failed. Err: %v\n", err)
				return err
			}
			regJob, err := regexp.Compile("(create)\\s(a|the)??\\sjob\\s(name|named|called)\\s{1}([a-z\\s]+)")
			if err != nil {
				klog.V(1).Infof("regexp.Compile failed. Err: %v\n", err)
				return err
			}

			if regJob.MatchString(text) {
				klog.V(2).Infof("Creating/activating a job...\n")

				matches := regJob.FindStringSubmatch(text)
				jobAction := matches[1]
				jobName := matches[4]

				klog.V(2).Infof("jobAction: %s, jobName: %s\n", jobAction, jobName)

				// check if task already exists
				a.activeTask = a.jobs[jobName]
				if a.activeTask == nil {
					// create chatgpt client
					personaConfig, err := personas.DefaultConfig("", "")
					if err != nil {
						klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
						return err
					}

					persona, err := personas.NewAdvancedChatStreamWithOptions(personaConfig)
					if err != nil {
						klog.V(1).Infof("personas.NewAdvancedChatStreamWithOptions failed. Err: %v\n", err)
						return err
					}

					(*persona).Init(gpeasyinterfaces.SkillTypeGeneric, "")
					err = (*persona).AddDirective("Try using the information provided in this conversation thread before going to other sources when answering question.")
					if err != nil {
						klog.V(1).Infof("personas.AddDirective failed error: %v\n", err)
					}

					a.jobs[jobName] = persona
					a.activeJob = persona
				}

				// clear active task
				a.activeTask = nil

				err = ParrotReply(a.options.SpeechOptions, fmt.Sprintf("The job called %s has been %sd. What would you like me to research?", jobName, jobAction))
				if err != nil {
					klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
					return err
				}

				return nil
			} else if regTask.MatchString(text) {
				klog.V(2).Infof("Creating/activating a task...\n")

				matches := regTask.FindStringSubmatch(text)
				taskAction := matches[1]
				taskName := matches[4]

				klog.V(2).Infof("taskAction: %s, taskName: %s\n", taskAction, taskName)

				// check if task already exists
				a.activeTask = a.tasks[taskName]
				if a.activeTask == nil {
					// create chatgpt client
					personaConfig, err := personas.DefaultConfig("", "")
					if err != nil {
						klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
						return err
					}

					persona, err := personas.NewAdvancedChatStreamWithOptions(personaConfig)
					if err != nil {
						klog.V(1).Infof("personas.NewAdvancedChatStreamWithOptions failed. Err: %v\n", err)
						return err
					}

					(*persona).Init(gpeasyinterfaces.SkillTypeGeneric, "")
					err = (*persona).AddDirective("Try using the information provided in this conversation thread before going to other sources when answering question.")
					if err != nil {
						klog.V(1).Infof("personas.AddDirective failed error: %v\n", err)
					}

					a.tasks[taskName] = persona
					a.activeTask = persona
				}

				// clear active job
				a.activeJob = nil

				err = ParrotReply(a.options.SpeechOptions, fmt.Sprintf("The task called %s has been %sd.", taskName, taskAction))
				if err != nil {
					klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
					return err
				}

				return nil
			}

			// activate task or job?
			regActionSkip, err := regexp.Compile("(activate|create|resume)+")
			if err != nil {
				klog.V(1).Infof("regexp.Compile failed. Err: %v\n", err)
				return err
			}
			regItemSkip, err := regexp.Compile("(credit|task|job)+")
			if err != nil {
				klog.V(1).Infof("regexp.Compile failed. Err: %v\n", err)
				return err
			}

			if regActionSkip.MatchString(text) || regItemSkip.MatchString(text) {
				klog.V(2).Infof("Skip using KITT...\n")
				return nil
			}

			// active task?
			if a.activeTask != nil {
				klog.V(2).Infof("Active task found. Asking kitt.\n")

				err := a.activetaskQuestion(text)
				if err != nil {
					klog.V(1).Infof("activetaskQuestion failed. Err: %v\n", err)
				} else {
					klog.V(4).Infof("activetaskQuestion succeeded. text: %s\n", text)
				}

				return nil
			}

			// throwaway but need to answer
			klog.V(2).Infof("No active task found. Creating a throwaway.\n")

			err = a.throwawayQuestion(text)
			if err != nil {
				klog.V(1).Infof("throwawayQuestion failed. Err: %v\n", err)
			}
			return err
		} else if a.activeJob != nil {
			klog.V(2).Infof("This is not a message for Kitt. This is the start to launching a job.\n")

			// TODO: commenting this out for demo purposes since this is SUPER long running
			// cmdline := fmt.Sprintf("python ./babyagi/babyapi.py", text)

			// cmdline := fmt.Sprintf("./llama.cpp/main -m models/llama-2-7b-chat.Q4_K_M.gguf --color -c 4096 --temp 0.7 --repeat_penalty 1.1 -n -1 -p \"[INST] <<SYS>>You are a helpful, respectful and honest assistant. Always answer as helpfully as possible, while being safe.  Your answers should not include any harmful, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature. If a question does not make any sense, or is not factually coherent, explain why instead of answering something not correct. If you dont know the answer to a question, please dont share false information.<</SYS>>%s[/INST]\" 2>/dev/null",
			// 	text)

			// stopChan := make(chan struct{})
			// err := command(cmdline, a.speech, stopChan)
			// if err != nil {
			// 	klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
			// 	return err
			// }

			err := ParrotReply(a.options.SpeechOptions, fmt.Sprintf("Launching long running job will report back when finished. Prompt: %s", text))
			if err != nil {
				klog.V(1).Infof("personas.DefaultConfig error: %v\n", err)
				return err
			}

			// TODO: commenting this out for demo purposes
			a.activeJob = nil

		} else if a.activeTask != nil {
			klog.V(2).Infof("This is not a message for Kitt. Adding to activate task.\n")

			err := (*a.activeTask).AddUserContext(text)
			if err != nil {
				klog.V(1).Infof("activeTask.AddUserContext failed. Err: %v\n", err)
			}

			return err
		}
		// }
	}

	return nil
}

func (a *Assistant) throwawayQuestion(text string) error {
	// create chatgpt client
	personaConfig, err := personas.DefaultConfig("", "")
	if err != nil {
		klog.V(1).Infof("personas.DefaultConfig failed. Err: %v\n", err)
		return err
	}

	persona, err := personas.NewAdvancedChatStreamWithOptions(personaConfig)
	if err != nil {
		klog.V(1).Infof("personas.NewAdvancedChatStreamWithOptions failed. Err: %v\n", err)
		return err
	}

	(*persona).Init(gpeasyinterfaces.SkillTypeGeneric, "")

	stream, err := (*persona).Query(context.Background(), text)
	if err != nil {
		klog.V(1).Infof("personas.Query failed. Err: %v\n", err)
		return err
	}

	// convert stream to string
	sb := bytes.NewBufferString("")

	err = (*stream).Stream(sb)
	if err != nil {
		klog.V(1).Infof("stream.Stream failed. Err: %v\n", err)
		return err
	}
	(*stream).Close()

	trimSentence := strings.TrimSpace(sb.String())

	err = ParrotReply(a.options.SpeechOptions, trimSentence)
	if err != nil {
		klog.V(1).Infof("stream.Stream failed. Err: %v\n", err)
		return err
	}

	klog.V(4).Infof("throwawayQuestion succeeded. text: %s\n", trimSentence)
	return nil
}

func (a *Assistant) activetaskQuestion(text string) error {
	text = strings.TrimSpace(text)

	if a.activeTask == nil {
		klog.V(1).Infof("personas.Query failed\n")
		return ErrNoActiveTask
	}

	stream, err := (*a.activeTask).Query(context.Background(), text)
	if err != nil {
		klog.V(1).Infof("personas.Query failed. Err: %v\n", err)
		return err
	}

	// convert stream to string
	sb := bytes.NewBufferString("")

	err = (*stream).Stream(sb)
	if err != nil {
		klog.V(1).Infof("stream.Stream failed. Err: %v\n", err)
		return err
	}
	(*stream).Close()

	trimSentence := strings.TrimSpace(sb.String())

	err = ParrotReply(a.options.SpeechOptions, trimSentence)
	if err != nil {
		klog.V(1).Infof("stream.Stream failed. Err: %v\n", err)
		return err
	}

	klog.V(4).Infof("throwawayQuestion succeeded. text: %s\n", trimSentence)
	return nil
}

func (a *Assistant) Stop() error {
	err := (*a.transcriber).Stop()
	if err != nil {
		klog.V(1).Infof("transcriber.Stop failed. Err: %v\n", err)
	}
	return err
}
