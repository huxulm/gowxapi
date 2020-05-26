package codesandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	mgp "github.com/gobeam/mongo-go-pagination"
	"github.com/jackdon/gowxapi/api"
	"github.com/jackdon/gowxapi/helper"
	"github.com/jackdon/gowxapi/lesson"
	"github.com/jackdon/gowxapi/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListSandBox provides a list of sandbox
func ListSandBox(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	exampleDesc := "社区精选示例代码"
	tourDesc := "Golang官方示例tour,带你完整的浏览"
	list := []struct {
		Name      string `json:"name"`
		Desc      string `json:"desc"`
		DirCount  int64  `json:"dir_count"`
		FileCount int64  `json:"file_count"`
	}{
		{"Go By Example", exampleDesc, 0, 0}, {"Go Tour之旅", tourDesc, 0, 0},
	}
	colcExample := helper.ConnectDB("examples")
	if exmapleCount, err := colcExample.CountDocuments(context.TODO(), bson.M{}); err == nil {
		list[0].DirCount = exmapleCount
		list[0].FileCount = exmapleCount
	}
	colcLesson := helper.ConnectDB("lessons")
	colcSec := helper.ConnectDB("lesson_sections")
	if tourCount, err := colcLesson.CountDocuments(context.TODO(), bson.M{}); err == nil {
		list[1].DirCount = tourCount
		list[1].FileCount = countSections(colcSec)
	}
	var resp models.Resp
	code := 0
	msg := "ok"
	resp.Code = &code
	resp.Msg = &msg
	resp.Data = &list
	if b, err := json.Marshal(&resp); err == nil {
		fmt.Fprintln(w, string(b))
	} else {
		fmt.Fprintln(w, nil)
	}
}

func countSections(c *mongo.Collection) int64 {
	if count, err := c.CountDocuments(context.TODO(), bson.M{}); err == nil {
		return count
	}
	return 0
}

// LessonSectionPaging is ...
func LessonSectionPaging(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	colSec := helper.ConnectDB("lesson_sections")
	id := ps.ByName("lesson")
	oid, _ := primitive.ObjectIDFromHex(id)
	//match query
	filter := bson.D{{"lesson", oid}}

	pageStr, pageSizeStr := r.FormValue("page"), r.FormValue("pageSize")
	if len(pageStr) == 0 {
		pageStr = "1"
	}
	if len(pageSizeStr) == 0 {
		pageSizeStr = "5"
	}
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
	if pageSize == 0 {
		pageSize = 5
	}

	//match query
	// match := bson.M{"$match": bson.M{ /* "id": bson.M{"$regex": "fu"} */ }}
	selector := bson.M{"_id": 1, "title": 1, "lesson": 1}
	// projection := bson.M{"$project": bson.M{"id": 1, "name": 1, "no": 1, "_id": 1}}
	aggPaginatedData, err := mgp.New(colSec).Limit(pageSize).Filter(filter).Select(selector).Page(page).Find()
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	var secs []*lesson.Page
	for _, raw := range aggPaginatedData.Data {
		var p *lesson.Page
		if marshallErr := bson.Unmarshal(raw, &p); marshallErr == nil {
			secs = append(secs, p)
		}
	}
	msg := "ok"
	code := 0
	result := &models.Resp{Code: &code, Msg: &msg}
	data := make(map[string]interface{}, 2)
	data = map[string]interface{}{"docs": secs, "pagination": &aggPaginatedData.Pagination}
	result.Data = &data
	resultB, err := json.Marshal(result)
	fmt.Fprintf(w, fmt.Sprint(string(resultB)))
}

// LessonPaging is ...
func LessonPaging(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	colLesson := helper.ConnectDB("lessons")
	category := r.FormValue("category")
	// id := ps.ByName("id")
	// oid, _ := primitive.ObjectIDFromHex(id)
	//match query
	filter := bson.D{{"category", category}}

	pageStr, pageSizeStr := r.FormValue("page"), r.FormValue("pageSize")
	if len(pageStr) == 0 {
		pageStr = "1"
	}
	if len(pageSizeStr) == 0 {
		pageSizeStr = "5"
	}
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
	if pageSize == 0 {
		pageSize = 5
	}

	//match query
	// match := bson.M{"$match": bson.M{ /* "id": bson.M{"$regex": "fu"} */ }}
	selector := bson.M{"_id": 1, "title": 1, "description": 1, "category": 1, "seq": 1, "create_time": 1}
	// projection := bson.M{"$project": bson.M{"id": 1, "name": 1, "no": 1, "_id": 1}}
	aggPaginatedData, err := mgp.New(colLesson).Sort("seq", 1).Limit(pageSize).Filter(filter).Select(selector).Page(page).Find()
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	var lessons []*lesson.Lesson
	for _, raw := range aggPaginatedData.Data {
		var lesson *lesson.Lesson
		if marshallErr := bson.Unmarshal(raw, &lesson); marshallErr == nil {
			lessons = append(lessons, lesson)
		}
	}
	msg := "ok"
	code := 0
	result := &models.Resp{Code: &code, Msg: &msg}
	data := make(map[string]interface{}, 2)
	data = map[string]interface{}{"docs": lessons, "pagination": &aggPaginatedData.Pagination}
	result.Data = &data
	resultB, err := json.Marshal(result)
	fmt.Fprintf(w, fmt.Sprint(string(resultB)))
}

// GetLessonSectionDetail is ...
func GetLessonSectionDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	colSec := helper.ConnectDB("lesson_sections")
	id := ps.ByName("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	//match query
	filter := bson.D{{"_id", oid}}

	var s *lesson.Page
	if err := colSec.FindOne(context.TODO(), filter, options.FindOne()).Decode(&s); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			code := 1
			msg := "fail"
			resp := &models.Resp{Code: &code, Msg: &msg, Data: s}
			resultB, _ := json.Marshal(resp)
			fmt.Fprintln(w, string(resultB))
		}
		log.Println(err)
	}

	// highlight code
	highlightFilesContent(s)
	msg := "ok"
	code := 0
	result := &models.Resp{Code: &code, Msg: &msg}
	result.Data = s
	resultB, _ := json.Marshal(s)
	fmt.Fprintf(w, fmt.Sprint(string(resultB)))
}

func highlightFilesContent(s *lesson.Page) {
	if s != nil {
		for i, f := range s.Files {
			w := new(bytes.Buffer)
			api.HlightCode(w, f.Content)
			s.Files[i].ContentHL = w.String()
		}
	}
}
