package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/go-gosh/gask/app/conf"
)

var client = resty.New()

type response[T any] struct {
	Code    int
	Message string
	Data    T
}

func DecodeBody[T any](res *resty.Response) (data T, err error) {
	if !res.IsSuccess() {
		err = errors.New(res.String())
		return
	}
	var b response[T]
	err = json.Unmarshal(res.Body(), &b)
	if err != nil {
		return
	}
	if b.Code > 399 {
		err = errors.New(b.Message)
		return
	}
	data = b.Data
	return
}

func ShouldDecodeBody[T any](res *resty.Response, err error) (T, error) {
	var t T
	if err != nil {
		return t, err
	}
	return DecodeBody[T](res)
}

func Create[T any](api string, body any) (T, error) {
	return ShouldDecodeBody[T](client.R().SetBody(body).Post(api))
}

func Update(api string, id uint, body any) (any, error) {
	return ShouldDecodeBody[any](client.R().SetPathParam("id", fmt.Sprintf("%v", id)).SetBody(body).Put(api))
}

func Delete(api string, id uint) (any, error) {
	return ShouldDecodeBody[any](client.R().SetPathParam("id", fmt.Sprintf("%v", id)).Delete(api))
}

func Retrieve[T any](api string, id uint) (T, error) {
	return ShouldDecodeBody[T](client.R().SetPathParam("id", fmt.Sprintf("%v", id)).Get(api))
}

func Paginate[T any](api string, params map[string]string) (T, error) {
	return ShouldDecodeBody[T](client.R().SetQueryParams(params).Get(api))
}

func ApiRoute(api string) string {
	return fmt.Sprintf("http://localhost:%d/api/v1/%s", conf.GetConfig().Port, api)
}
