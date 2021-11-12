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

package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/megaease/easemeshctl/cmd/transformer/generator"
)

func main() {

	flag.Parse()

	resourceType := flag.Arg(0)
	if resourceType == "" {
		common.ExitWithErrorf("usage: generator <resourceType> [resource=url]")
	}

	var resourceMappings []string
	if len(flag.Args()) > 1 {
		resourceMappings = flag.Args()[1:]
	}

	var subResource = ""
	if strings.Contains(resourceType, ".") {
		r := strings.Split(resourceType, ".")
		resourceType = r[0]
		subResource = r[1]
	}

	spec := initialSpec(generator.ResourceType(resourceType), subResource, resourceMappings)
	err := generator.New(spec).Accept(generator.NewVisitor(generator.ResourceType(spec.ResourceType)))
	if err != nil {
		common.ExitWithError(err)
	}

	err = spec.Buf.Save(spec.GenerateFileName)
	if err != nil {
		common.ExitWithError(err)
	}
}

func initialSpec(resourceName generator.ResourceType, subResource string, mapping []string) *generator.InterfaceFileSpec {

	spec := generator.InterfaceFileSpec{}
	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	spec.Buf = jen.NewFile(goPackage)
	spec.SourceFile = os.Getenv("GOFILE")
	spec.PkgName = goPackage
	ext := filepath.Ext(spec.SourceFile)
	baseFilename := spec.SourceFile[0 : len(spec.SourceFile)-len(ext)]
	spec.GenerateFileName = "zz_" + baseFilename + "_gen.go"
	spec.ResourceType = resourceName
	spec.ResourceMapping = extractResourceMapping(mapping)
	spec.SubResource = subResource
	return &spec
}

func extractResourceMapping(mappings []string) map[string]string {
	result := make(map[string]string)
	for _, mapping := range mappings {
		r := strings.Split(mapping, "=")
		if len(r) == 2 {
			result[r[0]] = r[1]
		} else {
			result[r[0]] = r[0]
		}
	}
	return result
}
