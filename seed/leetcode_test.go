package seed

import (
	"testing"
)

func TestTransformPipe(t *testing.T)  {
	root := "/home/bx/Public/git_repos/awesome-golang-leetcode"
	if err := TransformPipe(root + "/src", "/home/bx/Public/workspace/gowxapi/.dist"); err != nil {
		t.Fatal(err)
	}
}