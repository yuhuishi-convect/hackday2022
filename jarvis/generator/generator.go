package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func autoGenerator() {
	log.Println("event: Code Generation")
	err := entc.Generate("../ent/schema/", &gen.Config{
		Header: "// Auto generated",
		IDType: &field.TypeInfo{Type: field.TypeInt},
	})
	if err != nil {
		log.Fatal("running ent codegen:", err)
	}
}

func main() {
	autoGenerator()
}
