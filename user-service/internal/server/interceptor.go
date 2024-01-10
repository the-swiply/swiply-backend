package server

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func withLogAndRecover() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		loggy.Errorln("PANIC:", p)
		loggy.Errorln(string(debug.Stack()))
		return status.Errorf(codes.Internal, "%v", p)
	})
}
