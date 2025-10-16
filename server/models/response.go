package models

type Apiresponse struct {
	Message string      `json:"message" validate:"required,min=1,max=255"`
	Code    int32       `json:"code" validate:"required,min=1"`
	Status  bool        `json:"status" validate:"required"`
	Data    interface{} `json:"data"`
}
