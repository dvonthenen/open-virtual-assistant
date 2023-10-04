// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

type ResponseCallback interface {
	Response(sentence string) error
}

type Transcriber interface {
	Start() error
	Stop() error
}
