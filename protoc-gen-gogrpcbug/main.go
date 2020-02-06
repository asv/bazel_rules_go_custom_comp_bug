package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"

	// To be link gRPC plugin.
	_ "github.com/golang/protobuf/protoc-gen-go/grpc"
)

type gogrpcbug struct {
	gen *generator.Generator
}

func (g *gogrpcbug) Init(gen *generator.Generator) {
	g.gen = gen
}

func (g *gogrpcbug) Name() string {
	return "nextgen"
}

func (g *gogrpcbug) Generate(file *generator.FileDescriptor) {
	// Nothing to-do
}

func (g *gogrpcbug) GenerateImports(file *generator.FileDescriptor) {
	// Nothing to-do
}

func init() {
	generator.RegisterPlugin(new(gogrpcbug))
}

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
