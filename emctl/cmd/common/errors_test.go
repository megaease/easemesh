/*
 * Copyright (c) 2017, MegaEase
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

package common

import (
	"errors"
	"os"
	"testing"

	"bou.ke/monkey"
)

func TestErrors(t *testing.T) {
	fakeExit := func(int) {
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer monkey.Unpatch(patch)
	err := errors.New("an error")
	ExitWithError(err)
	ExitWithErrorf("found an error: %s", err)
	OutputError(err)
	OutputErrorf("found an error: %s", err)

}
