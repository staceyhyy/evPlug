package main

import (
	"context"
	"golive/app"
	"golive/data"
	"golive/envvar"
	"golive/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//load env parameter
	env, err := envvar.Load()
	if err != nil {
		log.Fatal(err)
	}

	//setting logfile
	logFile, err := logger.NewLogger(env.LogFile)
	if err != nil {
		logger.Fatal.Println("Failed to open error log file: ", err)
	}
	defer logFile.Close()

	//connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := data.Connect(ctx, env)
	if err != nil {
		logger.Fatal.Println("Failed to open error log file: ", err)
	}
	defer db.Client.Disconnect(ctx)

	r := app.SetRouter(db) //register handlers
	s := app.NewServer(r, env.Address)
	logger.Info.Println("server started on ", env.Address)
	go func() {
		err := s.ListenAndServeTLS("cert/cert.pem", "cert/key.pem")
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal.Println("error starting the server: ", err)
			os.Exit(1)
		}
	}()
	waitForShutdown(s) //handle proper shutdown
}

//waitForShutdown shutdown the server properly by checking the sigterm or interupt.
//it will wait for 30 seconds for current operation to complete before proceeding with the shutdown
func waitForShutdown(s *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info.Println("received signal on server: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	logger.Info.Println("shutting down")
	os.Exit(0)
}
