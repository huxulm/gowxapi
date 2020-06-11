package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var css = DefaultCSS

func TestInjectStyle(t *testing.T) {
	fp := filepath.Join("../", "api/hl.html")
	if srcHtlm, err := ioutil.ReadFile(fp); err == nil {
		if out := InjectStyle(srcHtlm, css); out != nil {
			os.Remove(fp)
			ioutil.WriteFile(fp, out, 0777)
		}
	}
}
