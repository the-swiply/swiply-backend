package grut

import (
	"fmt"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(msg string, err error) error {
	errMsg := fmt.Sprintf("%s: %v", msg, err)
	loggy.Errorln(errMsg)
	return status.Error(codes.Internal, errMsg)
}
