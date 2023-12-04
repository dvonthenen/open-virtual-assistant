// Copyright 2023 The dvonthenen Open-Virtual-Assistant Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package impl

import (
	"io"
	"os/exec"

	klog "k8s.io/klog/v2"
)

func command(cmdLine string, w io.Writer, stopChan chan struct{}) error {
	klog.V(6).Infof("command ENTER\n")
	klog.V(2).Infof("Cmdline: %s\n", cmdLine)

	cmd := exec.Command("bash", "-c", cmdLine)
	if cmd == nil {
		klog.V(1).Infof("Error creating cmd\n")
		klog.V(6).Infof("command LEAVE\n")
		return ErrCommandCreateFailed
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		klog.V(1).Infof("StdoutPipe failed. Err: %v\n", err)
		klog.V(6).Infof("command LEAVE\n")
		return err
	}

	chunk := make([]byte, CHUNK_SIZE)
	for {
		select {
		case <-stopChan:
			return nil
		default:

			bytesRead, err := out.Read(chunk)
			if err != nil {
				klog.V(1).Infof("out.Read failed. Err: %v\n", err)
				return err
			}

			if bytesRead == 0 {
				continue
			}

			_, err = w.Write(chunk[:bytesRead])
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				return err
			}
			klog.V(7).Infof("io.Writer succeeded. Bytes written: %d\n", bytesRead)
		}
	}
}
