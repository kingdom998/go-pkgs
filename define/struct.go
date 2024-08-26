package define

type HttpResp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type GenerateParam struct {
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
	Route    string `json:"route,omitempty"`
	MsgId    string `json:"msgId,omitempty"`
	Body     []byte `json:"body"`
}

type GenerateResp struct {
	Images []string `json:"images"`
	MsgId  string   `json:"msg_id"`
}

type QueryParam struct {
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
	MsgId    string `json:"msg_id,omitempty" form:"msg_id"`
}
