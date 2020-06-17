package seed

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackdon/gowxapi/helper"
	"github.com/jackdon/gowxapi/lesson"
	"go.mongodb.org/mongo-driver/bson"
)

func decideSequence(dir string) int {
	splits := strings.Split(dir, ".")
	if len(splits) > 0 {
		r, err := strconv.Atoi(splits[1])
		if err == nil {
			return r
		}
	}
	return 0
}

// SeedLeetcode is execute to parsing leetcode into lessons
func SeedLeetcode(root *string) (*[]lesson.Lesson, error) {
	errCount := 0
	lessons := make([]lesson.Lesson, 0)
	rootDir, err := os.Open(*root)
	if err != nil {
		return nil, err
	}
	dirs, err := rootDir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		df, err := os.Open(filepath.Join(*root, dir))
		if err != nil {
			return nil, err
		}
		files, err := df.Readdirnames(0)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if filepath.Ext(f) != ".md" {
				continue
			}
			fp := filepath.Join(df.Name(), f)
			if _, ls, err := lesson.PasrseLesson(&fp, true); err == nil {
				log.Println("Lesson:", ls.Title, "parsed.")
				lessons = append(lessons, *ls)
				colc := helper.ConnectDB("lessons")
				colcS := helper.ConnectDB("lesson_sections")

				doc := bson.D{
					{"title", ls.Title},
					{"description", ls.Description},
					{"category", "leetcode"},
					{"seq", decideSequence(dir)},
					{"create_time", time.Now()},
				}
				ior, err := colc.InsertOne(context.TODO(), doc)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println("Insert lesson:", *ior)
				sections := make([]interface{}, len(ls.Pages))
				for i, p := range ls.Pages {
					p.Lesson = ior.InsertedID
					sections[i] = p
				}
				iors, err := colcS.InsertMany(context.TODO(), sections)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Insert sections: %v\n", iors)
			} else {
				log.Println(err)
				errCount++
			}
		}
	}
	fmt.Printf(`error counts: %v\n`, errCount)
	return &lessons, nil
}
