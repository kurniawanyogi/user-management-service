package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-service/cmd/middleware"
	"user-management-service/common"
	"user-management-service/config"
	"user-management-service/delivery"
)

type IServer interface {
	Serve(ctx context.Context)
}

type Server struct {
	router Router
}

func NewServer(
	common common.IRegistry,
	delivery delivery.IRegistry,
	middleware middleware.IMiddleware,
) *Server {
	return &Server{
		router: NewRouter(
			common,
			delivery,
			middleware,
		),
	}
}

func (s *Server) Serve(ctx context.Context) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cold.AppPort),
		Handler: s.router.Register(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//logger.Error(ctx, "failed init http server", err)
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of config.Hot.ShutDownDelayInSecond seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//logger.Info(ctx, fmt.Sprintf("Shutdown Server in %d seconds", config.Hot.ShutDownDelayInSecond))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Hot.ShutDownDelayInSecond)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		//logger.Fatal(ctx, "Server Shutdown:", logger.Tag{
		//	Key:   "error",
		//	Value: err.Error(),
		//})
	}
	// catching ctx.Done(). timeout of config.Hot.ShutDownDelayInSecond seconds.
	select {
	case <-ctx.Done():
		log.Println(fmt.Sprintf("timeout of %d seconds.", config.Hot.ShutDownDelayInSecond))
	}
	log.Println("Server exiting")
}
