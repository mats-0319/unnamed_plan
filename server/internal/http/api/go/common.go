package api

type Pagination struct {
	Size int `json:"size"`
	Num  int `json:"num"`
}

type ResBase struct {
	IsSuccess bool   `json:"is_success"`
	Err       string `json:"err"`
}

var ResBaseSuccess = ResBase{IsSuccess: true}
