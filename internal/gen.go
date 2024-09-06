package internal

import (
	"fmt"

	"github.com/pubgo/protoc-gen-go-mongo/pkg/lavamongo"
	"github.com/pubgo/protoc-gen-go-mongo/pkg/mongopbv1"
	"google.golang.org/protobuf/proto"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".pb.mongo.go"
	genFile := jen.NewFile(string(file.GoPackageName))
	genFile.HeaderComment("Code generated by protoc-gen-go-mongo. DO NOT EDIT.")
	genFile.HeaderComment("versions:")
	genFile.HeaderComment(fmt.Sprintf("- protoc-gen-go-mongo %s", version))
	genFile.HeaderComment(fmt.Sprintf("- protoc              %s", protocVersion(gen)))
	if file.Proto.GetOptions().GetDeprecated() {
		genFile.HeaderComment(fmt.Sprintf("%s is a deprecated file.", file.Desc.Path()))
	} else {
		genFile.HeaderComment(fmt.Sprintf("source: %s", file.Desc.Path()))
	}

	genFile.Comment("This is a compile-time assertion to ensure that this generated file")
	genFile.Comment("is compatible with the grpc package it is being compiled against.")
	genFile.Comment("Requires gRPC-Go v1.32.0 or later.")
	genFile.Id("const _ =").Qual("google.golang.org/grpc", "SupportPackageIsVersion7")

	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.Skip()

	if len(file.Messages) == 0 {
		return g
	}

	//var events = make(map[string]any)
	for _, msg := range file.Messages {
		opts, ok := proto.GetExtension(msg.Desc.Options(), mongopbv1.E_Options).(*lavamongo.Options)
		if !ok || !opts.GetEnable() {
			continue
		}

		g.Unskip()
		for _, mm := range msg.Fields {
			genFile.Commentf("%sField%s", msg.GoIdent.GoName, mm.GoIdent.GoName)
		}
	}

	g.P(genFile.GoString())
	return g
}

func genName(name1, name2 string) map[string]string {
	return map[string]string{
		"Field":                    name1 + "." + name2,
		"FieldForUpdate":           name1 + ".$." + name2,
		"FieldForUpdateMany":       name1 + ".$[]." + name2,
		"FieldForUpdateSpecified":  name1 + ".$[elem]." + name2,
		"FieldForUpdateSpecified1": name1 + ".$[elem]." + name2,
	}
}
