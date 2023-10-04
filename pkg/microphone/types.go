package microphone

import (
	"sync"

	"github.com/gordonklaus/portaudio"
)

// AudioConfig init config for library
type AudioConfig struct {
	InputChannels int
	SamplingRate  float32
}

// Microphone...
type Microphone struct {
	// microphone
	stream *portaudio.Stream

	// buffer
	intBuf []int16

	// operational
	stopChan chan struct{}
	mute     sync.Mutex
	muted    bool
}
