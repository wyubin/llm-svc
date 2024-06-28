package ctl

type RespModels struct {
	Models []RespModel `json:"models"`
}

type RespModel struct {
	Name  string `json:"name"`
	Model string `json:"model"`
	Size  int    `json:"size"`
}

type RespEmbedding struct {
	Embedding []float32 `json:"embedding"`
}

type RespChat struct {
	Model   string  `json:"model"`
	Message RespMsg `json:"message"`
}

type RespMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
