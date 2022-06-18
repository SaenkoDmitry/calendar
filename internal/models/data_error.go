package models

type DataError struct {
	Data interface{}    `json:"data"`
	Err  *InternalError `json:"err"`
}

type InternalError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
