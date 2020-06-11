package codesandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jackdon/gowxapi/config"
	"github.com/jackdon/gowxapi/helper"
	"github.com/jackdon/gowxapi/lesson"
	"github.com/jackdon/gowxapi/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SharePost create a post image by lesson_page's id
func SharePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userInfo := map[string]string{"nick": r.FormValue("nick"), "avatar": r.FormValue("avatar")}

	colSec := helper.ConnectDB("lesson_sections")
	id := ps.ByName("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	var s *lesson.Page
	if err := colSec.FindOne(context.TODO(), filter, options.FindOne()).Decode(&s); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			handleError(err, w)
		}
		return
	}

	htmlBuf := new(bytes.Buffer)
	helper.HlightCode2Html(htmlBuf, s.Files[0].Content)
	html, err := helper.GenHTMLImage(helper.InjectStyle(htmlBuf.Bytes(), helper.DefaultCSS), &map[string]string{
		"format": "png", "quality": "100", "width": "380",
	}) // getCode()

	/* avatarURL := r.FormValue("avatar")  */
	avatar, err := downloadAvatar(userInfo["avatar"]) // getAvatar() //
	if err != nil {
		handleError(err, w)
		return
	}
	logo, err := getLogo()
	if err != nil {
		handleError(err, w)
		return
	}
	if post, err := helper.Generate(html, avatar, logo, userInfo); err != nil {
		handleError(err, w)
		return
	} else {
		w.Header().Add("Content-Type", "image/png")
		w.Write(post)
	}
}

func handleError(err error, w io.Writer) {
	code := 1
	msg := "fail"
	resp := &models.Resp{Code: &code, Msg: &msg, Data: nil}
	resultB, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(resultB))
}

func getCode() ([]byte, error) {
	return helper.LoadLogo("/Users/xulingming/Public/gowork/gowxapi/test/code.png")
}
func getAvatar() ([]byte, error) {
	return helper.LoadLogo(config.C.StaticResource.AvatarPath)
}
func getLogo() ([]byte, error) {
	return helper.LoadLogo(config.C.StaticResource.LogoPath)
}

func downloadAvatar(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
