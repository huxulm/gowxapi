package codesandbox

import (
	"io/ioutil"
	"testing"

	"github.com/jackdon/gowxapi/api"
	"github.com/jackdon/gowxapi/helper"
)

func TestGenHTMLImage(t *testing.T) {
	if b, err := helper.GenHTMLImage([]byte(api.Code), nil); err == nil {
		ioutil.WriteFile("share.png", b, 0777)
	} else {
		t.Fatal(err)
	}
}
