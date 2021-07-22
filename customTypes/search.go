package customTypes

type SearchResponse struct {
	Code    string   `json:"code"`
	Keyword string   `json:"keyword"`
	Datas   []string `json:"datas"`
}
