package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/server/http/middleware"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
	sqlStorage "github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

type Server struct {
	conf       config.ServerConf
	log        *zap.Logger
	app        Application
	httpServer *http.Server
}

type Application interface {
	CreateEvent(ctx context.Context, event app.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	UpdateEvent(ctx context.Context, event app.Event) error
	GetAllEvents(ctx context.Context) ([]app.Event, error)
	GetAllEventsOfUser(ctx context.Context, userID int64) ([]app.Event, error)
}

func NewServer(config config.ServerConf, logger *zap.Logger, app Application) *Server {
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

	s.log.Info("listening", zap.String("addr", s.conf.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func handler(logger *zap.Logger) http.Handler {
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
		} else {
			_, _ = io.WriteString(w, "connected")
		}

		title := r.FormValue("title")
		desc := r.FormValue("desc")
		userId, _ := strconv.Atoi(r.FormValue("user_id"))
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
			UserID: int64(userId),
			Begin:  begin,
			End:    end,
		}

		err = st.CreateEvent(context.Background(), event)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		_, _ = io.WriteString(w, "success")
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
		} else {
			_, _ = io.WriteString(w, "connected")
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
		_, _ = io.WriteString(w, "success")
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
		} else {
			_, _ = io.WriteString(w, "connected")
		}

		id, _ := strconv.Atoi(r.FormValue("id"))

		err = st.DeleteEvent(context.Background(), int64(id))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		_, _ = io.WriteString(w, "success")
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
		} else {
			_, _ = io.WriteString(w, "connected \n")
		}

		userId, _ := strconv.Atoi(r.FormValue("user_id"))

		events, err := st.GetAllEventsOfUser(context.Background(), int64(userId))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		for _, event := range events {
			fmt.Fprintf(w, "%v \n", event)
		}
	})

	mux.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		st := sqlStorage.New(config.SQLStorage{
			DSN:          "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable",
			MaxOpenConns: 2,
			MaxIdleConns: 2,
		})
		err := st.Connect(context.Background())
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		} else {
			_, _ = io.WriteString(w, "connected \n")
		}

		events, err := st.GetAllEvents(context.Background())
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
