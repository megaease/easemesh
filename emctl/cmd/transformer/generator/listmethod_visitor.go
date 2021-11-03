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
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type (
	statementBuilder interface {
		build(resourceName, subResource string) ([]jen.Code, error)
	}

	listMethodVisitor interface {
		visitorStatusCodeJudgement1() error
		visitorStatusCodeJudgement2() error
		visitorUnmarshalObject() error
		visitorAssignResult() error
		visitorReturn() error
	}

	baseListMethodVisitor struct {
		resourceType ResourceType
		resourceName string
		subResource  string
		group        *jen.Group
		//
		unmarshalStatementMappings   map[ResourceType]statementBuilder
		resultAssignStatementMapping map[ResourceType]statementBuilder
	}

	statementBuilderFunc func(resourceName, subResource string) ([]jen.Code, error)
)

func newListMethodVisitor(resourceType ResourceType, resourceName string,
	subResource string, group *jen.Group) listMethodVisitor {
	return &baseListMethodVisitor{
		resourceType: resourceType,
		resourceName: resourceName,
		subResource:  subResource,
		group:        group,
		unmarshalStatementMappings: map[ResourceType]statementBuilder{
			Service: statementBuilderFunc(serviceUnmarshalObject),
			Global:  statementBuilderFunc(globalUnmarshalObject),
		},

		resultAssignStatementMapping: map[ResourceType]statementBuilder{
			Service: statementBuilderFunc(serviceResultAssign),
			Global:  statementBuilderFunc(globalResultAssign),
		},
	}
}

func (s *baseListMethodVisitor) visitorStatusCodeJudgement1() error {
	s.group.Add(
		jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
			jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
				jen.Id("NotFoundError"),
				jen.Lit("list service"),
			))),
	)

	return nil
}

func (s *baseListMethodVisitor) visitorStatusCodeJudgement2() error {
	s.group.Add(
		jen.If(jen.Id("statusCode").Op(">=").Lit(300)).Op("&&").Id("statusCode").Op("<").Lit(200).Block(
			jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Errorf").Call(
				jen.Lit("call GET %s failed, return statuscode %d text %+v"),
				jen.Id("url"),
				jen.Id("statusCode"),
				jen.Id("b"),
			)),
		),
	)
	return nil
}

func (s *baseListMethodVisitor) visitorReturn() error {
	s.group.Add(jen.Return(jen.Id("results"), jen.Nil()))
	return nil
}

func (s *baseListMethodVisitor) visitorUnmarshalObject() error {
	sb := s.unmarshalStatementMappings[s.resourceType]
	if sb == nil {
		return errors.Errorf("unknown resourceType %s", s.resourceType)
	}
	codes, err := sb.build(s.resourceName, s.subResource)
	if err != nil {
		return err
	}
	for _, c := range codes {
		s.group.Add(c)
	}
	return nil
}

func (s *baseListMethodVisitor) visitorAssignResult() error {
	sb := s.resultAssignStatementMapping[s.resourceType]
	codes, err := sb.build(s.resourceName, s.subResource)
	if err != nil {
		return err
	}
	for _, c := range codes {
		s.group.Add(c)
	}
	return nil
}

func serviceUnmarshalObject(resourceName, subResource string) (result []jen.Code, err error) {
	stmt3 := jen.Id("services").Op(":=").Op("[]").Qual(v1alpha1Pkg, "Service").Block()
	result = append(result, stmt3)
	stmt4 := jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
		jen.Id("b"), jen.Op("&").Id("services"),
	)
	result = append(result, stmt4)
	stmt5 := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
			jen.Id("err"),
			jen.Lit("unmarshal data to v1alpha1.")),
		))
	result = append(result, stmt5)
	return
}

func globalUnmarshalObject(resourceName, subResource string) (codes []jen.Code, err error) {
	resourceVarName := strings.ToLower(resourceName[0:1]) + resourceName[1:]
	def := jen.Id(resourceVarName).Op(":=").Op("[]").Qual(v1alpha1Pkg, resourceName).Op("{}")
	codes = append(codes, def)

	unmarshal := jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
		jen.Id("b"), jen.Op("&").Id(resourceVarName),
	)

	codes = append(codes, unmarshal)
	judgeUnmarshal := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
			jen.Id("err"),
			jen.Lit("unmarshal data to v1alpha1.")),
		))
	codes = append(codes, judgeUnmarshal)
	return
}

func serviceResultAssign(resourceName, subResource string) (codes []jen.Code, err error) {

	capResourceName := strings.ToUpper(resourceName[0:1]) + resourceName[1:]
	stmtDef := jen.Id("results").Op(":=").Op("[]").Op("*").Qual(resourcePkg, capResourceName).Block()
	codes = append(codes, stmtDef)

	var stmtLoop jen.Code
	if subResource != "" {
		fields := strings.Split(capResourceName, subResource)
		if len(fields) != 2 {
			err = errors.Errorf("resource %s must contain sub resource %s", capResourceName, subResource)
			return
		}
		stmtLoop = jen.For(jen.Id("_").Op(",").Id("service").Op(":=").Range().Id("services")).Block(
			jen.If(jen.Id("service").Dot(subResource).Op("!=").Nil()).Block(
				jen.Id("results").Op("=").Append(jen.Id("results"), jen.Qual(resourcePkg, "To"+capResourceName).Call(
					jen.Id("service").Dot("Name"),
					jen.Id("service").Dot(subResource).Dot(fields[1]),
				)),
			))
	} else {
		stmtLoop = jen.For(jen.Id("_").Op(",").Id("service").Op(":=").Range().Id("services")).Block(
			jen.If(jen.Id("service").Dot(capResourceName).Op("!=").Nil()).Block(
				jen.Id("results").Op("=").Append(jen.Id("results"), jen.Qual(resourcePkg, "To"+capResourceName).Call(
					jen.Id("service").Dot("Name"),
					jen.Id("service").Dot(capResourceName),
				)),
			),
		)
	}
	codes = append(codes, stmtLoop)
	return
}

func globalResultAssign(resourceName, subResource string) (codes []jen.Code, err error) {
	resourceVarName := strings.ToLower(resourceName[0:1]) + resourceName[1:]
	def := jen.Id("results").Op(":=").Op("[]*").Qual(resourcePkg, resourceName).Block()
	codes = append(codes, def)
	forLoop := jen.For(jen.Id("_").Op(",").Id("item").Op(":=").Range().Id(resourceVarName)).Block(
		jen.Id("copy").Op(":=").Id("item"),
		jen.Id("results").Op("=").Append(jen.Id("results").Op(",").Qual(resourcePkg, "To"+resourceName).Call(jen.Op("&").Id("copy"))),
	)
	codes = append(codes, forLoop)
	return
}

func (s statementBuilderFunc) build(resourceName, subResource string) ([]jen.Code, error) {
	if s == nil {
		return nil, errors.Errorf("statementBuilderFunc is nil")
	}
	return s(resourceName, subResource)
}

func listMethodAcceptor(visitor listMethodVisitor) error {
	visitorMethods := []func() error{
		visitor.visitorStatusCodeJudgement1,
		visitor.visitorStatusCodeJudgement2,
		visitor.visitorUnmarshalObject,
		visitor.visitorAssignResult,
		visitor.visitorReturn,
	}
	for _, v := range visitorMethods {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}
