package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/the-swiply/swiply-backend/event/pkg/api/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"path"
)

const (
	swaggerFileName = "swagger.json"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(ctx context.Context, cfg HTTPConfig) (*HTTPServer, error) {
	httpMux := http.NewServeMux()

	err := registerGRPCGateway(ctx, httpMux, cfg.GRPCEndpoint)
	if err != nil {
		return nil, fmt.Errorf("can't register grpc gateway handler: %w", err)
	}

	registerSwagger(httpMux, cfg.SwaggerPath)

	srv := &HTTPServer{}
	httpServer := &http.Server{
		Addr:    cfg.ServeAddr,
		Handler: httpMux,
	}

	srv.Server = httpServer

	return srv, nil
}

func registerGRPCGateway(ctx context.Context, mux *http.ServeMux, grpcAddr string) error {
	gwMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}), runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		// Change if custom headers matching is needed.
		return runtime.DefaultHeaderMatcher(key)
	}))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := event.RegisterEventHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("can't register handler for grpc endpoint: %w", err)
	}

	mux.Handle("/", gwMux)

	return nil
}

func registerSwagger(mux *http.ServeMux, swaggerPath string) {
	mux.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(swaggerPath, swaggerFileName))
	})
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(swaggerPath))))
}
