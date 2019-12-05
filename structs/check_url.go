package structs

type ResAndTDK struct {
	ResStatusCode int    `json:"res_status_code"`
	ResContent    string `json:"res_content"`
	TDK           TDK    `json:"tdk"`
}

type TDK struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}
