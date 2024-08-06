package models

type Anime struct {
	ID          uint64 `json:"mal_id"`
	Title       string `json:"title"`
	Description string `json:"synopsis"`
	Episodes    uint16 `json:"episodes"`
}

type Data struct {
	Data Anime `json:"data"`
}
