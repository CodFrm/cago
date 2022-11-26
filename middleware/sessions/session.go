package sessions

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	gorillaSession "github.com/gorilla/sessions"
)

const errorFormat = "[sessions] ERROR! %s\n"

type Store interface {
	sessions.Store
	Refresh(r *http.Request, name string, session *gorillaSession.Session) error
}

type Session interface {
	sessions.Session
	Refresh() error
}

// Middleware 在gin-contrib/sessions上做封装
func Middleware(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(sessions.DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

func Ctx(ctx *gin.Context) Session {
	return sessions.Default(ctx).(Session)
}

var _ Session = (*session)(nil)

type session struct {
	name    string
	request *http.Request
	store   Store
	session *gorillaSession.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) Refresh() error {
	return s.store.Refresh(s.request, s.name, s.Session())
}

func (s *session) ID() string {
	return s.Session().ID
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options sessions.Options) {
	s.written = true
	s.Session().Options = options.ToGorillaOptions()
}

func (s *session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Session() *gorillaSession.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *session) Written() bool {
	return s.written
}
