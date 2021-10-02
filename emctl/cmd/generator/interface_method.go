package main

import (
	"go/ast"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type (
	interfaceMethodBuilder interface {
		buildGetMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error
		buildPatchMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error
		buildCreateMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error
		buildDeleteMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error
		buildListMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error
	}

	interfaceBuilder struct{}
)

const (
	clientPkg   = "github.com/megaease/easemeshctl/cmd/common/client"
	errorsPkg   = "github.com/pkg/errors"
	v1alpha1Pkg = "github.com/megaease/easemesh-api/v1alpha1"
	resourcePkg = "github.com/megaease/cmd/client/resource"
)

func (i *interfaceBuilder) buildGetMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) (err error) {
	var arguments, results []jen.Code
	var funcName string
	err = covertFuncType(method).
		extractArguments(&arguments).
		extractResults(&results).
		extractFuncName(&funcName).
		error()

	if err != nil {
		return errors.Wrapf(err, "extract arguments and result of the %s interface get methods error", resourceName)
	}
	buf.Func().Params(
		jen.Id(string(resourceName[0])).Op("*").Id(interfaceStructName),
	).Id(funcName).Params(arguments...).Params(results...).BlockFunc(func(g *jen.Group) {
		err = i.buildGetMethodBody(resourceName, interfaceStructName, g)
	})
	if err != nil {
		return errors.Wrapf(err, "build body of the %s interface Get method", resourceName)
	}
	return nil
}

func (i *interfaceBuilder) buildGetMethodBody(resourceName, interfaceStructName string, g *jen.Group) error {
	capResourceName := strings.ToUpper(string(resourceName[0])) + resourceName[1:]
	urlStmt := jen.Id("url").Op(":=").Qual("fmt", "Sprintf").Call(
		jen.Lit("http://").Op("+").
			Id("c").Dot("client").Dot("server").
			Op("+").Id("apiURL").Op("+").Lit("/mesh/services/%s/").Op("+").Lit(resourceName),
		jen.Id("args_1"),
	)
	httpCall := jen.Id("r").Op(",").Id("err").Op(":=").Qual(clientPkg, "NewHTTPJSON").
		Call().
		Dot("GetByContext").Call(
		jen.Id("args_0"),
		jen.Id("url"),
		jen.Nil(),
		jen.Nil()).Dot("HandleResponse").Call(
		jen.Func().Params(
			jen.Id("b").Op("[]").Byte(),
			jen.Id("statusCode").Int(),
		).Params(jen.Interface(), jen.Error()).BlockFunc(func(g1 *jen.Group) {
			stmt1 := jen.If(jen.Id("statusCode").Op("==").Qual("http", "StatusNotFound")).Block(
				jen.Return(jen.Nil(), jen.Qual(errorsPkg, "Wrapf").Call(
					jen.Id("NotFoundError"),
					jen.Lit("unmarshal data to v1alpha1."+resourceName),
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

	judgeStmt := jen.If(jen.Id("err").Op("!=").Nil()).Block(
		jen.Return(jen.Nil(), jen.Id("err")),
	)

	returnStmt := jen.Return(jen.Id("r").Op(".").Parens(jen.Op("*").Qual(resourcePkg, capResourceName)).Op(",").Nil())

	g.Add(urlStmt)
	g.Add(httpCall)
	g.Add(judgeStmt)
	g.Add(returnStmt)
	return nil
}

func (i *interfaceBuilder) buildPatchMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error {
	return nil
}

func (i *interfaceBuilder) buildCreateMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error {
	return nil
}

func (i *interfaceBuilder) buildDeleteMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error {
	return nil
}

func (i *interfaceBuilder) buildListMethod(resourceName, interfaceStructName string, method *ast.Field, buf *jen.File) error {
	return nil
}
