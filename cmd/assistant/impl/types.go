// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import (
	gpeasyinterfaces "github.com/dvonthenen/chat-gpeasy/pkg/personas/interfaces"

	interfaces "github.com/dvonthenen/open-virtual-assistant/pkg/speech/interfaces"
)

type MyAssistant struct {
	speech *interfaces.Speech

	tasks      map[string]*gpeasyinterfaces.AdvancedChatStream
	jobs       map[string]*gpeasyinterfaces.AdvancedChatStream
	activeTask *gpeasyinterfaces.AdvancedChatStream
	activeJob  *gpeasyinterfaces.AdvancedChatStream
}
