package deepgram

import (
	"strings"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"

	microphone "github.com/dvonthenen/open-virtual-assistant/pkg/microphone"
	config "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/config"
)

type InsightOptions struct {
	TranscribeOptions *config.TranscribeOptions
	Microphone        *microphone.Microphone
}

type Insights struct {
	options *InsightOptions

	sb strings.Builder
}

func NewInsightHandler(opts *InsightOptions) *Insights {
	return &Insights{
		options: opts,
	}
}

func (i *Insights) Message(mr *api.MessageResponse) error {
	klog.V(5).Infof("\n\n")
	klog.V(5).Infof("---------------------------------------------\n")
	klog.V(5).Infof("recv: %s\n", mr.Channel.Alternatives[0].Transcript)
	klog.V(5).Infof("is_final: %t\n", mr.IsFinal)
	klog.V(5).Infof("speech_final: %t\n", mr.SpeechFinal)
	klog.V(5).Infof("---------------------------------------------\n")
	klog.V(5).Infof("\n\n")

	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		// klog.V(7).Infof("DEEPGRAM - no transcript")
		return nil
	}

	isFinal := mr.SpeechFinal
	sentence = strings.ToLower(sentence)
	i.sb.WriteString(sentence)

	// // debug
	// klog.V(4).Infof("transcription result: text = %s, final = %t", i.sb.String(), isFinal)

	if !isFinal {
		// klog.V(7).Infof("DEEPGRAM - not final")
		return nil
	}

	// debug
	klog.V(2).Infof("Deepgram transcription: text = %s, final = %t", i.sb.String(), isFinal)

	// perform callback
	i.options.Microphone.Mute()
	(*i.options.TranscribeOptions.Callback).Response(i.sb.String())
	i.options.Microphone.Unmute()

	// clear for new sentence
	i.sb.Reset()

	return nil
}

func (i *Insights) Metadata(md *api.MetadataResponse) error {
	klog.V(3).Infof("\nMetadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	klog.V(3).Infof("Metadata.Channels: %d\n", md.Channels)
	klog.V(3).Infof("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))

	return nil
}

func (i *Insights) Error(er *api.ErrorResponse) error {
	klog.V(1).Infof("\nError.Type: %s\n", er.Type)
	klog.V(1).Infof("Error.Message: %s\n", er.Message)
	klog.V(1).Infof("Error.Description: %s\n\n", er.Description)

	return nil
}
