// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import "errors"

// constants...
const (
	CHUNK_SIZE = 1024 * 2
)

var (
	//ErrCommandCreateFailed creating the command failed
	ErrCommandCreateFailed = errors.New("Unable to create the command object")

	// ErrNoActiveTask no active task
	ErrNoActiveTask = errors.New("no active task")
)

var (
	// Naive trigger/activation implementation
	GreetingWords = []string{"hi", "hello", "hey", "hallo", "salut", "bonjour", "hola", "eh", "ey"}
	NameWords     = []string{"kit", "chatgpt", "gpt", "kitt", "kid", "kate", "kent", "kiss"}
)
