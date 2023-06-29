package token_auth

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/cache/cache"
)

type AccessToken struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	KvMap        map[string]string `json:"kv_map"`
	ExpireAt     int64             `json:"expire_at"`
	RefreshAt    int64             `json:"refresh_at"`
}

type Storage interface {
	Save(ctx context.Context, accessToken *AccessToken) error
	FindByAccessToken(ctx context.Context, accessToken string) (*AccessToken, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*AccessToken, error)
	Delete(ctx context.Context, accessToken *AccessToken) error
}

type cacheStorage struct {
	cache  cache.Cache
	prefix string
}

func NewCacheStorage(cache cache.Cache, prefix string) Storage {
	return &cacheStorage{
		cache:  cache,
		prefix: prefix,
	}
}

func (r *cacheStorage) key(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

func (r *cacheStorage) Serialize(token *AccessToken) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(token)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

func (r *cacheStorage) Deserialize(d []byte, token *AccessToken) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&token)
}

func (r *cacheStorage) Save(ctx context.Context, accessToken *AccessToken) error {
	b, err := r.Serialize(accessToken)
	if err != nil {
		return err
	}
	err = r.cache.Set(ctx, r.key(accessToken.AccessToken), b, cache.Expiration(
		time.Duration(accessToken.RefreshAt-time.Now().Unix())*time.Second,
	)).Err()
	if err != nil {
		return err
	}
	// 保存refresh_token->access_token的映射
	err = r.cache.Set(ctx, r.key("ref:"+accessToken.RefreshToken), []byte(accessToken.AccessToken), cache.Expiration(
		time.Duration(accessToken.RefreshAt-time.Now().Unix())*time.Second,
	)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *cacheStorage) FindByAccessToken(ctx context.Context, accessToken string) (*AccessToken, error) {
	data, err := r.cache.Get(ctx, r.key(accessToken)).Bytes()
	if err != nil {
		if cache.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	token := &AccessToken{}
	err = r.Deserialize(data, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *cacheStorage) FindByRefreshToken(ctx context.Context, refreshToken string) (*AccessToken, error) {
	accessToken, err := r.cache.Get(ctx, r.key("ref:"+refreshToken)).Result()
	if err != nil {
		if cache.IsNil(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.FindByAccessToken(ctx, accessToken)
}

func (r *cacheStorage) Delete(ctx context.Context, accessToken *AccessToken) error {
	err := r.cache.Del(ctx, r.key(accessToken.AccessToken))
	if err != nil {
		return err
	}
	err = r.cache.Del(ctx, r.key("ref:"+accessToken.RefreshToken))
	if err != nil {
		return err
	}
	return nil
}
