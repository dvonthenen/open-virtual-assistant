// Copyright 2022. All Rights Reserved.
// SPDX-License-Identifier: MIT

package assistantimpl

import (
	"errors"
	"io"
	"os/exec"

	klog "k8s.io/klog/v2"
)

const (
	CHUNK_SIZE = 1024 * 2
)

var (
	//ErrBufferCreateFailed creating the buffer failed
	ErrBufferCreateFailed = errors.New("Unable to create the buffer object")

	//ErrScannerCreateFailed creating the scanner failed
	ErrScannerCreateFailed = errors.New("Unable to create the scanner object")

	//ErrReaderCreateFailed creating the reader failed
	ErrReaderCreateFailed = errors.New("Unable to create the reader object")

	//ErrCommandCreateFailed creating the command failed
	ErrCommandCreateFailed = errors.New("Unable to create the command object")

	//ErrExecuteFailed installation package failed
	ErrExecuteFailed = errors.New("The command line failed to execute correctly")
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

	klog.V(1).Infof("Cmdline implicitly failed to execute correctly")
	klog.V(6).Infof("command LEAVE\n")
	return ErrExecuteFailed
}
