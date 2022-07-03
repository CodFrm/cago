package loki

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap/zapcore"
)

type lokiCore struct {
	zapcore.Core
}

func NewLokiCore(ctx context.Context, lokiUrl *url.URL, enab zapcore.LevelEnabler, opt ...Option) (*lokiCore, error) {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}
	w := newLokiWriter(ctx, lokiUrl)
	sync := zapcore.AddSync(w)
	if options.sync != nil {
		sync = zapcore.NewMultiWriteSyncer(options.sync, sync)
	}
	return &lokiCore{
		Core: zapcore.NewCore(
			NewLokiEncode(),
			sync,
			enab,
		),
	}, nil
}

func (l *lokiCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if l.Enabled(ent.Level) {
		return ce.AddCore(ent, l)
	}
	return ce
}

type lokiWriter struct {
	lokiUrl *url.URL
	ctx     context.Context
	c       *http.Client
	ch      chan []byte
}

func newLokiWriter(ctx context.Context, url *url.URL) *lokiWriter {
	w := &lokiWriter{
		ctx:     ctx,
		lokiUrl: url,
		c: &http.Client{
			Timeout: time.Second * 2,
		},
		ch: make(chan []byte, 1024),
	}
	go w.loop()
	return w
}

func (l *lokiWriter) loop() {
	for {
		select {
		case b := <-l.ch:
			resp, err := l.c.Post(l.lokiUrl.String(), "application/json", bytes.NewBuffer(b))
			if err != nil {
				log.Printf("loki push err: %v", err)
				break
			}
			buf := bytes.NewBuffer([]byte{})
			func() {
				defer resp.Body.Close()
				_, err = buf.ReadFrom(resp.Body)
				if err != nil {
					log.Printf("loki push err: %v", err)
					resp.Body.Close()
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
