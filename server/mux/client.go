package mux

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type DoOptions struct {
	path string
}

type DoOption func(*DoOptions)

func NewDoOptions(opts ...DoOption) *DoOptions {
	options := &DoOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithPath(path string) DoOption {
	return func(options *DoOptions) {
		options.path = path
	}
}

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) Request(ctx context.Context, req, resp any, opts ...DoOption) (*http.Request, error) {
	options := NewDoOptions(opts...)
	route, ok := reflect.TypeOf(req).Elem().FieldByName("Meta")
	if !ok {
		return nil, errors.New("invalid method, second parameter must have Meta field")
	}
	// 必须有Route字段
	if !ok || route.Type != reflect.TypeOf(Meta{}) {
		return nil, errors.New("invalid method, second parameter must have Meta field")
	}
	method := route.Tag.Get("method")
	path := route.Tag.Get("path")
	if options.path != "" {
		path = options.path
	}

	ptrVal := reflect.ValueOf(req)
	ptrElem := ptrVal.Elem()
	ptrType := ptrElem.Type()

	query := url.Values{}
	data := make(map[string]interface{})
	var form func(key string, value any)
	switch method {
	case http.MethodGet, http.MethodDelete:
		form = func(key string, value any) {
			switch value.(type) {
			case string:
				query.Add(key, value.(string))
			case bool:
				query.Add(key, strconv.FormatBool(value.(bool)))
			case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
				query.Add(key, fmt.Sprintf("%d", value))
			case float32, float64:
				query.Add(key, fmt.Sprintf("%f", value))
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
		} else if uri := tag.Get("uri"); uri != "" {
			path = strings.ReplaceAll(path, ":"+uri, ptrElem.Field(i).String())
		}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(b)

	path = c.baseURL + path
	if len(query) != 0 {
		path += "?" + query.Encode()
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func (c *Client) Do(ctx context.Context, req any, resp any, opts ...DoOption) error {
	httpReq, err := c.Request(ctx, req, resp, opts...)
	if err != nil {
		return err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()
	b, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	jsonResp := &httputils.JSONResponse{
		Data: resp,
	}
	if err := json.Unmarshal(b, jsonResp); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}
	if jsonResp.Code != 0 {
		return httputils.NewError(httpResp.StatusCode, jsonResp.Code, jsonResp.Msg)
	}
	return nil
}
