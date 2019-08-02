package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/jinzhu/configor"
	"github.com/nats-io/go-nats"
	"github.com/paulvasilenko/nats-http-adapter/handler"
	internal_nats "github.com/paulvasilenko/nats-http-adapter/nats"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	HTTP struct {
		Port string `yaml:"Port" default:"80"`
	} `yaml:"HTTP"`
	NATS struct {
		Endpoint       string        `yaml:"Endpoint" required:"true"`
		RequestTimeout time.Duration `yaml:"RequestTimeout" default:"1s"`
	} `yaml:"NATS"`
}

var configPath = flag.String("c", "config.yaml", "path to config file")

func main() {
	flag.Parse()
	ctx := context.Background()

	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	conf := Config{}
	err := configor.
		New(&configor.Config{
			ErrorOnUnmatchedKeys: true,
		}).
		Load(&conf, *configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Infof("Connecting to NATS %v", conf.NATS.Endpoint)

	natsConn, err := nats.Connect(conf.NATS.Endpoint)
	if err != nil {
		log.Fatalf("failed to open nats conn: %v", err)
	}

	log.Info("Connected to NATS")

	service := handler.Service{
		NATS: internal_nats.Conn{
			Timeout: conf.NATS.RequestTimeout,
			Conn:    natsConn,
		},
	}

	http.HandleFunc("/nats", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Receiving request %v", r)
		w.Header().Add("Content-Type", "application/json")
		if r.Body == nil {
			processErr(handler.BadRequest("empty request body"), w)
			return
		}

		req := &handler.NATSMessage{}
		rawBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			processErr(err, w)
			return
		}

		if len(rawBody) == 0 {
			processErr(handler.BadRequest("empty request body"), w)
			return
		}

		if err := json.Unmarshal(rawBody, req); err != nil {
			processErr(handler.BadRequest(errors.Wrap(err, "bad body structure").Error()), w)
			return
		}

		res, err := service.SendNATSMessage(ctx, req)
		if err != nil {
			processErr(err, w)
			return
		}

		if res != nil {
			rawRespBody, err := json.Marshal(res)
			if err != nil {
				processErr(err, w)
				return
			}
			w.Write(rawRespBody)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Infof("Starting HTTP service on port %v", conf.HTTP.Port)
	if err := http.ListenAndServe(":"+conf.HTTP.Port, http.DefaultServeMux); err != nil {
		log.Fatalf("failed to run fileserver: %v", err)
	}
}

func processErr(err error, w http.ResponseWriter) {
	if httpErr, ok := err.(handler.HTTPError); ok {
		if httpErr.Code() != 0 {
			w.WriteHeader(httpErr.Code())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}
