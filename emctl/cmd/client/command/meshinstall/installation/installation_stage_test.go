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
package installation

import (
	"testing"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
)

func stepOneDescribe(*installbase.StageContext, installbase.InstallPhase) string {
	return "stepOneDescribe"

}

func stepOneClear(s *installbase.StageContext) error {
	return nil
}

func stepOneDeploy(s *installbase.StageContext) error {
	return nil
}

func stepOnePreCheck(s *installbase.StageContext) error {
	return nil
}

func stepTwoDescribe(*installbase.StageContext, installbase.InstallPhase) string {
	return "stepTwoDescribe"

}

func stepTwoClear(s *installbase.StageContext) error {
	return nil
}

func stepTwoDeploy(s *installbase.StageContext) error {
	return nil
}

func stepTwoPreCheck(s *installbase.StageContext) error {
	return nil
}

func TestInstallation(t *testing.T) {

	installStages := []InstallStage{
		Wrap(stepOnePreCheck, stepOneDeploy, stepOneClear, stepOneDescribe),
		Wrap(stepTwoPreCheck, stepTwoDeploy, stepTwoClear, stepTwoDescribe),
	}

	installations := New(installStages...)

	err := installations.DoInstallStage(&installbase.StageContext{})
	if err != nil {
		t.Fatalf("Run mock installage failed %s", err)
	}

	installContext := installbase.StageContext{
		ClearFuncs: []func(*installbase.StageContext) error{
			stepOneClear,
			stepTwoClear,
		},
	}

	installations.ClearResource(&installContext)

	installStages[0].Clear(&installContext)
}
