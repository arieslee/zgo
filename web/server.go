package web

import (
	"net/http"
	"regexp"
	"time"

	"github.com/JoveYu/zgo/log"
)

const (
	PoweredBy string = "zgo/0.0.1"
)

type Server struct {
	Addr    string
	Routers []Router
	Logger  *log.LevelLogger
	Charset string
}

func NewServer() *Server {
	server := Server{
		Logger: log.DefaultLog,
	}
	return &server
}

func (s *Server) Router(method string, path string, handler ContextHandlerFunc) {
	cr, err := regexp.Compile(path)
	if err != nil {
		s.Logger.Warn("can not add route [%s] %s", path, err)
		return
	}

	s.Routers = append(s.Routers, Router{
		r:       path,
		cr:      cr,
		method:  method,
		handler: handler,
	})

}

func (s *Server) ProcessHandler(w http.ResponseWriter, r *http.Request) {
	tstart := time.Now()

	ctx := NewContext(w, r)
	ctx.Params = make(map[string]string)

	path := ctx.URL().Path

	// default header
	ctx.SetHeader("X-Powered-By", PoweredBy)
	ctx.SetContextType("text/plain")

	for _, router := range s.Routers {
		if ctx.Method() != router.method {
			continue
		}
		if !router.cr.MatchString(path) {
			continue
		}
		match := router.cr.FindStringSubmatch(path)
		if len(match[0]) != len(path) {
			continue
		}

		if len(match) > 1 {
			for idx, name := range router.cr.SubexpNames()[1:] {
				ctx.Params[name] = match[idx+1]
			}
			log.Debug(ctx.Params)
		}

		router.handler(ctx)
		s.LogRequest(tstart, ctx)
		return
	}

	ctx.Abort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	s.LogRequest(tstart, ctx)
}

func (s *Server) Run(addr string) {
	http.HandleFunc("/", s.ProcessHandler)
	http.ListenAndServe(addr, nil)
}

func (s *Server) LogRequest(tstart time.Time, ctx Context) {
	if s.Logger == nil {
		s.Logger = log.DefaultLog
	}

	s.Logger.Info("%d|%s|%s|%s|%s|%d",
		ctx.ResponseWriter.status, ctx.Method(), ctx.URL().Path,
		ctx.Query().Encode(), ctx.ClientIP(),
		time.Since(tstart)/time.Microsecond,
	)
}