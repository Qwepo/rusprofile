package main

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Qwepo/rusprofile/gen/rusprof"
	"github.com/Qwepo/rusprofile/internal"
	"github.com/Qwepo/rusprofile/internal/services"
	"github.com/Qwepo/rusprofile/pkg/logger"
	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed swagger/rusprof.swagger.json
var swg []byte

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	conf, err := internal.NewConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(conf.Logger.Level, conf.Logger.Filename)
	service := services.NewServices(log)
	mux := runtime.NewServeMux()
	// Start GRPCServer
	g.Go(func() error {
		grpcServer := grpc.NewServer()
		rusprof.RegisterRusprofServer(grpcServer, service)
		rpcaddr := fmt.Sprintf("%s:%s", conf.GRPCServer.Addr, conf.GRPCServer.Port)
		listn, err := ListenWithContext(ctx, rpcaddr)
		if err != nil {
			return errors.WithMessage(err, "GRPCListener")
		}

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = rusprof.RegisterRusprofHandlerFromEndpoint(ctx, mux, rpcaddr, opts)
		if err != nil {
			return errors.WithMessage(err, "GRPCRegisterEndpoint")
		}
		log.Info().Msgf("Run GRPC: %s", listn.Addr())
		err = grpcServer.Serve(listn)
		return errors.WithMessage(err, "GRPCServer")
	})
	// Start HTTPServer
	g.Go(func() error {

		gatewayaddr := fmt.Sprintf("%s:%s", conf.GRPCGateway.Addr, conf.GRPCGateway.Port)

		listn, err := ListenWithContext(ctx, gatewayaddr)
		if err != nil {
			return errors.WithMessage(err, "GatewayListener")
		}
		log.Info().Msgf("Run GRPC-gateway: %s", listn.Addr())
		err = http.Serve(listn, mux)
		return errors.WithMessage(err, "GatewayServer")

	})
	// Start swagger
	g.Go(func() error {
		swaggeraddr := fmt.Sprintf("%s:%s", conf.Swagger.Addr, conf.Swagger.Port)
		listn, err := ListenWithContext(ctx, swaggeraddr)
		if err != nil {
			return errors.WithMessage(err, "SwaggerListener")
		}

		http.Handle("/swaggerui/", http.StripPrefix("/swaggerui", swaggerui.Handler(swg)))
		log.Info().Msgf("Run swagger: %s", listn.Addr())
		err = http.Serve(listn, nil)
		return errors.WithMessage(err, "SwaggerServer")

	})
	errExit := errors.New("program exit")
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-ctx.Done():
		case <-c:
		}
		return errExit

	})

	if err = g.Wait(); err != nil && !errors.Is(err, errExit) {
		log.Fatal().Err(err).Send()
	}
}

func ListenWithContext(ctx context.Context, addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()
	return listener, nil
}
