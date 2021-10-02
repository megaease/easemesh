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

package main

import (
	"flag"
	"go/ast"
	"io"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/megaease/easemeshctl/cmd/common"
)

type interfaceFileSpec struct {
	buf              *jen.File
	source           []*ast.TypeSpec
	writer           io.Writer
	pkgName          string
	resourceName     string
	generateFileName string
	sourceFile       string
}

func main() {

	flag.Parse()

	resourceName := flag.Arg(0)
	var spec *interfaceFileSpec
	if len(flag.Args()) > 1 {
		spec = fromArgs(flag.Args())
	} else if resourceName != "" {
		spec = initialSpec(resourceName)
	} else {
		spec = fromArgs(flag.Args())
		//common.ExitWithErrorf("usage: generator <resourceName>")
	}

	err := newGenerator(spec).accept(&meshClientVisitor{resourceName: spec.resourceName, builder: &interfaceBuilder{}})
	if err != nil {
		common.ExitWithError(err)
	}
	spec.buf.Render(os.Stdout)
	err = spec.buf.Save(spec.generateFileName)
	if err != nil {
		common.ExitWithError(err)
	}
}

func initialSpec(resourceName string) *interfaceFileSpec {

	spec := interfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	spec.buf = jen.NewFile(goPackage)
	spec.sourceFile = os.Getenv("GOFILE")
	spec.pkgName = goPackage
	ext := filepath.Ext(spec.sourceFile)
	baseFilename := spec.sourceFile[0 : len(spec.sourceFile)-len(ext)]
	spec.generateFileName = baseFilename + "_gen.go"
	spec.resourceName = resourceName
	return &spec
}

func fromArgs(args []string) *interfaceFileSpec {
	// for test
	spec := interfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := "github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	spec.buf = jen.NewFile(goPackage)
	spec.sourceFile = "../client/command/meshclient/canary_interface.go"
	spec.pkgName = goPackage
	ext := filepath.Ext(spec.sourceFile)
	baseFilename := spec.sourceFile[0 : len(spec.sourceFile)-len(ext)]
	spec.generateFileName = baseFilename + "_gen.go"
	spec.resourceName = "canary"
	return &spec
}
