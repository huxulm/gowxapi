package api

import (
	"fmt"
	"log"
	"testing"
)

func TestWXBizDataCrypt(t *testing.T) {

	encryptData := "CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZMQmRzooG2xrDcvSnxIMXFufNstNGTyaGS9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+3hVbJSRgv+4lGOETKUQz6OYStslQ142dNCuabNPGBzlooOmB231qMM85d2/fV6ChevvXvQP8Hkue1poOFtnEtpyxVLW1zAo6/1Xx1COxFvrc2d7UL/lmHInNlxuacJXwu0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn/Hz7saL8xz+W//FRAUid1OksQaQx4CMs8LOddcQhULW4ucetDf96JcR3g0gfRK4PC7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns/8wR2SiRS7MNACwTyrGvt9ts8p12PKFdlqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYVoKlaRv85IfVunYzO0IKXsyl7JCUjCpoG20f0a04COwfneQAGGwd5oa+T8yO5hzuyDb/XcxxmK01EpqOyuxINew=="
	iv := "r7BXXKkLb8qrSNn05n0qiA=="

	crypter := WXBizDataCrypt{
		AppID:  "wx4f4bc4dec97d474b",
		Secret: "tiihtNczf5v6AKRyjwEUhQ==",
	}
	result, _ := crypter.DecryptData(encryptData, iv)
	fmt.Println(result)
}

func TestMarkdownIt(t *testing.T) {
	var CodeMd = "### ok\n" +
		"	```go\n" +
		"	func ok() string {\n" +
		"	}\n" +
		"	```\n"

	log.Println(MarkdownIt(CodeMd))
}
