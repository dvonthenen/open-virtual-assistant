// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	micinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone/interfaces"
)

type AssistantCapabilities interface {
	sdkinterfaces.InsightCallback

	SetMicrophone(micinterfaces.Microphone)
}
