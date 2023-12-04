// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

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
