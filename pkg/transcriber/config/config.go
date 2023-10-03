// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package config

import (
	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/interfaces"
)

type TranscribeOptions struct {
	InputChannels int
	SamplingRate  int

	Callback *interfaces.ResponseCallback
}
