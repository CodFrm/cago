package cache

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"
	"time"

	cache2 "github.com/codfrm/cago/database/cache/cache"
	sessions2 "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Store struct {
	cache   cache2.Cache
	options *Options
}

func NewCacheStore(cache cache2.Cache, prefix string, opts ...Option) sessions2.Store {
	options := &Options{
		prefix:        "session",
		defaultMaxAge: 86400 * 30,
		refreshTime:   86400,
		sessionOptions: &sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   86400 * 30,
			Secure:   true,
			HttpOnly: true,
			SameSite: 0,
		},
		codecs: securecookie.CodecsFromPairs(
			[]byte("secret-auth"),
		),
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Store{
		cache:   cache,
		options: options,
	}
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.options.sessionOptions
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.options.codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

func (s *Store) key(session *sessions.Session) string {
	return fmt.Sprintf("%s:%s", s.options.prefix, session.ID)
}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s *Store) Serialize(ss *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[string]interface{}
func (s *Store) Deserialize(d []byte, ss *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Values)
}

func (s *Store) load(session *sessions.Session) (bool, error) {
	data, err := s.cache.Get(context.Background(), s.key(session)).Bytes()
	if err != nil {
		if err == cache2.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	if err := s.Deserialize(data, session); err != nil {
		return false, err
	}
	// 检查session是否快过期
	if session.Values["created"] == nil ||
		time.Now().Unix()-session.Values["created"].(int64) > int64(s.options.refreshTime) {
		// 这里可能会存在并发问题,但是不影响使用
		// 删除原来的session并重新生成一个新的session
		if err := s.delete(session); err != nil {
			return false, err
		}
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}
	return true, nil
}

func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.options.codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *Store) save(session *sessions.Session) error {
	age := session.Options.MaxAge
	if age == 0 {
		age = s.options.defaultMaxAge
	}
	session.Values["created"] = time.Now().Unix()
	b, err := s.Serialize(session)
	if err != nil {
		return err
	}
	if err := s.cache.Set(context.Background(), s.key(session), b, cache2.Expiration(
		time.Duration(age)*time.Second)).Err(); err != nil {
		return err
	}
	return nil
}

// delete removes keys from redis if MaxAge<0
func (s *Store) delete(session *sessions.Session) error {
	return s.cache.Del(context.Background(), s.key(session))
}

// Options gin-contrib/sessions的选项
func (s *Store) Options(options sessions2.Options) {
}
