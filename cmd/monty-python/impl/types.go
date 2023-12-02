// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import (
	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
)

// handles question...
type ResponseFunc func() string

type QuestionResponse struct {
	keys     []string
	callback ResponseFunc
}

// My Assistant
type MyAssistant struct {
	speech *interfaces.Speech
}
