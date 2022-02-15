/*
 * Copyright (c) 2021, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package flags

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/spf13/cobra"
	utiltesting "k8s.io/client-go/util/testing"
)

func TestGetServerAddress(t *testing.T) {
	homeDir, err := utiltesting.MkTmpdir("getserveraddress")
	if err != nil {
		t.Fatalf("create tempDir error")
	}

	fakeUserHomeDir := func() (string, error) {
		return homeDir, nil
	}
	patch := monkey.Patch(os.UserHomeDir, fakeUserHomeDir)
	defer patch.Unpatch()

	GetServerAddress()
}

func TestResetFlag(t *testing.T) {
	cmd := &cobra.Command{}
	r := Reset{}
	r.AttachCmd(cmd)
}

func TestDeleteFlag(t *testing.T) {
	cmd := &cobra.Command{}
	d := Delete{}
	d.AttachCmd(cmd)
}

func TestGetFlag(t *testing.T) {
	cmd := &cobra.Command{}
	g := Get{}
	g.AttachCmd(cmd)
}

func TestApplyFlag(t *testing.T) {
	cmd := &cobra.Command{}
	a := Apply{}
	a.AttachCmd(cmd)
}

func TestInstallFlag(t *testing.T) {
	cmd := &cobra.Command{}
	a := Install{}
	a.AttachCmd(cmd)
}
