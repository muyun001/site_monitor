package services

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/panwenbin/ghttpclient"
	"github.com/samclarke/robotstxt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"regexp"
	"site-monitor/services/request_out/fxt_sites_api"
	"site-monitor/structs"
	"site-monitor/util"
	"strconv"
	"strings"
	"time"
)

func CheckUrl() error {
	needCheckSites, err := fxt_sites_api.Sites()
	if err != nil {
		return err
	}

	checkUrlResponse := fxt_sites_api.CheckUrlFeedback{}

	needCheckSiteChan := make(chan fxt_sites_api.FxtSite)
	eachUrlFeedbackChan := make(chan fxt_sites_api.EachUrlFeedback)
	for i := 0; i < 5; i++ {
		go func() {
			for {
				idAndUrl, ok := <-needCheckSiteChan
				if !ok {
					break
				}
				eachUrlFeedback, _ := CheckEachUrl(idAndUrl)
				eachUrlFeedbackChan <- eachUrlFeedback
			}
		}()
	}
	go func() {
		for _, idAndUrl := range needCheckSites {
			needCheckSiteChan <- idAndUrl
		}
		close(needCheckSiteChan)
	}()

	for i := 0; i < len(needCheckSites); i++ {
		checkUrlResponse.Data = append(checkUrlResponse.Data, <-eachUrlFeedbackChan)
	}

	err = fxt_sites_api.Feedback(checkUrlResponse)
	if err != nil {
		return err
	}

	return nil
}

// CheckEachUrl: 处理每个url
func CheckEachUrl(idAndUrl fxt_sites_api.FxtSite) (fxt_sites_api.EachUrlFeedback, error) {
	openFailResponse := fxt_sites_api.EachUrlFeedback{
		ID:         idAndUrl.ID,
		Url:        idAndUrl.Url,
		IsOpenOK:   false,
		IsIndexed:  true,
		IsHacked:   false,
		IsRobotsOK: true,
		TDK:        structs.TDK{},
	}

	urlResAndTDK, err := urlResAndTDK(idAndUrl.Url)
	if err != nil {
		return openFailResponse, err
	}

	isStatusCodeOk := urlResAndTDK.ResStatusCode < 400
	if isStatusCodeOk == false {
		return openFailResponse, errors.New(fmt.Sprintf("status code %d", urlResAndTDK.ResStatusCode))
	}

	eachUrlResponse := fxt_sites_api.EachUrlFeedback{
		ID:       idAndUrl.ID,
		IsOpenOK: true,
		TDK:      urlResAndTDK.TDK,
		Url:      idAndUrl.Url,
	}

	eachUrlResponse.IsIndexed, _ = isIncludedInBaidu(idAndUrl.Url)
	eachUrlResponse.IsRobotsOK, _ = isRobotOK(idAndUrl.Url)
	eachUrlResponse.IsHacked, _ = isHacked(idAndUrl.Url)

	return eachUrlResponse, nil
}

// urlResAndTDK 获取访问页面的返回结果,状态码和TDK
func urlResAndTDK(url string) (*structs.ResAndTDK, error) {
	resContent, resStatusCode, err := visitUrl(url)
	if err != nil {
		resContent, resStatusCode, err = visitUrl(url)
		if err != nil {
			return &structs.ResAndTDK{}, err
		}
	}

	tdk, err := TDK(resContent)
	if err != nil {
		return &structs.ResAndTDK{}, err
	}

	response := &structs.ResAndTDK{
		ResStatusCode: resStatusCode,
		ResContent:    resContent,
		TDK:           tdk,
	}

	return response, nil
}

// TDK 获取TDK
func TDK(html string) (structs.TDK, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return structs.TDK{}, err
	}

	title := dom.Find("title").Text()

	keywords := ""
	keywordItem := dom.Find("meta[name='keywords']")
	if len(keywordItem.Nodes) == 0 {
		keywordItem = dom.Find("meta[name='Keywords']")
	}

	if len(keywordItem.Nodes) != 0 {
		for _, attr := range keywordItem.Nodes[0].Attr {
			if attr.Key == "content" {
				keywords = attr.Val
			}
		}
	}

	description := ""
	descriptionItem := dom.Find("meta[name='description']")
	if len(descriptionItem.Nodes) == 0 {
		descriptionItem = dom.Find("meta[name='Description']")
	}

	if len(descriptionItem.Nodes) != 0 {
		for _, attr := range descriptionItem.Nodes[0].Attr {
			if attr.Key == "content" {
				description = attr.Val
			}
		}
	}

	tdk := structs.TDK{
		Title:       title,
		Keywords:    keywords,
		Description: description,
	}

	return tdk, nil
}

// isIncludedInBaidu 判断是否有收录
func isIncludedInBaidu(url string) (bool, error) {
	isIncluded := false
	domain := strings.Replace(strings.Replace(strings.Replace(url, "http:", "", -1), "https:", "", -1), "/", "", -1)
	baiduUrl := "https://www.baidu.com/s?wd=site%3A" + domain
	responseBody, _, err := visitUrl(baiduUrl)
	if err != nil {
		responseBody, _, err = visitUrl(baiduUrl)
		if err != nil {
			return false, err
		}
	}

	r := regexp.MustCompile(`百度为您找到相关结果约(.*?)个`)
	subMatch := r.FindStringSubmatch(responseBody)
	if len(subMatch) == 2 {
		include, err := strconv.Atoi(strings.Replace(subMatch[1], ",", "", -1))
		if err != nil {
			return false, err
		}

		if include != 0 {
			isIncluded = true
		}
	}

	return isIncluded, nil
}

// isRobotOK 判断网站是否允许程序访问网站内的链接
func isRobotOK(url string) (bool, error) {
	isRobotsOK := true
	robotUrl := url + "robots.txt"
	contents := `
        User-agent: *
        Disallow: /dir/
    `
	robots, err := robotstxt.Parse(contents, robotUrl)
	if err != nil {
		return isRobotsOK, nil
	}

	allowed, _ := robots.IsAllowed("Sams-Bot/1.0", url)
	if !allowed {
		isRobotsOK = false
	}

	return isRobotsOK, nil
}

// isHacked 通过无头浏览器访问检查网址是否被入侵
func isHacked(url string) (bool, error) {
	headlessUrl := ""
	_, headlessUrl, err := headlessTDKAndUrl(url)
	if err != nil {
		return false, err
	}

	domain := strings.Replace(url, "http://", "", -1)
	if strings.Contains(domain, ".") {
		if headlessUrl != "" && !strings.Contains(headlessUrl, strings.Split(domain, ".")[1]) {
			return true, nil
		}
	} else if headlessUrl != "" && !strings.Contains(headlessUrl, domain) {
		return true, nil
	}

	return false, nil
}

// headlessTDKAndUrl 获取无头浏览器访问url的TDK
func headlessTDKAndUrl(url string) (structs.TDK, string, error) {
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args:  []string{"--headless"},
	}

	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, os.Getenv("HEADLESS_URL_PREFIX"))
	if err != nil {
		return structs.TDK{}, "", err
	}

	defer wd.Quit()
	wd.SetPageLoadTimeout(time.Second * 30)
	err = wd.Get(url)
	if err != nil {
		return structs.TDK{}, "", err
	}

	sourceCode, err := wd.PageSource()
	if err != nil {
		return structs.TDK{}, "", err
	}

	headlessUrl, err := wd.CurrentURL()
	if err != nil {
		return structs.TDK{}, "", err
	}

	headlessTDK, err := TDK(sourceCode)
	if err != nil {
		return structs.TDK{}, "", err
	}

	return headlessTDK, headlessUrl, nil
}

// visitUrl 访问url,获取返回内容和状态码
func visitUrl(url string) (string, int, error) {
	client := ghttpclient.NewClient().
		SslSkipVerify(true).
		Timeout(time.Second * 10).
		UserAgent(util.RandomUserAgent()).
		Url(url).
		Get()
	response, err := client.Response()
	if err != nil {
		return "", 0, err
	}
	statusCode := response.StatusCode

	body, err := client.TryUTF8ReadBodyClose()
	if err != nil {
		return "", 0, err
	}

	return string(body), statusCode, nil
}
