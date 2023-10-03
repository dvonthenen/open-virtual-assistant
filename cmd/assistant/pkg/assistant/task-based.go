// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"context"
	"fmt"

	klog "k8s.io/klog/v2"

	speech "github.com/dvonthenen/open-virtual-assistant/pkg/speech"
)

func ParrotReply(options *speech.SpeechOptions, text string) error {
	ctx := context.Background()

	speechClient, err := speech.New(ctx, options)
	if err != nil {
		klog.V(1).Infof("speech.New failed. Err: %v\n", err)
		return err
	}

	speechClient.Play(ctx, text)
	fmt.Printf("Response:\n%s\n\n", text)
	speechClient.Close()

	return nil
}
