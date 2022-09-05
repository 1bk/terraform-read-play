package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var (
	configFileSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}

	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "description",
			},
			{
				Name: "type",
			},
			{
				Name: "sensitive",
			},
			{
				Name: "sensitiveTwo",
			},
			{
				Name: "nesting",
			},
		},
	}
	nestingBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "val",
			},
		},
	}
)

type Config struct {
	Variables []*Variable
}

type Variable struct {
	Name         string
	Description  string
	Type         string
	Sensitive    bool
	SensitiveTwo string
	Nestings     *Nesting
}

type Nesting struct {
	Val int
}

func main() {
	config := configFromFile("examples/string.tf")
	fmt.Printf("Running\n")
	for _, v := range config.Variables {
		fmt.Printf("%+v\n", v)
	}
	fmt.Printf("Finished")
}

func configFromFile(filePath string) *Config {
	content, err := os.ReadFile(filePath) // go 1.16
	if err != nil {
		log.Fatal(err)
	}

	file, diags := hclsyntax.ParseConfig(content, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		log.Fatal("ParseConfig", diags)
	}

	bodyCont, diags := file.Body.Content(configFileSchema)
	if diags.HasErrors() {
		log.Fatal("file content", diags)
	}

	res := &Config{}

	for _, block := range bodyCont.Blocks {
		v := &Variable{
			Name: block.Labels[0],
		}

		blockCont, diags := block.Body.Content(variableBlockSchema)
		if diags.HasErrors() {
			log.Fatal("block content ", diags)
		}

		if attr, exists := blockCont.Attributes["description"]; exists {
			diags := gohcl.DecodeExpression(attr.Expr, nil, &v.Description)
			if diags.HasErrors() {
				log.Fatal("description attr ", diags)
			}
		}

		if attr, exists := blockCont.Attributes["sensitive"]; exists {
			diags := gohcl.DecodeExpression(attr.Expr, nil, &v.Sensitive)
			if diags.HasErrors() {
				log.Fatal("sensitive attr ", diags)
			}
		}

		if attr, exists := blockCont.Attributes["sensitiveTwo"]; exists {
			diags := gohcl.DecodeExpression(attr.Expr, nil, &v.SensitiveTwo)
			if diags.HasErrors() {
				log.Fatal("sensitiveTwo attr ", diags)
			}
		}

		fmt.Println("Herere000")

		fmt.Printf("blockCont.Attributes: %v\n", blockCont.Attributes)
		// bs, _ := json.Marshal(blockCont.Attributes)
		// fmt.Println(string(bs))

		fmt.Printf("blockCont.Attributes[\"nesting\"]: %v\n", blockCont.Attributes["nesting"])
		// bs, _ = json.Marshal(blockCont.Attributes["nesting"])
		// fmt.Println(string(bs))

		fmt.Printf("blockCont.Blocks: %v\n", blockCont.Blocks)
		fmt.Printf("blockCont.Blocks: %v\n", blockCont)
		if attr, exists := blockCont.Attributes["nesting"]; exists {
			fmt.Println("Herere00")

			fmt.Printf("attr.Expr: %v\n", attr.Expr)
			fmt.Printf("attr.Expr.Variables(): %v\n", attr.Expr.Variables())
			// diags := gohcl.DecodeExpression(attr.Expr, nil, &v.Nestings)
			// fmt.Println("Herere01")
			// if diags.HasErrors() {
			// 	log.Fatal("nesting attr ", diags)
			// }

			// for _, n_block := range bodyCont.Blocks {
			// 	fmt.Println("Herere0")
			// 	_, n_diags := n_block.Body.Content(nestingBlockSchema)
			// 	if n_diags.HasErrors() {
			// 		log.Fatal("n_block content ", n_diags)
			// 	}
			// 	fmt.Println("Herere1")

			// 	// n := &Nesting{}
			// 	// if n_attr, exists := nBlockCont.Attributes["val"]; exists {
			// 	// 	fmt.Println("Herere2")
			// 	// 	nn_diags := gohcl.DecodeExpression(n_attr.Expr, nil, &n.Val)
			// 	// 	if nn_diags.HasErrors() {
			// 	// 		log.Fatal("description attr ", nn_diags)
			// 	// 	}
			// 	// }

			// }
		}

		if attr, exists := blockCont.Attributes["type"]; exists {
			v.Type = hcl.ExprAsKeyword(attr.Expr)
			if v.Type == "" {
				log.Fatal("type attr ", "invalid value")
			}
		}

		res.Variables = append(res.Variables, v)
	}
	return res
}
