package http

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/davidbelliott/nyx/config"
	"github.com/davidbelliott/nyx/http/admin"
	"github.com/davidbelliott/nyx/http/board"
	"github.com/davidbelliott/nyx/http/errw"
	"github.com/davidbelliott/nyx/http/middle"
	"net/http"
	"time"
)

func Start(config *config.Config) error {
	err := admin.LoadTemplates()
	if err != nil {
		return err
	}
	err = board.LoadTemplates()
	if err != nil {
		return err
	}
	err = errw.LoadTemplates()
	if err != nil {
		return err
	}
	middle.SetupSessionManager(config)

	r := chi.NewRouter()

	fmt.Println("Setting up Router")
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(middleware.ThrottleBacklog(1000, 6000, 10*time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middle.LimitSize(config))
	r.Use(middleware.DefaultCompress)

	r.Use(middle.ConfigCtx(config))

	r.Use(middle.CSRFProtect)
	{
		mw, err := middle.Database(config)
		if err != nil {
			return err
		}
		r.Use(mw)
	}

	r.Route("/admin/", admin.AdminRouter)
	r.Route("/mod/", admin.ModRouter)
	{
		box, err := rice.FindBox("res/")
		if err != nil {
			return err
		}
		atFileServer := http.StripPrefix("/@/", http.FileServer(box.HTTPBox()))
		r.Mount("/@/", atFileServer)
	}
	r.Group(board.Router)

	fmt.Println("Setup Complete, Starting Web Server...")
	srv := &http.Server{
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		Handler:        r,
		Addr:           config.ListenOn,
		MaxHeaderBytes: 1 * 1024 * 1024,
	}
	return srv.ListenAndServe()
}
