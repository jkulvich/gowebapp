package main

import (
	"GOFirst/app"
	"GOFirst/server"
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	// reads the app flags
	flagConfigName := flag.String("c", "./config.yaml", "config file name")
	flag.Parse()

	// trying to read the config
	conf, err := ReadConfig(*flagConfigName)
	if err != nil {
		logrus.Fatalf("can't read config: %s", err)
	}

	// logger configuring
	lg, err := NewLogger(conf.Logger)
	if err != nil {
		logrus.Fatalf("can't configure logger: %s", err)
	}

	// the global app context
	ctx := context.Background()

	// new server instance (backend of the WebApp)
	serv := server.NewServer(ctx, lg, conf.Server)
	if err := serv.Run(); err != nil {
		lg.Fatalf("can't run the server: %s", err)
	}

	// chanel of signals for stopping the WebApp
	stopChan := make(chan os.Signal, 1)

	// new instance and run the WebApp (frontend)
	app := app.NewApp(ctx, lg, conf.App)
	go func() {
		if err := app.Run(); err != nil {
			lg.Fatalf("can't run the app: %s", err)
		}
		stopChan <- nil
	}()

	// waiting for stop signal
	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<- stopChan

	// send stop signal to frontend
	app.Close()

	// send stop signal to backend
	if err := serv.Stop(); err != nil {
		lg.Fatalf("can't stop the server: %s", err)
	}
}
