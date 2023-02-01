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

	"github.com/codfrm/cago/database/cache/cache"
	cagoSession "github.com/codfrm/cago/middleware/sessions"
	ginSession "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type store struct {
	cache   cache.Cache
	options *Options
}

func NewCacheStore(cache cache.Cache, prefix string, opts ...Option) cagoSession.Store {
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
	return &store{
		cache:   cache,
		options: options,
	}
}

func (s *store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *store) New(r *http.Request, name string) (*sessions.Session, error) {
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

func (s *store) key(session *sessions.Session) string {
	return fmt.Sprintf("%s:%s", s.options.prefix, session.ID)
}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s *store) Serialize(ss *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(ss.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[string]interface{}
func (s *store) Deserialize(d []byte, ss *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Values)
}

func (s *store) load(session *sessions.Session) (bool, error) {
	data, err := s.cache.Get(context.Background(), s.key(session)).Bytes()
	if err != nil {
		if err == cache.Nil {
			return false, nil
		}
		return false, err
	}
	if err := s.Deserialize(data, session); err != nil {
		return false, err
	}
	return true, nil
}

func (s *store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
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

func (s *store) save(session *sessions.Session) error {
	age := session.Options.MaxAge
	if age == 0 {
		age = s.options.defaultMaxAge
	}
	b, err := s.Serialize(session)
	if err != nil {
		return err
	}
	if err := s.cache.Set(context.Background(), s.key(session), b, cache.Expiration(
		time.Duration(age)*time.Second)).Err(); err != nil {
		return err
	}
	return nil
}

// delete removes keys from redis if MaxAge<0
func (s *store) delete(session *sessions.Session) error {
	return s.cache.Del(context.Background(), s.key(session))
}

// Options gin-contrib/sessions的选项
func (s *store) Options(options ginSession.Options) {
}

func (s *store) Refresh(r *http.Request, name string, session *sessions.Session) error {
	// 删除原来的
	if err := s.delete(session); err != nil {
		return err
	}
	// 重新生成id
	session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	return nil
}
