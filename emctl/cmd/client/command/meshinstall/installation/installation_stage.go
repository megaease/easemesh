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
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
)

// InstallStage holds operations in installation of a stage
type InstallStage interface {
	Do(*installbase.StageContext, Installation) error
	Clear(*installbase.StageContext) error
}

// Installation represents installing or cleaning infrastructure
// components for the EaseMesh
type Installation interface {
	DoInstallStage(*installbase.StageContext) error
	ClearResource(*installbase.StageContext)
}

type installation struct {
	stages []InstallStage
	step   int
}

// New creates a new Installation
func New(stages ...InstallStage) Installation {
	return &installation{stages: stages, step: 0}
}

func (i *installation) DoInstallStage(context *installbase.StageContext) error {
	if i.step >= len(i.stages) {
		return nil
	}
	current := i.step
	i.step++
	return i.stages[current].Do(context, i)
}

func (i *installation) ClearResource(context *installbase.StageContext) {
	for _, f := range context.ClearFuncs {
		err := f(context)
		if err != nil {
			common.OutputErrorf("clear resource error:%s", err)
		}
	}
}

// InstallFunc is the type of install function
type InstallFunc func(*installbase.StageContext) error

// HookFunc is the type of hook function
type HookFunc InstallFunc

// ClearFunc is the type of clean function which cleans installed resources when installation failed
type ClearFunc HookFunc

// PreCheckFunc is the type of function previously checking condition whether is satisfied with the installation
type PreCheckFunc HookFunc

// DescribeFunc is the type of function describing what's the situation of the installation
type DescribeFunc func(*installbase.StageContext, installbase.InstallPhase) string

// Wrap creates new InstallStage via wraping functions
func Wrap(preCheckFunc HookFunc, installFunc InstallFunc, clearFunc HookFunc, description DescribeFunc) InstallStage {
	return &baseInstallStage{preCheck: PreCheckFunc(preCheckFunc), installFunc: installFunc, clearFunc: ClearFunc(clearFunc), description: description}
}

type baseInstallStage struct {
	preCheck    PreCheckFunc
	installFunc InstallFunc
	clearFunc   ClearFunc
	description DescribeFunc
}

var _ InstallStage = &baseInstallStage{}

func (b *baseInstallStage) Do(context *installbase.StageContext, install Installation) error {
	fmt.Printf("%s\n", b.description(context, installbase.BeginPhase))
	if b.preCheck != nil {
		if err := b.preCheck(context); err != nil {
			return errors.Wrap(err, "pre check installation condition failed")
		}
	}
	err := b.installFunc(context)
	context.ClearFuncs = append(context.ClearFuncs, b.clearFunc)
	if err != nil {
		return errors.Wrap(err, "invoke install func")
	}

	fmt.Printf("Install successfully end, following resource are deployed successfully: %s\n", b.description(context, installbase.EndPhase))
	return install.DoInstallStage(context)
}

func (b *baseInstallStage) Clear(context *installbase.StageContext) error {
	if b.clearFunc != nil {
		return b.clearFunc(context)
	}
	return nil
}
