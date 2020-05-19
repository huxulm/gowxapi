package seed

import (
	"context"
	"log"
	"testing"

	"github.com/jackdon/gowxapi/helper"
)

func TestSeed(t *testing.T) {
	examples, _ := Seed("../generate.json")
	if len(*examples) > 0 {
		// b, _ := json.Marshal(examples)
		colc := helper.ConnectDB("examples")
		docs := make([]interface{}, len(*examples))
		for i, exa := range *examples {
			no := int64(i + 1)
			exa.NO = &no
			docs[i] = exa
		}
		imr, err := colc.InsertMany(context.TODO(), docs)
		if err != nil {
			log.Fatalln("Seed failed.")
			panic(err)
		}
		log.Println("successfuly seed", imr.InsertedIDs)
	} else {
		return
	}
}

func TestQueryPage(t *testing.T) {
	println("test query")
}
