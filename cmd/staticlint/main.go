// Package main демонстрирует создание пользовательского анализатора для запрещения вызова os.Exit в функции main.
// Включает в себя механизм чтения конфигурации из файла и запуска анализаторов через multichecker.
package main

import (
	"encoding/json"
	"go/ast"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// Config — константа, содержащая имя файла конфигурации.
const Config = `lint_config.json`

// ConfigData описывает ожидаемую структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string // Список названий анализаторов staticcheck для активации.
}

// ExitCheckAnalyzer представляет пользовательский анализатор, который ищет прямые вызовы os.Exit в функции main.
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "prohibits direct calls to os.Exit within the main function to ensure graceful shutdown",
	Run:  run,
}

// run выполняет проверку каждого файла исходного кода в пакете на наличие вызовов os.Exit.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, что мы находимся в пакете main и в функции main
		if pass.Pkg.Name() == "main" && file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.CallExpr:
					if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
						if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "os" && fun.Sel.Name == "Exit" {
							pass.Reportf(x.Pos(), "direct call to os.Exit in main function is prohibited")
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}

// main читает конфигурацию, подготавливает и запускает анализаторы.
func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(cwd, Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		ExitCheckAnalyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
