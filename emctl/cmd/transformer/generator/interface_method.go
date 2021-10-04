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

	genCodeFactory func() jen.Code

	buildInfo struct {
		resourceName        string
		interfaceStructName string
		method              *ast.Field
		imports             []*ast.ImportSpec
		buf                 *jen.File
	}
)

var _ interfaceMethodBuilder = &interfaceBuilder{}

const (
	clientPkg   = "github.com/megaease/easemeshctl/cmd/common/client"
	errorsPkg   = "github.com/pkg/errors"
	v1alpha1Pkg = "github.com/megaease/easemesh-api/v1alpha1"
	resourcePkg = "github.com/megaease/easemeshctl/cmd/client/resource"
)

func (g genCodeFactory) generate() jen.Code {
	return g()
}

func (i *interfaceBuilder) buildGetMethod(info *buildInfo) (err error) {
	factories := []genCodeFactory{
		buildURLStatement(info.resourceName, info.interfaceStructName),
		buildGetByContextHTTPCallStatement(info.resourceName, info.interfaceStructName),
		buildJudgeResponseStatement(info.resourceName, info.interfaceStructName),
		buildReturnStatement(info.resourceName, info.interfaceStructName),
	}

	err = i.buildCommonMethodBody(info, factories)
	if err != nil {
		return errors.Wrapf(err, "build get method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildPatchMethod(info *buildInfo) (err error) {
	factories := []genCodeFactory{
		buildURLStatement(info.resourceName, info.interfaceStructName),
		buildResourceToObjectStatement(info.resourceName, info.interfaceStructName),
		buildPutByContextStatement(info.resourceName, info.interfaceStructName),
		buildReturnErrStatement(info.resourceName, info.interfaceStructName),
	}
	err = i.buildCommonMethodBody(info, factories)
	if err != nil {
		return errors.Wrapf(err, "build patch method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildCreateMethod(info *buildInfo) (err error) {

	factories := []genCodeFactory{
		buildURLStatement(info.resourceName, info.interfaceStructName),
		buildCreateByContextStatement(info.resourceName, info.interfaceStructName),
		buildReturnErrStatement(info.resourceName, info.interfaceStructName),
	}
	err = i.buildCommonMethodBody(info, factories)
	if err != nil {
		return errors.Wrapf(err, "build create method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildDeleteMethod(info *buildInfo) (err error) {

	factories := []genCodeFactory{
		buildURLStatement(info.resourceName, info.interfaceStructName),
		buildDeleteByContextStatement(info.resourceName, info.interfaceStructName),
		buildReturnErrStatement(info.resourceName, info.interfaceStructName),
	}
	err = i.buildCommonMethodBody(info, factories)
	if err != nil {
		return errors.Wrapf(err, "build delete method of the interface error")
	}
	return nil

}

func (i *interfaceBuilder) buildListMethod(info *buildInfo) error {
	factories := []genCodeFactory{
		buildListURLStatement(info.resourceName, info.interfaceStructName),
		buildListByContextStatement(info.resourceName, info.interfaceStructName),
		buildListJudgeErrReturnStatement(info.resourceName, info.interfaceStructName),
		buildListReturnStatement(info.resourceName, info.interfaceStructName),
	}
	err := i.buildCommonMethodBody(info, factories)
	if err != nil {
		return errors.Wrapf(err, "build list method of the interface error")
	}
	return nil
}

func (i *interfaceBuilder) buildCommonMethodBody(info *buildInfo, factories []genCodeFactory) (err error) {
	var arguments, results []jen.Code
	var funcName string
	err = covertFuncType(info.method, info.imports).
		extractArguments(&arguments).
		extractResults(&results).
		extractFuncName(&funcName).
		error()

	if err != nil {
		return errors.Wrapf(err, "extract arguments and result fo the %s interface error", info.resourceName)
	}

	info.buf.Func().Params(
		jen.Id(string(info.resourceName[0])).Op("*").Id(info.interfaceStructName),
	).Id(funcName).Params(arguments...).Params(results...).BlockFunc(func(g *jen.Group) {
		err = i.buildMethodBody(factories, g)
	})
	if err != nil {
		return errors.Wrapf(err, "build body of the %s interface method", info.resourceName)
	}
	return nil
}

func (i *interfaceBuilder) buildMethodBody(factories []genCodeFactory, g *jen.Group) error {
	for _, factory := range factories {
		g.Add(factory.generate())
	}
	return nil
}
func buildURLStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
		return jen.Id("url").Op(":=").Qual("fmt", "Sprintf").Call(
			jen.Lit("http://").Op("+").
				Id("c").Dot("client").Dot("server").
				Op("+").Id("apiURL").Op("+").Lit("/mesh/services/%s/").Op("+").Lit(resourceName),
			jen.Id("args_1"),
		)
	}
}
func buildResourceToObjectStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
		return jen.Id("object").Op(":=").Id("args_1").Dot("ToV1Alpha1").Call()
	}
}

func buildGetByContextHTTPCallStatement(resourceName, interfaceStructName string) func() jen.Code {
	capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
	return func() jen.Code {
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

func buildJudgeResponseStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
		return jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		)
	}
}

func buildReturnStatement(resourceName, interfaceStructName string) func() jen.Code {
	capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
	return func() jen.Code {
		return jen.Return(jen.Id("r").Op(".").Parens(jen.Op("*").Qual(resourcePkg, capResourceName)).Op(",").Nil())
	}
}

func buildPutByContextStatement(resourceName, interfaceStructName string) func() jen.Code {

	return func() jen.Code {
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

func buildReturnErrStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
		return jen.Return(jen.Id("err"))
	}
}

func buildDeleteByContextStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
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

func buildCreateByContextStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
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

func buildListURLStatement(resourceName, interfaceStructName string) func() jen.Code {
	return func() jen.Code {
		return jen.Id("url").Op(":=").Lit("http://").Op("+").
			Id("c").Dot("client").Dot("server").
			Op("+").Id("apiURL").Op("+").Lit("/mesh/services")
	}
}
func buildListByContextStatement(resourceName, interfaceStructName string) func() jen.Code {
	capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
	return func() jen.Code {
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

func buildListJudgeErrReturnStatement(resourceName, interfaceStructName string) func() jen.Code {

	return func() jen.Code {
		return jen.If(jen.Id("err").Op("!=").Nil().Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		))
	}
}
func buildListReturnStatement(resourceName, interfaceStructName string) func() jen.Code {

	capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
	return func() jen.Code {
		return jen.Return(jen.Id("result").Op(".").Parens(jen.Op("[]").Op("*").Qual(resourcePkg, capResourceName)), jen.Nil())
	}
}
