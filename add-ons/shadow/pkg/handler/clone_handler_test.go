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

package handler

import (
	"testing"

	shadowfake "github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
)

func TestShadowServiceCloner_Clone(t *testing.T) {
	cloner := &ShadowServiceCloner{
		KubeClient: prepareClientForTest(),
	}

	shadowService := shadowfake.NewShadowService()
	sourceDeployment := shadowfake.NewSourceDeployment()

	serviceCloneBlock := ShadowServiceBlock{
		service:   shadowService,
		deployObj: sourceDeployment,
	}
	cloner.Clone(serviceCloneBlock)
}
