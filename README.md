# The Open Virtual Assistant Framework

The Open Virtual Assistant Framework is an Google/Siri/Alexa-style Virtual Assistant library written in Go which can be installed natively on a desktop or even be placed on a edge device since Go supports cross-compilation onto many different CPU architectures.

The idea is that anyone can create their own Virtual Assistant using this Framework by only writing a small subset of code. Most of the heavy lifting like the Microphone, Speech-to-Text, and Text-to-Speech (Assistant Replies) are already handled for you. You only need to implement triggers to what is being heard by your assistant.

## Motivations

I created the Open Source Virtual Assistant Framework in an effort to:
- prove out breaking down speech (our commands to the Assistant) in unique ways
- showcase something cool for people who might want to create their own assistants

The high-level components of this project are:
- using the open source [PortAudio](http://www.portaudio.com/) to handle the Microphone input of our Assistant
- using Google's Speech-to-Text to handle the human Speech-to-Text aspects of the project
- using Google's Text-to-Speech so that we can give our Assistant a voice

## Requirements

Here are something you are going to need to install on your device/laptop/platform you intent to run this on...

### PortAudio

The Microphone makes use of a [microphone package](https://github.com/dvonthenen/open-virtual-assistant/tree/main/pkg/audio/microphone) contained within this project. That package makes use of the [PortAudio library](http://www.portaudio.com/) which is a cross-platform open source audio library. If you are on Linux, you can install this library using whatever package manager is available (yum, apt, etc.) on your operating system. If you are on macOS, you can install this library using [brew](https://brew.sh/).

### Google Cloud Account

You are also going to need a [Google Cloud account](https://cloud.google.com/text-to-speech) which you can create one for free and get $300 in credits for their Text-to-Speech library. If you already have a Google Cloud account, the cost for using the Text-To-Speech is fractional pennies for converting text or in our case strings to minutes of audio/speech.

## Project Structure

The overall project structure...

### Assistant Framework

The [pkg](https://github.com/dvonthenen/open-virtual-assistant/tree/main/pkg) folder really contains most of the code to do the heavy lifting. You can initialize the Assistant using this library.

The only thing that really needs to be provided are:
- OpenAI (ChatGPT) API KEY to access ChatGPT
- Google Creds to access their cloud to handle the Speech-To-Text and Text-to-Speech.

And an implementation of `TODO`. This interface tells the Framework how to route the transcription to your code. Think of this interface as a callback to the action you want to do.

```
TODO
```

### Example Virtual Assistant

There is an example virtual assistant that is currently located in the [cmd/assistant](https://github.com/dvonthenen/open-virtual-assistant/tree/main/cmd/assistant) directory which provides a simplistic example of a virtual assistant.

> **_IMPORTANT:_** This project structure is subject to change. In order to meet the deadlines of a presentation showcasing this project, certain shortcuts were taken. Future structure change mainly entail moving the example virtual assistant to an `example` folder at the root of the directory.

## Community

This Framework is actively being developed, and we love to hear from you! Please feel free to [create an issue][issues] or [open a pull request][pulls] with your questions, comments, suggestions, and feedback. If you liked our integration guide, please star our repo!

This library is released under the [MIT License][license]

[issues]: https://github.com/dvonthenen/open-virtual-assistant/issues
[pulls]: https://github.com/dvonthenen/open-virtual-assistant/pulls
[license]: LICENSE
