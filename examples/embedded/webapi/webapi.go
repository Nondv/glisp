package main

import (
	// _ "embed"
	"fmt"
	"io"
	"net/http"
	"os"

	"nondv.io/glisp/interpreter"
	. "nondv.io/glisp/types"
)

// //go:embed router.lisp
// var routerCode string

func main() {
	pathToRouter := "examples/embedded/webapi/router.lisp"
	routerCode, err := os.ReadFile(pathToRouter)
	routerCodeFileStats, err := os.Stat(pathToRouter)

	if err != nil {
		panic(err)
	}

	baseBindings := interpreter.BuildBaseBindings()
	interpreter.ReadEval(baseBindings, `(load "lang/core.lisp")`)
	interpreter.ReadEval(baseBindings, `(load "lang/alist.lisp")`)
	interpreter.ReadEvalAll(baseBindings, string(routerCode))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// auto-reloading
		fileStats, err := os.Stat(pathToRouter)
		if err == nil && fileStats.ModTime().After(routerCodeFileStats.ModTime()) {
			routerCode, err = os.ReadFile(pathToRouter)
			if err == nil {
				fmt.Println("Reloading router.lisp")
				routerCodeFileStats = fileStats
				interpreter.ReadEvalAll(baseBindings, string(routerCode))
			}
		}

		handle(baseBindings, w, r)
	})
	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handle(baseBindings *Bindings, w http.ResponseWriter, r *http.Request) {
	bindings := baseBindings.AssocSym("request-data", prepareRequestData(r))
	result, err := interpreter.ReadEval(bindings, "(router)")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Something went wrong")
		return
	}

	bindings = bindings.AssocSym("result", result)
	body, _ := interpreter.ReadEval(bindings, `(alist/get "body" result)`)
	if !body.IsString() {
		w.WriteHeader(500)
		fmt.Fprint(w, "Something went wrong")
		return
	}

	status, _ := interpreter.ReadEval(bindings, `(alist/get "status" result)`)
	if status.IsInteger() {
		w.WriteHeader(status.ToInt())
	}

	fmt.Fprint(w, body.ToStr())
}

func prepareRequestData(r *http.Request) *Value {
	result := BuildEmptyList()
	result = alistAssoc(result, "path", BuildString(r.URL.Path))
	result = alistAssoc(result, "method", BuildString(r.Method))

	body, _ := io.ReadAll(r.Body)
	if string(body) != "" {
		result = alistAssoc(result, "body", BuildString(string(body)))
	}

	query := BuildEmptyList()
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			val := BuildString(values[len(values)-1])
			query = alistAssoc(query, key, val)
		}
	}
	result = alistAssoc(result, "query", query)

	return result
}

func alistAssoc(alist *Value, key string, value *Value) *Value {
	return BuildCons(BuildCons(BuildString(key), value), alist)
}
