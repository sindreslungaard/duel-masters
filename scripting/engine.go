package scripting

import (
	"fmt"
	"os"
	"strings"

	"github.com/clarkmcc/go-typescript"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

func Setup() {
	vm := goja.New()
	new(require.Registry).Enable(vm)
	console.Enable(vm)

	script, err := os.ReadFile("./scripts/example.ts")

	if err != nil {
		panic(err)
	}

	result, err := typescript.Transpile(strings.NewReader(string(script)))
	if err != nil {
		panic(err)
	}

	println("Running:", result)

	prog, err := goja.Compile("", result, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}
	_, err = vm.RunProgram(prog)
}
