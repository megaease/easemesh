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

package generator

import (
	"go/ast"
	"io"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type (
	// Verb indicated how to access mesh resource
	Verb string
	// ResourceType indicates what type of resource, there are three types of resources Global
	// (tenant, service, ingress), Service (canary, loadbalance, resilience),
	// ServiceSubresource (observabilityTracings, observabilityMetrics, observabilityOutputServer).
	ResourceType string

	// InterfaceFileSpec hold the interface basic information to help visitor to generate concreate struct and method.
	InterfaceFileSpec struct {
		Buf              *jen.File
		Source           []*ast.TypeSpec
		Writer           io.Writer
		PkgName          string
		ResourceType     ResourceType
		SubResources     []string
		GenerateFileName string
		SourceFile       string
		ResourceMapping  map[string]string
		SubResource      string
	}
	generator struct {
		spec   *InterfaceFileSpec
		finder *interfaceFinder
	}

	// Generator is generate code.
	Generator interface {
		Accept(visitor InterfaceVisitor) error
	}

	// InterfaceVisitor generate concret struct and method while traveling the interface.
	InterfaceVisitor interface {
		visitorBegin(imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
		visitorResourceGetterConcreatStruct(name string, spec *InterfaceFileSpec) error
		visitorInterfaceConcreatStruct(name string, spec *InterfaceFileSpec) error
		visitorResourceGetterMethod(name string, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
		visitorIntrefaceMethod(concreateStruct string, verb Verb, method *ast.Field, imports []*ast.ImportSpec, spec *InterfaceFileSpec) error
		visitorEnd(spec *InterfaceFileSpec) error
		onError(e error)
	}
)

const (
	// Get get a specific resource.
	Get Verb = "Get"
	// List query a list of a kind of the mesh resource.
	List Verb = "List"
	// Create create a mesh resource.
	Create Verb = "Create"
	// Delete delete a specific mesh resource.
	Delete Verb = "Delete"
	// Patch update a specific mesh resource.
	Patch Verb = "Patch"

	// Global indicates the resource is global resource.
	Global ResourceType = "Global"
	// Service indicates the resource is service resource.
	Service ResourceType = "Service"

	//CustomResource indicates the resource is the custom resource.
	CustomResource ResourceType = "CustomResource"
)

// New create a code generator
func New(spec *InterfaceFileSpec) Generator {
	return &generator{spec: spec, finder: newInterfaceFinder(spec.SourceFile)}
}

func (g *generator) Accept(visitor InterfaceVisitor) error {
	if visitor == nil {
		return errors.Errorf("illegal argument(s): visitor is required")
	}

	err := g.finder.parseFile()
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "parseFile error")
	}
	err = visitor.visitorBegin(g.finder.imports, g.spec)
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "visitorBegin failed")
	}

	for _, name := range g.finder.typeNameResults() {
		name = lowerFirstChar(name)
		if strings.HasSuffix(name, "Getter") {
			err = visitor.visitorResourceGetterConcreatStruct(name, g.spec)
		} else if strings.HasSuffix(name, "Interface") {
			err = visitor.visitorInterfaceConcreatStruct(name, g.spec)
		}
		if err != nil {
			visitor.onError(err)
			return errors.Wrapf(err, "visitor concrete struct %s", name)
		}
	}

	for _, result := range g.finder.types() {
		name := lowerFirstChar(result.Name.Name)
		if strings.HasSuffix(name, "Getter") {
			if interType, ok := result.Type.(*ast.InterfaceType); ok {
				for _, method := range interType.Methods.List {
					if len(method.Names) < 1 {
						err = errors.Errorf("length of method names should greater than 1 ")
						visitor.onError(err)
						return err
					}
					err = visitor.visitorResourceGetterMethod(method.Names[0].Name, method, g.finder.imports, g.spec)
					if err != nil {
						visitor.onError(err)
						return errors.Wrapf(err, "visitorResourceGetterMethod name:%s, method:%+v error", method.Names[0].Name, *method)
					}
				}
			}
		} else if strings.HasSuffix(name, "Interface") {
			if interType, ok := result.Type.(*ast.InterfaceType); ok {
				for _, method := range interType.Methods.List {
					if len(method.Names) < 1 {
						err = errors.Errorf("length of method names should greater than 1 ")
						visitor.onError(err)
						return err
					}
					err = visitor.visitorIntrefaceMethod(name, Verb(method.Names[0].Name), method, g.finder.imports, g.spec)
					if err != nil {
						visitor.onError(err)
						return errors.Wrapf(err, "visitorInterfaceMethod name:%s, method:%+v error", method.Names[0].Name, *method)
					}
				}
			}

		}
	}

	err = visitor.visitorEnd(g.spec)
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "visitorEnd failed")
	}
	return nil
}

func lowerFirstChar(str string) string {
	first := true
	return strings.Map(func(r rune) rune {
		if !first {
			return r
		}
		first = false
		return unicode.ToLower(r)
	}, str)
}
