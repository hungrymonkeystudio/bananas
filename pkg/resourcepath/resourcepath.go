// properly find the resource path

package resourcepath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetResourcePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {	
		panic("No caller information")
	}
	path := strings.Split(filepath.Dir(filename), "/")
	textPath := strings.Join(path[0:len(path)-2], "/")
	if textPath == os.Getenv("HOME")+"/.local/bin" {
		fmt.Println("On local install")
		return os.Getenv("HOME")+"/.local/share/bananas"
	} else if textPath == "/usr/local/bin" {
		return "/usr/local/share/bananas"
	}
	return textPath + "/resources"
}
