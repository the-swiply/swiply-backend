package server

import (
	"crypto/md5"
	"encoding/hex"
	"google.golang.org/grpc/metadata"
	"strings"
)

var fingerprintHeaders = []string{
	"sec-ch-ua-platform", "sec-ch-ua",
	"grpcgateway-accept-language", "grpcgateway-user-agent",
	"user-agent",
}

func createFingerprintFromMeta(md metadata.MD) string {
	extractedHeaders := make([]string, len(fingerprintHeaders))
	for i, header := range fingerprintHeaders {
		h := md.Get(header)
		if len(h) > 0 {
			extractedHeaders[i] = h[0]
		}
	}
	fingerprint := strings.Join(extractedHeaders, ":")

	hash := md5.Sum([]byte(fingerprint))

	return hex.EncodeToString(hash[:])
}
