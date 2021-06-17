package main

import (
	"io/ioutil"
	"kurs/internal"
	"kurs/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func Parse(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	is, err := antlr.NewFileStream(filename)
	if err != nil {
		panic(err)
	}
	lexer := parser.NewLuaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewLuaParser(stream)

	collector := internal.NewInfoCollector(string(data))

	antlr.ParseTreeWalkerDefault.Walk(collector, p.Chunk())

	collector.Funcs.MustVisualize("assets/out", "func_calls.dot")
	internal.FuncTableMustVisualize(collector.Funcs, "assets/out", "vars.dot")
}

func main() {
	Parse("assets/lua/rich_funcs.lua")
}
