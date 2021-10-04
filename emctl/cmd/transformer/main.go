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
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/megaease/easemeshctl/cmd/transformer/generator"
)

func main() {

	flag.Parse()

	resourceName := flag.Arg(0)
	if resourceName == "" {
		common.ExitWithErrorf("usage: generator <resourceName>")
	}

	spec := initialSpec(resourceName)
	err := generator.New(spec).Accept(generator.NewVisitor(spec.ResourceName))
	if err != nil {
		common.ExitWithError(err)
	}
	spec.Buf.Render(os.Stdout)
	err = spec.Buf.Save(spec.GenerateFileName)
	if err != nil {
		common.ExitWithError(err)
	}
}

func initialSpec(resourceName string) *generator.InterfaceFileSpec {

	spec := generator.InterfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	spec.Buf = jen.NewFile(goPackage)
	spec.SourceFile = os.Getenv("GOFILE")
	spec.PkgName = goPackage
	ext := filepath.Ext(spec.SourceFile)
	baseFilename := spec.SourceFile[0 : len(spec.SourceFile)-len(ext)]
	spec.GenerateFileName = baseFilename + "_gen.go"
	spec.ResourceName = resourceName
	return &spec
}
