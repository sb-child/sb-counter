package consts

type CounterUser struct {
	Path string `json:"path"`
	View string `json:"view"`
	DB   string `json:"db"`
}

type CounterData struct {
	All int
	Today int
	Yesterday int
	BeforeYesterday int
}
