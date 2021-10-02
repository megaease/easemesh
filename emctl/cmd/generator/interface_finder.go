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
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

type interfaceFinder struct {
	filename  string
	found     []*ast.TypeSpec
	typeNames map[string]string
}

func newInterfaceFinder(filename string) *interfaceFinder {
	return &interfaceFinder{filename: filename}
}

func (i *interfaceFinder) parseFile() error {

	find := func(n ast.Node) bool {

		var ts *ast.TypeSpec
		var ok bool

		if ts, ok = n.(*ast.TypeSpec); !ok {
			return true
		}

		if ts.Name == nil {
			return true
		}
		if _, ok := ts.Type.(*ast.InterfaceType); !ok {
			return true
		}

		i.typeNames[ts.Name.Name] = ts.Name.Name

		i.found = append(i.found, ts)

		return false
	}

	node, err := parser.ParseFile(token.NewFileSet(), i.filename, nil, 0)
	if err != nil {
		return errors.Wrapf(err, "parse file %s", i.filename)
	}

	ast.Inspect(node, find)
	return nil
}

func (i *interfaceFinder) typeNameResults() (types []string) {
	for k := range i.typeNames {
		types = append(types, k)
	}
	return
}

func (i *interfaceFinder) types() []*ast.TypeSpec {
	return i.found
}
