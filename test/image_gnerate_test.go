package test

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestWkhtml2Image(t *testing.T) {
	// length: 270,329

	// 270329
	input := "/Users/xulingming/Public/gowork/gowxapi/api/hl.html"
	output := "/Users/xulingming/Public/gowork/gowxapi/test/code.png"
	wkpath, err := exec.LookPath("wkhtmltoimage")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("wkhtmltoimage is availiable at %s\n", wkpath)
	cmd := exec.Command("wkhtmltoimage", "--format", "png", "--quality", "100", "--width", "380", "-", "-")

	cmd.Stdin, _ = os.Open(input) //bytes.NewReader([]byte("a big shark"))

	f, err := os.Create(output)
	defer f.Close()

	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if b, err := io.Copy(f, stdout); err == nil {
		fmt.Println("copy success!")
		fmt.Println(b)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
