package api

type NewsResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type NotifyResponse struct {
	Status int `json:"status"`
}
