package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx    context.Context
	cancel func()
	svrs   []http.Server
}

func New(ctx context.Context, addr2Handler map[string]http.Handler) *App {
	servers := []http.Server{}
	for addr, handler := range addr2Handler {
		newSrv := http.Server{
			Addr:    addr,
			Handler: handler,
		}
		servers = append(servers, newSrv)
	}
	defaultCtx := context.Background()
	ctx, cancel := context.WithTimeout(defaultCtx, 10*time.Second)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		svrs:   servers,
	}
}

func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	for _, svr := range a.svrs {
		svr := svr
		g.Go(func() error {
			<-ctx.Done()
			return svr.Shutdown(context.TODO())
		})

		g.Go(func() error {
			return svr.ListenAndServe()
		})
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-sig:
				a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func main() {
	addr2Handler := make(map[string]http.Handler)
	addr2Handler[":8080"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "hello")
	})

	defaultCtx := context.Background()
	app := New(defaultCtx, addr2Handler)
	fmt.Println(app.Run())
}
