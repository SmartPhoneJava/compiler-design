package main

import (
	"kurs/internal"
	"kurs/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func Parse(filename string) {

	is, err := antlr.NewFileStream(filename)
	if err != nil {
		panic(err)
	}
	lexer := parser.NewLuaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewLuaParser(stream)

	collector := internal.NewInfoCollector()

	antlr.ParseTreeWalkerDefault.Walk(collector, p.Chunk())

	collector.Funcs.MustVisualize("assets", "func_calls.dot")
	internal.FuncTableMustVisualize(collector.Funcs, "assets", "vars.dot")
}

func main() {
	Parse("tables.lua")
}
