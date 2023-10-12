package model

type EmailRequest struct {
	Id          int64 `json:"id"`
	TemplateId  int   `json:"templateId"`
	Data        []KeyValue
	UserEmail   string `json:"userEmail"`
	Recipient   string `json:"recipient"`
	Sent        string `json:"sent"`
	Receivetime string `json:"receivetime"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EmailRequests []*EmailRequest
