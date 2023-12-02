// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
	transcriber "github.com/dvonthenen/open-virtual-assistant/pkg/transcriber/interfaces"
)

type AssistantImpl interface {
	transcriber.ResponseCallback

	SetSpeech(s *speech.Speech)
}
