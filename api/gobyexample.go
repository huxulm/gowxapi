package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	mgp "github.com/gobeam/mongo-go-pagination"
	"github.com/jackdon/gowxapi/helper"
	"github.com/jackdon/gowxapi/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PageGoExample provides a paging list with a standard pagination interface.
func PageGoExample(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := helper.ConnectDB("examples")
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

	// projection := bson.M{"$project": bson.M{"id": 1, "name": 1, "no": 1, "_id": 1}}
	aggPaginatedData, err := mgp.New(collection).Sort("no", 1).Limit(pageSize).Filter(bson.M{}).Select(bson.M{"id": 1, "name": 1, "no": 1, "_id": 1}).Page(page).Find()
	if err != nil {
		panic(err)
	}

	var aggExampleList []models.ExampleBase
	for _, raw := range aggPaginatedData.Data {
		var example *models.ExampleBase
		if marshallErr := bson.Unmarshal(raw, &example); marshallErr == nil {
			aggExampleList = append(aggExampleList, *example)
		}
	}
	msg := "ok"
	code := 0
	result := &models.Resp{Code: &code, Msg: &msg}
	data := make(map[string]interface{}, 2)
	data = map[string]interface{}{"docs": &aggExampleList, "pagination": &aggPaginatedData.Pagination}
	result.Data = &data

	resultB, err := json.Marshal(result)
	fmt.Fprintf(w, fmt.Sprint(string(resultB)))
}

func highlightCode(result *models.ExampleBase, comment string) string {
	var buf bytes.Buffer
	source := make([]string, 0)
	if len(comment) != 0 && comment == "on" {
		source = append(source, *result.GoCode)
	} else {
		for _, e := range *&result.Segs[0] {
			source = append(source, e.Code)
		}
	}
	HlightCode(&buf, strings.Join(source, "\n"))
	return buf.String()
}
func markupDocs(result *models.ExampleBase) string {
	source := make([]string, 0)
	for _, e := range *&result.Segs[0] {
		source = append(source, e.Docs)
	}
	return MarkdownIt(strings.Join(source, "\n"))
}
func highlightRunDocs(result *models.ExampleBase) string {
	source := make([]string, 0)
	for _, e := range *&result.Segs[1] {
		source = append(source, e.Code)
	}
	var buf bytes.Buffer
	HlightCode(&buf, strings.Join(source, "\n"))
	return buf.String()
}

// GoExampleDetail is ...
func GoExampleDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collection := helper.ConnectDB("examples")
	comment := r.FormValue("comment")
	id := ps.ByName("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	//match query
	filter := bson.D{{"_id", oid}}
	var result models.ExampleBase
	var resp *models.Resp
	err := collection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			code := 1
			msg := "fail"
			resp = &models.Resp{Code: &code, Msg: &msg, Data: &result}
			fmt.Fprintln(w, resp)
		}
		log.Println(err)
	}
	code := 0
	msg := "ok"
	hlc := highlightCode(&result, comment)
	hlcClean := highlightCode(&result, "")
	result.HighlightCode = &hlc
	result.HighlightCodeClean = &hlcClean
	docs := markupDocs(&result)
	result.DocsMarkup = &docs
	runDocs := highlightRunDocs(&result)
	result.RunDocs = &runDocs
	resp = &models.Resp{Code: &code, Msg: &msg, Data: &result}
	respB, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(respB))
}
