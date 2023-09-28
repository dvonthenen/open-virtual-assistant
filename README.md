# The Open Virtual Assistant Framework

The Open Virtual Assistant Framework is an Alexa-style Virtual Assistant library written in Go which can be installed natively on a desktop or even be placed on a edge device since Go supports cross-compilation onto many different CPU architectures.

The idea is that anyone can create their own Virtual Assistant using this Framework by only writing a small subset of code. Most of the heavy lifting like the Microphone, Speech-to-Text, and Text-to-Speech (Assistant Replies) are already handled for you. You only need to implement triggers to what is being heard by your assistant.

## Motivations

I created the Open Source Virtual Assistant Framework in an effort to:
- prove out the [Symbl.ai Go SDK Project](https://github.com/dvonthenen/symbl-go-sdk) which is used break down the speech (our commands to the Assistant) in unique ways
- showcase something cool when you introduce a SDK that is atypical in an ecosystem

The high-level components of this project are:
- using the open source [PortAudio](http://www.portaudio.com/) to handle the Microphone input of our Assistant
- using the [Symbl Platform](https://platform.symbl.ai/) to handle the human Speech-to-Text aspects of the project
- using Google's Text-to-Speech so that we can give our Assistant a voice

## Requirements

Here are something you are going to need to install on your device/laptop/platform you intent to run this on...

### PortAudio

The Microphone makes use of a [microphone package](https://github.com/dvonthenen/symbl-go-sdk/tree/main/pkg/audio/microphone) contained within the [Symbl.ai Go SDK Project](https://github.com/dvonthenen/symbl-go-sdk) repository. That package makes use of the [PortAudio library](http://www.portaudio.com/) which is a cross-platform open source audio library. If you are on Linux, you can install this library using whatever package manager is available (yum, apt, etc.) on your operating system. If you are on macOS, you can install this library using [brew](https://brew.sh/).

### Symbl Platform Account

The SDK needs to be initialized with your account's credentials `APP_ID` and `APP_SECRET`, which are available in your [Symbl.ai Platform][api-keys]. If you don't have a Symbl.ai Platform account, you can [sign up here][symbl_signup] for free.

You must add your `APP_ID` and `APP_SECRET` to your list of environment variables. We use environment variables because they are easy to configure, support PaaS-style deployments, and work well in containerized environments like Docker and Kubernetes.

```sh
export APP_ID=YOUR-APP-ID-HERE
export APP_SECRET=YOUR-APP-SECRET-HERE
```

### Google Cloud Account

You are also going to need a [Google Cloud account](https://cloud.google.com/text-to-speech) which you can create one for free and get $300 in credits for their Text-to-Speech library. If you already have a Google Cloud account, the cost for using the Text-To-Speech is fractional pennies for converting text or in our case strings to minutes of audio/speech.

## Project Structure

The overall project structure...

### Assistant Framework

The [pkg](https://github.com/dvonthenen/open-virtual-assistant/tree/main/pkg) folder really contains most of the code to do the heavy lifting. You can initialize the Assistant using this library. The only thing that really needs to be provided is a [Symbl Platform configuration](https://github.com/dvonthenen/open-virtual-assistant/blob/main/cmd/assistant/cmd.go#L35-L49) to handle the Speech-To-Text and conversation insights likes Topics, Trackers, and etc.

```
config := &cfginterfaces.StreamingConfig{
    InsightTypes: []string{"topic", "question", "action_item", "follow_up"},
    Config: cfginterfaces.Config{
        MeetingTitle:        "my-meeting",
        ConfidenceThreshold: 0.7,
        SpeechRecognition: cfginterfaces.SpeechRecognition{
            Encoding:        "LINEAR16",
            SampleRateHertz: 16000,
        },
    },
    Speaker: cfginterfaces.Speaker{
        Name:   "Jane Doe",
        UserID: "user@email.com",
    },
}
```

And an implementation of [InsightCallback interface](https://github.com/dvonthenen/symbl-go-sdk/blob/main/pkg/api/streaming/v1/interfaces/interface.go#L6-L13) found in the [Symbl.ai Go SDK Project](https://github.com/dvonthenen/symbl-go-sdk). This interface tells the Framework how to route these conversation insights (ie transcription in the form of sentences, trackers, topics, etc) to your code. Think of this interface as a callback to the action you want to do.

```
type InsightCallback interface {
	RecognitionResultMessage(rr *RecognitionResult) error
	MessageResponseMessage(mr *MessageResponse) error
	InsightResponseMessage(ir *InsightResponse) error
	TopicResponseMessage(tr *TopicResponse) error
	TrackerResponseMessage(tr *TrackerResponse) error
	UnhandledMessage(byMsg []byte) error
}
```

### Example Virtual Assistant

There is an example virtual assistant that is currently located in the [cmd/assistant](https://github.com/dvonthenen/open-virtual-assistant/tree/main/cmd/assistant) directory which provides a simplistic example of a virtual assistant.

> **_IMPORTANT:_** This project structure is subject to change. In order to meet the deadlines of a presentation showcasing this project, certain shortcuts were taken. Future structure change mainly entail moving the example virtual assistant to an `example` folder at the root of the directory.

## Community

This Framework is actively being developed, and we love to hear from you! Please feel free to [create an issue][issues] or [open a pull request][pulls] with your questions, comments, suggestions, and feedback. If you liked our integration guide, please star our repo!

This library is released under the [MIT License][license]


[api-keys]: https://platform.symbl.ai/#/login
[symbl_signup]: https://platform.symbl.ai/signup?utm_source=symbl&utm_medium=blog&utm_campaign=devrel&_ga=2.226597914.683175584.1662998385-1953371422.1659457591&_gl=1*mm3foy*_ga*MTk1MzM3MTQyMi4xNjU5NDU3NTkx*_ga_FN4MP7CES4*MTY2MzEwNDQyNi44Mi4xLjE2NjMxMDQ0MzcuMC4wLjA.
[issues]: https://github.com/dvonthenen/open-virtual-assistant/issues
[pulls]: https://github.com/dvonthenen/open-virtual-assistant/pulls
[license]: LICENSE
