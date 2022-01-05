package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http/middleware"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
	sqlStorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	UpdateEvent(ctx context.Context, event app.Event) error
	GetAllEventsOfUser(ctx context.Context, userID int64) ([]app.Event, error)
}

type Server struct {
	conf       config.ServerConf
	log        Logger
	app        Application
	httpServer *http.Server
}

func NewServer(config config.ServerConf, logger Logger, app Application) *Server {
	return &Server{
		conf: config,
		log:  logger,
		app:  app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:    s.conf.Addr,
		Handler: handler(s.log),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Stop(ctx)
		if err != nil {
			s.log.Error("failed to stop server")
		}
	}()

	s.log.Info(fmt.Sprintf("listening: %s", s.conf.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handler(logger Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		st := sqlStorage.New(config.SQLStorage{
			DSN:          "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		})
		err := st.Connect(context.Background())
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		title := r.FormValue("title")
		desc := r.FormValue("desc")
		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		begin, err := time.Parse(time.RFC822, r.FormValue("begin"))
		if err != nil {
			_, _ = io.WriteString(w, r.FormValue("begin")+"\n")
			_, _ = io.WriteString(w, err.Error())
			return
		}
		end, err := time.Parse(time.RFC822, r.FormValue("end"))
		if err != nil {
			_, _ = io.WriteString(w, r.FormValue("end")+"\n")
			_, _ = io.WriteString(w, err.Error())
			return
		}

		event := storage.Event{
			Title:  title,
			Desc:   desc,
			UserID: int64(userID),
			Begin:  begin,
			End:    end,
		}

		err = st.CreateEvent(context.Background(), event)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		_, _ = io.WriteString(w, "success \n")
	})

	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		st := sqlStorage.New(config.SQLStorage{
			DSN:          "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		})
		err := st.Connect(context.Background())
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		title := r.FormValue("title")
		desc := r.FormValue("desc")
		begin, err := time.Parse(time.RFC822, r.FormValue("begin"))
		if err != nil {
			_, _ = io.WriteString(w, r.FormValue("begin")+"\n")
			_, _ = io.WriteString(w, err.Error())
			return
		}
		end, err := time.Parse(time.RFC822, r.FormValue("end"))
		if err != nil {
			_, _ = io.WriteString(w, r.FormValue("end")+"\n")
			_, _ = io.WriteString(w, err.Error())
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))

		event := storage.Event{
			ID:    int64(id),
			Title: title,
			Desc:  desc,
			Begin: begin,
			End:   end,
		}

		err = st.UpdateEvent(context.Background(), event)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
	})

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		st := sqlStorage.New(config.SQLStorage{
			DSN:          "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		})
		err := st.Connect(context.Background())
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))

		err = st.DeleteEvent(context.Background(), int64(id))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
	})

	mux.HandleFunc("/user-events", func(w http.ResponseWriter, r *http.Request) {
		st := sqlStorage.New(config.SQLStorage{
			DSN:          "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		})
		err := st.Connect(context.Background())
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		userID, _ := strconv.Atoi(r.FormValue("user_id"))

		events, err := st.GetAllEventsOfUser(context.Background(), int64(userID))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		for _, event := range events {
			fmt.Fprintf(w, "%v \n", event)
		}
	})

	return middleware.NewLoggerMiddleware(logger, mux)
}
