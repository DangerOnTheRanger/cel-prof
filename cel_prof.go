package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

const (
	cpuProfPath = "cpu.prof"
	memProfPath = "mem.prof"
)

func main() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("l", decls.NewListType(decls.String)),
			decls.NewVar("s", decls.String),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s cel-file arg-map", os.Args[0])
	}
	cpuProf, err := os.Create(cpuProfPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := pprof.StartCPUProfile(cpuProf); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()
	uncompiledProg, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ast, issues := env.Compile(string(uncompiledProg))
	if issues != nil && issues.Err() != nil {
		log.Fatal(issues.Err())
	}
	compiledProg, err := env.Program(ast)
	if err != nil {
		log.Fatal(err)
	}
	rawArgs, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	var argMap map[string]interface{}
	json.Unmarshal([]byte(rawArgs), &argMap)
	out, _, err := compiledProg.Eval(argMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
	memProf, err := os.Create(memProfPath)
	if err != nil {
		log.Fatal(err)
	}
	runtime.GC()
	if err := pprof.WriteHeapProfile(memProf); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
