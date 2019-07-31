package websvr

import (
	"context"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

type Config struct {
	Host      string
	PProfHost string
	Thumbs    string
}

var logger = log.New(os.Stderr, "[websvr]", log.Lshortfile|log.LstdFlags)

func HTTPServer(ctx context.Context, config Config) error {
	if config.Host == "" {
		return errors.New("host cannot be empty")
	}
	mux, e := newServerMux(&config)
	if e != nil {
		return e
	}
	svr := http.Server{
		Addr:    config.Host,
		Handler: mux,
	}
	if config.PProfHost != "" {
		go func() {
			if e := http.ListenAndServe(config.PProfHost, nil); e != nil {
				logger.Printf("Error Serve PProf Server: %s\n", e)
			}
		}()
	}
	go func() {
		logger.Printf("Serve HTTP Server: %s\n", config.Host)
		if e := svr.ListenAndServe(); e != nil {
			logger.Printf("Error Serve HTTP Server: %s\n", e)
		}
	}()
	<-ctx.Done()
	c, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if e := svr.Shutdown(c); e != nil {
		return e
	}
	return nil
}

const thumbsPrefix = "/thumbs/"

func newServerMux(config *Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	replaceThumbs := false
	if config.Thumbs != "" {
		fi, e := os.Stat(config.Thumbs)
		if e != nil {
			return nil, e
		}
		if !fi.IsDir() {
			return nil, errors.New("Config.Thumbs should be a directory")
		}
		replaceThumbs = true
		mux.Handle(thumbsPrefix, http.StripPrefix(thumbsPrefix, http.FileServer(http.Dir(config.Thumbs))))
		logger.Printf("Enable local thumbs cache\n")
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/":
			galleries(replaceThumbs, config.Thumbs)(writer, request)
		default:
			http.NotFound(writer, request)
		}
	})
	return mux, nil
}
