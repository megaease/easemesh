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
	"bytes"
	"go/ast"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type (
	interfaceMethodBuilder interface {
		buildGetMethod(*buildInfo) error
		buildPatchMethod(*buildInfo) error
		buildCreateMethod(*buildInfo) error
		buildDeleteMethod(*buildInfo) error
		buildListMethod(*buildInfo) error
	}

	interfaceBuilder struct{}

	genCodeFactory  func(string) jen.Code
	resourceFetcher func(arguments, results []jen.Code) (string, error)
	buildInfo       struct {
		interfaceStructName string
		method              *ast.Field
		imports             []*ast.ImportSpec
		buf                 *jen.File
		resourceType        ResourceType
		subResource         string
		resource2UrlMapping map[string]string
	}
)

var _ interfaceMethodBuilder = &interfaceBuilder{}

const (
	clientPkg   = "github.com/megaease/easemeshctl/cmd/common/client"
	errorsPkg   = "github.com/pkg/errors"
	v1alpha1Pkg = "github.com/megaease/easemesh-api/v1alpha1"
	resourcePkg = "github.com/megaease/easemeshctl/cmd/client/resource"
)

func (g genCodeFactory) generate(resourceName string) jen.Code {
	return g(resourceName)
}

func (r resourceFetcher) do(arguments, results []jen.Code) (string, error) {
	if r == nil {
		return "", errors.Errorf("fetcher can't be nil")
	}
	return r(arguments, results)
}
func (i *interfaceBuilder) buildGetMethod(info *buildInfo) (err error) {
	factories := []genCodeFactory{
		buildURLStatement(info),
		buildGetByContextHTTPCallStatement(info),
		buildJudgeResponseStatement(info),
		buildReturnStatement(info),
	}

	err = i.buildCommonMethodBody(info, factories, readMethodFetcher)
	if err != nil {
		return errors.Wrapf(err, "build get method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildPatchMethod(info *buildInfo) (err error) {
	factories := []genCodeFactory{
		buildURLStatement(info),
		buildResourceToObjectStatement(info),
		buildPutByContextStatement(info),
		buildReturnErrStatement(info),
	}
	err = i.buildCommonMethodBody(info, factories, writeMethodFetcher)
	if err != nil {
		return errors.Wrapf(err, "build patch method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildCreateMethod(info *buildInfo) (err error) {

	factories := []genCodeFactory{
		buildURLStatement(info),
		buildCreateByContextStatement(info),
		buildReturnErrStatement(info),
	}
	err = i.buildCommonMethodBody(info, factories, writeMethodFetcher)
	if err != nil {
		return errors.Wrapf(err, "build create method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildDeleteMethod(info *buildInfo) (err error) {

	factories := []genCodeFactory{
		buildURLStatement(info),
		buildDeleteByContextStatement(info),
		buildReturnErrStatement(info),
	}

	err = i.buildCommonMethodBody(info, factories, deleteMethodFetcher(info.interfaceStructName))

	if err != nil {
		return errors.Wrapf(err, "build delete method of the interface error")
	}
	return nil

}

func (i *interfaceBuilder) buildListMethod(info *buildInfo) error {
	factories := []genCodeFactory{
		buildListURLStatement(info),
		buildListByContextStatement(info),
		buildListJudgeErrReturnStatement(info),
		buildListReturnStatement(info),
	}
	err := i.buildCommonMethodBody(info, factories, readMethodFetcher)
	if err != nil {
		return errors.Wrapf(err, "build list method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildCommonMethodBody(info *buildInfo, factories []genCodeFactory, fetcher resourceFetcher) (err error) {
	var arguments, results []jen.Code
	var funcName string
	err = covertFuncType(info.method, info.imports).
		extractArguments(&arguments).
		extractResults(&results).
		extractFuncName(&funcName).
		error()

	if err != nil {
		return errors.Wrapf(err, "extract arguments and result fo the interface error")
	}

	resourceName, err := fetcher.do(arguments, results)
	if err != nil {
		return errors.Wrapf(err, "fetch resource name failed")
	}

	info.buf.Func().Params(
		jen.Id(string(info.interfaceStructName[0])).Op("*").Id(info.interfaceStructName),
	).Id(funcName).Params(arguments...).Params(results...).BlockFunc(func(g *jen.Group) {
		err = i.buildMethodBody(resourceName, factories, g)
	})
	if err != nil {
		return errors.Wrapf(err, "build body of the %s interface method", resourceName)
	}
	return nil
}

func (i *interfaceBuilder) buildMethodBody(resourceName string, factories []genCodeFactory, g *jen.Group) error {
	for _, factory := range factories {
		g.Add(factory.generate(resourceName))
	}
	return nil
}
func buildURLStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		resourceFirstName := strings.ToLower(resourceName[0:1])

		return jen.Id("url").Op(":=").Qual("fmt", "Sprintf").Call(
			jen.Lit("http://").Op("+").
				Id(resourceFirstName).Dot("client").Dot("server").
				Op("+").Id("apiURL").Op("+").Lit("/mesh/services/%s/").Op("+").Lit(resourceName),
			jen.Id("args_1"),
		)
	}
}
func buildResourceToObjectStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		return jen.Id("object").Op(":=").Id("args_1").Dot("ToV1Alpha1").Call()
	}
}

func buildGetByContextHTTPCallStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
		return jen.Id("r").Op(",").Id("err").Op(":=").
			Qual(clientPkg, "NewHTTPJSON").Call().
			Dot("GetByContext").Call(
			jen.Id("args_0"),
			jen.Id("url"),
			jen.Nil(),
			jen.Nil(),
		).Dot("HandleResponse").Call(
			jen.Func().Params(
				jen.Id("b").Op("[]").Byte(),
				jen.Id("statusCode").Int(),
			).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
				stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("NotFoundError"),
						jen.Lit("get "+resourceName+" %s"),
						jen.Id("args_1"),
					)))

				stmt2 := jen.If(jen.Id("statusCode").Op(">=").Lit(300)).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Errorf").Call(
						jen.Lit("call %s failed, return status code %d text %+v"),
						jen.Id("url"), jen.Id("statusCode"), jen.Id("b"),
					)),
				)
				stmt3 := jen.Id(resourceName).Op(":=").Op("&").Qual(v1alpha1Pkg, capResourceName).Block()
				stmt4 := jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
					jen.Id("b"), jen.Id(resourceName),
				)
				stmt5 := jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("err"),
						jen.Lit("unmarshal data to v1alpha1."+resourceName)),
					))

				returnStmt := jen.Return(jen.Qual(resourcePkg, "To"+capResourceName).Call(
					jen.Id("args_1"), jen.Id(resourceName),
				).Op(",").Nil())
				g1.Add(stmt1)
				g1.Add(stmt2)
				g1.Add(stmt3)
				g1.Add(stmt4)
				g1.Add(stmt5)
				g1.Add(returnStmt)
			}),
		)
	}
}

func buildJudgeResponseStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		return jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		)
	}
}

func buildReturnStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
		return jen.Return(jen.Id("r").Op(".").Parens(jen.Op("*").Qual(resourcePkg, capResourceName)).Op(",").Nil())
	}
}

func buildPutByContextStatement(info *buildInfo) func(string) jen.Code {

	return func(resourceName string) jen.Code {
		return jen.Id("_").Op(",").Id("err").Op(":=").
			Qual(clientPkg, "NewHTTPJSON").Call().
			Dot("PutByContext").Call(
			jen.Id("args_0"),
			jen.Id("url"),
			jen.Id("object"),
			jen.Nil(),
		).Dot("HandleResponse").Call(
			jen.Func().Params(
				jen.Id("b").Op("[]").Byte(),
				jen.Id("statusCode").Int(),
			).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
				stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("NotFoundError"),
						jen.Lit("patch "+resourceName+" %s"),
						jen.Id("args_1").Dot("Name").Call(),
					)))

				stmt2 := jen.If(jen.Id("statusCode").Op("<").Lit(300)).Op("&&").Id("statusCode").Op(">=").Lit(200).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				)
				returnStmt := jen.Return(jen.Nil(),
					jen.Qual(errorsPkg, "Errorf").Call(
						jen.Lit("call PUT %s failed, return statuscode %d text %+v"),
						jen.Id("url"), jen.Id("statusCode"), jen.Id("b"),
					),
				)
				g1.Add(stmt1)
				g1.Add(stmt2)
				g1.Add(returnStmt)
			}),
		)
	}
}

func buildReturnErrStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		return jen.Return(jen.Id("err"))
	}
}

func buildDeleteByContextStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		return jen.Id("_").Op(",").Id("err").Op(":=").
			Qual(clientPkg, "NewHTTPJSON").Call().
			Dot("DeleteByContext").Call(
			jen.Id("args_0"),
			jen.Id("url"),
			jen.Nil(),
			jen.Nil(),
		).Dot("HandleResponse").Call(
			jen.Func().Params(
				jen.Id("b").Op("[]").Byte(),
				jen.Id("statusCode").Int(),
			).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
				stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("NotFoundError"),
						jen.Lit("Delete "+resourceName+" %s"),
						jen.Id("args_1"),
					)))

				stmt2 := jen.If(jen.Id("statusCode").Op("<").Lit(300)).Op("&&").Id("statusCode").Op(">=").Lit(200).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				)
				returnStmt := jen.Return(jen.Nil(),
					jen.Qual(errorsPkg, "Errorf").Call(
						jen.Lit("call Delete %s failed, return statuscode %d text %+v"),
						jen.Id("url"), jen.Id("statusCode"), jen.Id("b"),
					),
				)
				g1.Add(stmt1)
				g1.Add(stmt2)
				g1.Add(returnStmt)
			}))
	}
}

func buildCreateByContextStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		return jen.Id("_").Op(",").Id("err").Op(":=").
			Qual(clientPkg, "NewHTTPJSON").Call().
			Dot("PostByContext").Call(
			jen.Id("args_0"),
			jen.Id("url"),
			jen.Nil(),
			jen.Nil(),
		).Dot("HandleResponse").Call(
			jen.Func().Params(
				jen.Id("b").Op("[]").Byte(),
				jen.Id("statusCode").Int(),
			).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
				stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusConflict")).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("ConflictError"),
						jen.Lit("create "+resourceName+" %s"),
						jen.Id("args_1").Dot("Name").Call(),
					)))

				stmt2 := jen.If(jen.Id("statusCode").Op("<").Lit(300)).Op("&&").Id("statusCode").Op(">=").Lit(200).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				)
				returnStmt := jen.Return(jen.Nil(),
					jen.Qual(errorsPkg, "Errorf").Call(
						jen.Lit("call Post %s failed, return statuscode %d text %+v"),
						jen.Id("url"), jen.Id("statusCode"), jen.Id("b"),
					),
				)
				g1.Add(stmt1)
				g1.Add(stmt2)
				g1.Add(returnStmt)
			}))
	}
}

func buildListURLStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		resourceFirstName := strings.ToLower(resourceName[0:1])
		return jen.Id("url").Op(":=").Lit("http://").Op("+").
			Id(resourceFirstName).Dot("client").Dot("server").
			Op("+").Id("apiURL").Op("+").Lit("/mesh/services")
	}
}
func buildListByContextStatement(info *buildInfo) func(string) jen.Code {
	return func(resourceName string) jen.Code {
		capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
		return jen.Id("result").Op(",").Id("err").Op(":=").
			Qual(clientPkg, "NewHTTPJSON").Call().
			Dot("GetByContext").Call(
			jen.Id("args_0"),
			jen.Id("url"),
			jen.Nil(),
			jen.Nil(),
		).Dot("HandleResponse").Call(
			jen.Func().Params(
				jen.Id("b").Op("[]").Byte(),
				jen.Id("statusCode").Int(),
			).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
				stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("NotFoundError"),
						jen.Lit("list service"),
					)))
				stmt2 := jen.If(jen.Id("statusCode").Op(">=").Lit(300)).Op("&&").Id("statusCode").Op("<").Lit(200).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Errorf").Call(
						jen.Lit("call GET %s failed, return statuscode %d text %+v"),
						jen.Id("url"),
						jen.Id("statusCode"),
						jen.Id("b"),
					)),
				)
				stmt3 := jen.Id("services").Op(":=").Op("[]").Qual(v1alpha1Pkg, "Service").Block()
				stmt4 := jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
					jen.Id("b"), jen.Op("&").Id("services"),
				)
				stmt5 := jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
						jen.Id("err"),
						jen.Lit("unmarshal data to v1alpha1.")),
					))
				stmt6 := jen.Id("results").Op(":=").Op("[]").Op("*").Qual(resourcePkg, capResourceName).Block()
				stmtLoop := jen.For(jen.Id("_").Op(",").Id("service").Op(":=").Range().Id("services")).Block(
					jen.If(jen.Id("service").Dot(capResourceName).Op("!=").Nil()).Block(
						jen.Id("results").Op("=").Append(jen.Id("results"), jen.Qual(resourcePkg, "To"+capResourceName).Call(
							jen.Id("service").Dot("Name"),
							jen.Id("service").Dot(capResourceName),
						)),
					),
				)
				stmtReturn := jen.Return(jen.Id("results"), jen.Nil())
				g1.Add(stmt1)
				g1.Add(stmt2)
				g1.Add(stmt3)
				g1.Add(stmt4)
				g1.Add(stmt5)
				g1.Add(stmt6)
				g1.Add(stmtLoop)
				g1.Add(stmtReturn)

			}),
		)
	}
}

func buildListJudgeErrReturnStatement(info *buildInfo) func(string) jen.Code {

	return func(resourceName string) jen.Code {
		return jen.If(jen.Id("err").Op("!=").Nil().Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		))
	}
}
func buildListReturnStatement(info *buildInfo) func(string) jen.Code {

	return func(resourceName string) jen.Code {
		capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
		return jen.Return(jen.Id("result").Op(".").Parens(jen.Op("[]").Op("*").Qual(resourcePkg, capResourceName)), jen.Nil())
	}
}

func readMethodFetcher(arguments, results []jen.Code) (string, error) {
	if len(results) < 2 {
		return "", errors.Errorf("read method should return two arguments")
	}
	return extractResourceName(results[0])
}

func writeMethodFetcher(arguments, results []jen.Code) (string, error) {
	if len(arguments) < 2 {
		return "", errors.Errorf("read method should return two arguments")
	}
	return extractResourceName(arguments[1])
}

func extractResourceName(statement jen.Code) (string, error) {
	buf := &bytes.Buffer{}

	s, ok := statement.(*jen.Statement)
	if !ok {
		return "", errors.Errorf("code should be a statements, but %+v", statement)
	}

	s1 := *s

	// Add {} to suppress render error
	s1.Block()
	err := s1.Render(buf)
	if err != nil {
		return "", errors.Wrapf(err, "can not render statement:%+v", *s)
	}

	sections := strings.Split(string(buf.Bytes()), "resource.")
	if len(sections) < 1 {
		return "", errors.Errorf("rendered statement should contains 'resource.' but %s", string(buf.Bytes()))
	}

	// Trimming added {}
	return strings.Trim(sections[len(sections)-1], "{}"), nil
}

// we extract resource name of delete method from interfaceStructName
func deleteMethodFetcher(interfaceStructName string) resourceFetcher {
	var resourceName string
	result := strings.Split(interfaceStructName, "Interface")
	var err error
	if len(result) < 2 {
		err = errors.Errorf("the interface name %s don't contain resource", interfaceStructName)
	} else {
		resourceName = result[0]
	}

	resourceName = strings.ToUpper(resourceName[0:1]) + resourceName[1:]

	return func(arguments, results []jen.Code) (string, error) {
		return resourceName, err
	}
}
