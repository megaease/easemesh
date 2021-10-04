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
	"fmt"
	"go/ast"
	"io"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

// Verb indicated how to access mesh resource
type Verb string

const (
	// Get get a specific resource
	Get Verb = "Get"
	// List query a list of a kind of the mesh resource
	List Verb = "List"
	// Create create a mesh resource
	Create Verb = "Create"
	// Delete delete a specific mesh resource
	Delete Verb = "Delete"
	// Patch update a specific mesh resource
	Patch Verb = "Patch"
)

type (
	// InterfaceFileSpec hold the interface basic information to help visitor to generate concreate struct and method
	InterfaceFileSpec struct {
		Buf              *jen.File
		Source           []*ast.TypeSpec
		Writer           io.Writer
		PkgName          string
		ResourceName     string
		GenerateFileName string
		SourceFile       string
	}
	generator struct {
		spec   *InterfaceFileSpec
		finder *interfaceFinder
	}

	// Generator is generate code
	Generator interface {
		Accept(visitor InterfaceVisitor) error
	}

	// InterfaceVisitor generate concret struct and method while traveling the interface
	InterfaceVisitor interface {
		visitorBegin(buf *jen.File, imports []*ast.ImportSpec) error
		visitorResourceGetterConcreatStruct(name string, buf *jen.File) error
		visitorInterfaceConcreatStruct(name string, buf *jen.File) error
		visitorResourceGetterMethod(name string, method *ast.Field, imports []*ast.ImportSpec, buf *jen.File) error
		visitorIntrefaceMethod(verb Verb, method *ast.Field, imports []*ast.ImportSpec, buf *jen.File) error
		visitorEnd(buf *jen.File) error
		onError(e error)
	}
)

// New create a code generator
func New(spec *InterfaceFileSpec) Generator {
	return &generator{spec: spec, finder: newInterfaceFinder(spec.SourceFile)}
}

func (g *generator) Accept(visitor InterfaceVisitor) error {
	if visitor == nil {
		return errors.Errorf("illegal argument(s): visitor is required")
	}

	fmt.Printf("visitor: %+v\n", visitor)

	err := g.finder.parseFile()
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "parseFile error")
	}
	err = visitor.visitorBegin(g.spec.Buf, g.finder.imports)
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "visitorBegin failed")
	}

	for _, name := range g.finder.typeNameResults() {
		name = lowerFirstChar(name)
		if strings.HasSuffix(name, "Getter") {
			err = visitor.visitorResourceGetterConcreatStruct(name, g.spec.Buf)
		} else if strings.HasSuffix(name, "Interface") {
			err = visitor.visitorInterfaceConcreatStruct(name, g.spec.Buf)
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
					err = visitor.visitorResourceGetterMethod(method.Names[0].Name, method, g.finder.imports, g.spec.Buf)
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
					err = visitor.visitorIntrefaceMethod(Verb(method.Names[0].Name), method, g.finder.imports, g.spec.Buf)
					if err != nil {
						visitor.onError(err)
						return errors.Wrapf(err, "visitorInterfaceMethod name:%s, method:%+v error", method.Names[0].Name, *method)
					}
				}
			}

		}
	}

	err = visitor.visitorEnd(g.spec.Buf)
	if err != nil {
		visitor.onError(err)
		return errors.Wrapf(err, "visitorEnd failed")
	}
	return nil
}

func (g *generator) acceptInterfaceMethod() {
	for _, typeSpec := range g.finder.types() {
		if strings.HasSuffix(typeSpec.Name.Name, "Interface") {
			if interType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				for _, method := range interType.Methods.List {
					fmt.Printf("interface: %+v, %+v\n", *interType, *method)
				}
			}
		}
	}
}

func lowerFirstChar(str string) string {
	return string(unicode.ToLower(rune(str[0]))) + str[1:]
}
