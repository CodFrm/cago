package loki

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type lokiCore struct {
	zapcore.Core
	options *Options
}

func NewLokiCore(ctx context.Context, opt ...Option) zapcore.Core {
	options := &Options{
		level: func(l zapcore.Level) bool {
			return l >= zap.InfoLevel
		},
	}
	for _, o := range opt {
		o(options)
	}
	w := newLokiWriter(ctx, options)
	sync := zapcore.AddSync(w)
	encode := NewLokiEncode(options.labels...)
	return &lokiCore{
		Core: zapcore.NewCore(
			encode,
			sync,
			options.level,
		),
		options: options,
	}
}

func (l *lokiCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if l.Enabled(ent.Level) {
		return ce.AddCore(ent, l)
	}
	return ce
}

type lokiWriter struct {
	ctx     context.Context
	c       *http.Client
	ch      chan []byte
	options *Options
}

func newLokiWriter(ctx context.Context, options *Options) *lokiWriter {
	w := &lokiWriter{
		ctx: ctx,
		c: &http.Client{
			Timeout: time.Second * 2,
		},
		ch:      make(chan []byte, 1024),
		options: options,
	}
	go w.loop()
	return w
}

func (l *lokiWriter) loop() {
	for {
		select {
		case b := <-l.ch:
			req, err := http.NewRequest(http.MethodPost, l.options.url.String(), bytes.NewBuffer(b))
			if err != nil {
				log.Printf("loki push request err: %v", err)
				break
			}
			if l.options.username != "" {
				req.SetBasicAuth(l.options.username, l.options.password)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := l.c.Do(req)
			buf := bytes.NewBuffer([]byte{})
			func() {
				defer resp.Body.Close()
				_, err = buf.ReadFrom(resp.Body)
				if err != nil {
					log.Printf("loki push response err: %v", err)
					return
				}
			}()
			if resp.StatusCode >= 400 {
				log.Printf("loki push err: %v,status code: %v", buf.String(), resp.StatusCode)
			}
		case <-l.ctx.Done():
			break
		}
	}
}

func (l *lokiWriter) Write(b []byte) (int, error) {
	l.ch <- b
	return len(b), nil
}
