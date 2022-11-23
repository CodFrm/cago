package cache

import (
	"context"
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codfrm/cago/database/cache"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Store struct {
	cache   cache.ICache
	options *Options
}

func NewCacheStore(cache cache.ICache, prefix string, opts ...Option) sessions.Store {
	options := &Options{
		prefix:        "session",
		defaultMaxAge: 86400 * 30,
		refreshTime:   86400,
		sessionOptions: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			Secure:   true,
			HttpOnly: true,
		},
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
		sessionID := ""
		err = securecookie.DecodeMulti(name, c.Value, &sessionID, s.options.codecs...)
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

func (s *Store) load(session *sessions.Session) (bool, error) {
	if err := s.cache.Get(context.Background(), s.key(session), session); err != nil {
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
	if err := s.cache.Set(context.Background(), s.key(session), session, cache.Expiration(
		time.Duration(session.Options.MaxAge)*time.Second)); err != nil {
		return err
	}
	return nil
}

// delete removes keys from redis if MaxAge<0
func (s *Store) delete(session *sessions.Session) error {
	return s.cache.Del(context.Background(), s.key(session))
}
