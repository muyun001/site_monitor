package fxt_sites_api

import "site-monitor/structs"

type CodeMessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SiteUrlsResponse struct {
	Data []FxtSite `json:"data"`
}

type FxtSite struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

type CheckUrlFeedback struct {
	Data []EachUrlFeedback `json:"data"`
}

type EachUrlFeedback struct {
	ID            int         `json:"id"`
	IsOpenOK      bool        `json:"is_open_ok"`
	IsIndexed     bool        `json:"is_indexed"`
	IsHacked      bool        `json:"is_hacked"`
	IsRobotsOK    bool        `json:"is_robots_ok"`
	ResStatusCode int         `json:"res_status_code"`
	TDK           structs.TDK `json:"tdk"`
	HeadlessTDK   structs.TDK `json:"headless_tdk"`
	Url           string      `json:"url"`
	HeadlessUrl   string      `json:"headless_url"`
	ResContent    string      `json:"res_content"`
}