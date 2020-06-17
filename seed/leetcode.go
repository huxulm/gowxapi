package seed

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TransformPipe pip leetcode to lesson
func TransformPipe(src, dest string) error {
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	// clean dest dir
	if err := os.RemoveAll(dest); err == nil {
		err = os.Mkdir(dest, 0777)
	}
	if err != nil {
		return err
	}

	subDirs, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	fileReg := regexp.MustCompile("^\\d{4}.")

	for _, subDirName := range subDirs {
		m := fileReg.Match([]byte(subDirName))
		if !m || strings.HasPrefix(subDirName, "0000.") {
			continue
		}
		title := getTitleByDir(subDirName)
		if len(title) <= 0 {
			continue
		}

		// copy dir
		destDir := filepath.Join(dest, subDirName)
		err := os.Mkdir(destDir, 0777)
		if err != nil {
			return nil
		}
		// copy dir files
		subDir, err := os.Open(filepath.Join(src, subDirName))
		defer subDir.Close()
		if err != nil {
			return err
		}
		if files, err := subDir.Readdirnames(0); err == nil {
			codes := make([]string, 0)
			desc := make([]string, 0)
			for _, f := range files {
				// copy file
				if strings.HasSuffix(f, ".go") {
					codes = append(codes, f)
				}
				if strings.HasSuffix(f, ".md") {
					desc = append(desc, f)
				}
			}
			for _, c := range codes {
				copyFile(filepath.Join(subDir.Name(), c), filepath.Join(destDir, c))
			}
			copyPresent(filepath.Join(subDir.Name(), desc[0]), filepath.Join(destDir, desc[0]), codes, title)
		}
	}
	return nil
}

func getTitleByDir(path string) string {
	if len(path) < 5 {
		return ""
	}
	return string([]rune(path)[5:])
}

func copyFile(src, dest string) error {
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	_, err = io.Copy(df, sf)
	return err
}

func copyPresent(src, dest string, codes []string, title string) error {
	play := ""
	for _, c := range codes {
		play = fmt.Sprintf("%s.play %s\n", play, c)
	}
	df, err := os.Create(dest)
	defer df.Close()
	if err != nil {
		return err
	}
	bwf := bufio.NewWriter(df)

	sf, err := os.Open(src)
	defer sf.Close()
	if err != nil {
		return err
	}
	bsf := bufio.NewReader(sf)
	for {
		line, _, err := bsf.ReadLine()
		bwf.Write(replaceLine(line, play, title))
		if err == io.EOF {
			break
		} //do something
	}
	bwf.Flush()
	return err
}

var (
	// replace **example** => example
	reg1         = regexp.MustCompile(`\*{2}(\S+|(\w+\s+\S+))\*{2}`)
	titleLineReg = regexp.MustCompile(`^# [\S][\S]`)
	titleReg     = regexp.MustCompile(`\[([\s\d.\w])+\]`)
	descReg      = regexp.MustCompile(`\#{2}\sDescription`)
	sumReg       = regexp.MustCompile(`\#{2}\s结语`)
	headerReg    = regexp.MustCompile(`(#+|>+)\s`)
	// omit line starts with:
	// > [!
	// * 1
	// > [!
	// * 2
	// > [!
	// [title]:
	// [me]:
	omitReg = regexp.MustCompile(`>\s\[!|\[title\]:|\[me\]:|(\* \d)`)
)

func replaceLine(line []byte, play string, title string) []byte {
	rs := line
	switch {
	case omitReg.Match(line):
		rs = []byte("")
	case reg1.Match(line):
		if bytes.Index(line, []byte("*")) == 0 {
			rs = bytes.ReplaceAll(line, []byte("**"), []byte(""))
		} else {
			rs = bytes.ReplaceAll(line, []byte("**"), []byte("`"))
		}
	case titleLineReg.Match(line):
		r := titleReg.FindAll(line, -1)
		if len(r) > 0 {
			tb := bytes.Trim(r[0], "[]")
			rs = tb[bytes.Index(tb, []byte(" "))+1:]
		}
	case descReg.Match(line):
		rs = append([]byte("* "), []byte(title)...)
	case sumReg.Match(line):
		rs = append([]byte(play), line[3:]...)
	case headerReg.Match(line):
		rs = line[bytes.Index(line, []byte(" "))+1:]
	}
	return append(rs, []byte("\n")...)
}
