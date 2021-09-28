package http

import (
	"context"
	"net/http"
)

// ContextFunc is factory of context.Context.
type ContextFunc func() (context.Context, context.CancelFunc)

// ListenAndServe is run http.ListenAndServe and cancellation by context.Context.
func ListenAndServe(ctx context.Context, addr string, handler http.Handler, shutdownContext ContextFunc) error {
	return ListenAndServeServer(ctx, &http.Server{
		Addr:    addr,
		Handler: handler,
	}, shutdownContext)
}

// ListenAndServeServer is run http.Server.ListenAndServe and cancellation by context.Context.
func ListenAndServeServer(ctx context.Context, s *http.Server, shutdownContext ContextFunc) error {
	return listenAndServe(ctx, s, shutdownContext, func(s *http.Server) error {
		return s.ListenAndServe()
	})
}

// ListenAndServeTLS is run http.ListenAndServeTLS and cancellation by context.Context.
func ListenAndServeTLS(ctx context.Context, addr string, certFile string, keyFile string, handler http.Handler, shutdownContext ContextFunc) error {
	return ListenAndServeTLSServer(ctx, &http.Server{
		Addr:    addr,
		Handler: handler,
	}, certFile, keyFile, shutdownContext)
}

// ListenAndServeTLSServer is run http.Server.ListenAndServeTLS and cancellation by context.Context.
func ListenAndServeTLSServer(ctx context.Context, s *http.Server, certFile string, keyFile string, shutdownContext ContextFunc) error {
	return listenAndServe(ctx, s, shutdownContext, func(s *http.Server) error {
		return s.ListenAndServeTLS(certFile, keyFile)
	})
}

func listenAndServe(ctx context.Context, s *http.Server, shutdownContext ContextFunc, serve func(s *http.Server) error) error {
	errorCh := make(chan error, 1)
	go func() {
		defer close(errorCh)
		<-ctx.Done()

		ctx2, cancel := shutdownContext()
		defer cancel()

		errorCh <- s.Shutdown(ctx2)
	}()

	err := serve(s)
	if err != http.ErrServerClosed {
		return err
	}

	shutdownErr, ok := <-errorCh
	if !ok {
		return err
	}
	if shutdownErr != nil {
		return shutdownErr
	}

	return err
}
