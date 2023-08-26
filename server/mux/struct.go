package mux

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/codfrm/cago/pkg/utils/httputils"
)

// Meta 路由
type Meta struct {
}

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) Do(ctx context.Context, req any, resp any) error {
	route, ok := reflect.TypeOf(req).Elem().FieldByName("Meta")
	if !ok {
		return errors.New("invalid method, second parameter must have Meta field")
	}
	// 必须有Route字段
	if !ok || route.Type != reflect.TypeOf(Meta{}) {
		return errors.New("invalid method, second parameter must have Meta field")
	}
	method := route.Tag.Get("method")
	path := route.Tag.Get("path")

	ptrVal := reflect.ValueOf(req)
	ptrElem := ptrVal.Elem()
	ptrType := ptrElem.Type()

	query := url.Values{}
	data := make(map[string]interface{})
	var form func(key string, value any)
	switch method {
	case http.MethodGet, http.MethodDelete:
		form = func(key string, value any) {
			if s, ok := value.(string); ok {
				query.Add(key, s)
			}
		}
	default:
		form = func(key string, value any) {
			data[key] = value
		}
	}

	for i := 0; i < ptrElem.NumField(); i++ {
		tag := ptrType.Field(i).Tag
		if tag == "" {
			continue
		}
		if key := tag.Get("form"); key != "" {
			if key == "-" {
				continue
			}
			if key == ",inline" {
				data[key] = ptrElem.Field(i).Interface()
			} else {
				// 处理key,label的情况,例如: key,default=1
				key, opts := head(key, ",")
				opts, val := head(opts, "=")
				if opts == "default" && ptrElem.Field(i).IsZero() {
					form(key, val)
				} else {
					form(key, ptrElem.Field(i).Interface())
				}
			}
		}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(b)

	path = c.baseURL + path
	if len(query) != 0 {
		path += "?" + query.Encode()
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()
	b, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	jsonResp := &httputils.JSONResponse{
		Data: resp,
	}
	if err := json.Unmarshal(b, jsonResp); err != nil {
		return err
	}
	if jsonResp.Code != 0 {
		return httputils.NewError(httpResp.StatusCode, jsonResp.Code, jsonResp.Msg)
	}
	return nil
}
