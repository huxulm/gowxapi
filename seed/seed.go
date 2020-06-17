package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jackdon/gowxapi/config"
	"github.com/jackdon/gowxapi/helper"
)

var Config = config.C

// Seg is a segment of an example
type Seg struct {
	Docs, DocsRendered              string
	Code, CodeRendered, CodeForJs   string
	CodeEmpty, CodeLeading, CodeRun bool
}

// ExampleBase is example info extracted from gernerate.json
type ExampleBase struct {
	ID, Name                    *string
	NO                          *int64
	GoCode, GoCodeHash, URLHash *string
	Segs                        [][]*Seg
}

// Seed is ...
func Seed(seedFilePath string) (*[]ExampleBase, error) {
	genB, err := ioutil.ReadFile(seedFilePath)
	if err != nil {
		return nil, err
	}
	examples := make([]ExampleBase, 0)
	errParse := json.Unmarshal(genB, &examples)
	if errParse != nil {
		panic(errParse)
	}
	fmt.Println("Parsed length:", len(examples))
	return &examples, nil
}

func seedExamples() {
	filePath := Config.Seed.SeedFile
	if len(filePath) == 0 || Config.Seed.Seed == false {
		// filePath = path.Join("../", "generate.jsonn")
		log.Println("Seed skipped.")
		return
	}
	examples, _ := Seed(filePath)
	if examples != nil && len(*examples) > 0 {
		// b, _ := json.Marshal(examples)
		colc := helper.ConnectDB("examples")
		docs := make([]interface{}, len(*examples))
		for i, exa := range *examples {
			docs[i] = exa
		}
		imr, err := colc.InsertMany(context.TODO(), docs)
		if err != nil {
			log.Fatalln("Seed failed.")
			panic(err)
		}
		log.Println("successfuly seed", imr.InsertedIDs)
	} else {
		log.Println("seed 0")
	}
}

func seedTour() {
	seed := Config.SeedTour.Seed
	if seed == false {
		log.Println("Seed tour skipped.")
		return
	}
	tourRoot := Config.SeedTour.TourPath
	if len(tourRoot) == 0 {
		log.Println("TourRoot is empty. Seed tour skipped.")
		return
	}
	if err := SeedTourZh(&tourRoot); err != nil {
		log.Println("Seed tour failed:", err)
	}
}

func seedLeetcode() {
	seed := Config.SeedLeetcode.Seed
	if seed == false {
		log.Println("Seed leetcode skipped.")
		return
	}
	root := Config.SeedLeetcode.LeetcodePath
	if len(root) == 0 {
		log.Println("Leetcode root is empty. Seed Leetcode skipped.")
		return
	}
	if _, err := SeedLeetcode(&root); err != nil {
		log.Println("Seed leetcode failed:", err)
	}
}

func init() {
	seedExamples()
	seedTour()
	seedLeetcode()
}
