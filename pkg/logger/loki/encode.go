package loki

import (
	"encoding/json"

	"github.com/codfrm/cago/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type lokiEncode struct {
	*zapcore.MapObjectEncoder
	pool buffer.Pool
}

func NewLokiEncode() zapcore.Encoder {
	return &lokiEncode{
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
		pool:             buffer.NewPool(),
	}
}

func (e *lokiEncode) Clone() zapcore.Encoder {
	return e.clone()
}

func (e *lokiEncode) clone() *lokiEncode {
	encode := zapcore.NewMapObjectEncoder()
	for k, v := range e.MapObjectEncoder.Fields {
		encode.Fields[k] = v
	}
	return &lokiEncode{
		MapObjectEncoder: encode,
		pool:             e.pool,
	}
}

type lokiPush struct {
	Streams [1]lokiPushStream `json:"streams"`
}

type lokiPushStream struct {
	Stream map[string]interface{} `json:"stream"`
	Values [1][2]string           `json:"values"`
}

func (e *lokiEncode) EncodeEntry(ent zapcore.Entry, fields []zap.Field) (*buffer.Buffer, error) {
	final := e.clone()
	push := &lokiPush{
		Streams: [1]lokiPushStream{
			{
				Stream: nil,
				Values: [1][2]string{{
					utils.ToString(ent.Time.UnixNano()),
					ent.Message,
				}},
			},
		},
	}

	for _, v := range fields {
		v.AddTo(final)
	}

	push.Streams[0].Stream = final.Fields

	buf := e.pool.Get()
	b, err := json.Marshal(push)
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(b); err != nil {
		return nil, err
	}
	return buf, nil
}
