package routers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

// Router echo
type Router struct {
	Echo *echo.Echo
	Name string
}

// Init init func
func (r *Router) Init() {

	r.Echo = echo.New()
	r.Echo.HideBanner = true
}

// RegisterPreMiddleware registers pre middleware
func (r *Router) RegisterPreMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Pre(middleware)
}

// RegisterMiddleware registers midlleware
func (r *Router) RegisterMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Use(middleware)
}

// Start starts router
func (r *Router) Start(host string, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		r.Echo.Logger.Info(fmt.Sprintf("Starting %s server on port: %s", r.Name, port))
		if err := r.Echo.Start(host + ":" + port); err != nil && err != http.ErrServerClosed {
			r.Echo.Logger.Fatal(err)
			r.Echo.Logger.Fatal(xconstants.MSG_SERVER_SHUTTING_DOWN)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := r.Echo.Shutdown(ctx); err != nil {
		r.Echo.Logger.Fatal(err)
	}
}
