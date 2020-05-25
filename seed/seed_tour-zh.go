package seed

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackdon/gowxapi/helper"
	"github.com/jackdon/gowxapi/lesson"
	"go.mongodb.org/mongo-driver/bson"
)

func decideTourSequence(title string) int {
	switch title {
	case "包、变量和函数":
		return 2
	case "流程控制语句：for、if、else、switch 和 defer":
		return 3
	case "并发":
		return 6
	case "更多类型：struct、slice 和映射":
		return 4
	case "欢迎！":
		return 1
	case "方法和接口":
		return 5
	default:
		return 0
	}
}

// ParseTourZh is execute to parsing tour-zh into lessons
func SeedTourZh(tourRoot *string) error {
	content := filepath.Join(*tourRoot, "content")
	dir, err := os.Open(content)
	if err != nil {
		return err
	}
	files, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, f := range files {
		if filepath.Ext(f) != ".article" {
			continue
		}
		fp := filepath.Join(content, f)
		if _, ls, err := lesson.PasrseLesson(&fp, true); err == nil {
			log.Println("Lesson:", ls.Title, "parsed.")
			colc := helper.ConnectDB("lessons")
			colcS := helper.ConnectDB("lesson_sections")

			doc := bson.D{
				{"title", ls.Title},
				{"description", ls.Description},
				{"category", "tour"},
				{"seq", decideTourSequence(ls.Title)},
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
			log.Println("Insert sections:", len(iors.InsertedIDs))
		}
	}
	return nil
}
