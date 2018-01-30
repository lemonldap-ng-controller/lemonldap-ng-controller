/*
Copyright 2018 Mathieu Parent <math.parent@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/golang/glog"
)

// StartProcess starts a new LemonLDAP::NG master process running in foreground
func (c *LemonLDAPNGController) StartProcess(stopCh <-chan struct{}) {
	cmd := exec.Command(c.controllerConfig.Command[0], c.controllerConfig.Command[1:]...)

	// put llng-fastcgi-server in another process group to prevent it
	// to receive signals meant for the controller
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	glog.Info("Starting LemonLDAP::NG process...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		glog.Fatalf("LemonLDAP::NG error: %v", err)
		c.llngErrCh <- err
		return
	}

	go func() {
		c.llngErrCh <- cmd.Wait()
	}()

}
