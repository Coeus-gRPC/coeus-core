// Simple trick to make all `go test` call to happen under project root (to circumvent from path problem)
// https://brandur.org/fragments/testing-go-project-root

package testing

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
