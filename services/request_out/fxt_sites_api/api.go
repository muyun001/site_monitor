package fxt_sites_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/panwenbin/ghttpclient"
	"net/http"
	"os"
	"strings"
)

const (
	GET_SITE_LIST = "/api/special/for-site-check/site-list?page=1"
	POST_FEEDBACK = "/api/special/for-site-check/feedback"
)

var baseUrl string

func init() {
	baseUrl = os.Getenv("RECEIVE_URLS_AND_SEND_RESULT_API")
}

// apiUrl: 填充API参数并返回完整API地址
func apiUrl(path string, params map[string]string) string {
	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return baseUrl + path
}

func Sites() ([]FxtSite, error) {
	apiUrl := apiUrl(GET_SITE_LIST, nil)
	fxtSites := SiteUrlsResponse{}
	err := ghttpclient.Get(apiUrl, nil).ReadJsonClose(&fxtSites)
	return fxtSites.Data, err
}

func Feedback(checkUrlFeedback CheckUrlFeedback) error {
	apiUrl := apiUrl(POST_FEEDBACK, nil)
	jsonBytes, _ := json.Marshal(checkUrlFeedback)
	res := CodeMessageResponse{}
	err := ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&res)

	if err != nil {
		return err
	}
	if res.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("接收端报错 %d", res.Code))
	}

	return nil
}
