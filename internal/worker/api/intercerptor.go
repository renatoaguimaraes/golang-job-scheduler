package api

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryAuthInterceptor intercept unary calls to authorize the user
// based on certification extension oid 1.2.840.10070.8.1.
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx, info.FullMethod); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(ctx, req)
}

// StreamAuthInterceptor intercept stream calls to authorize the user
// based on certification extension oid 1.2.840.10070.8.1.
func StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(stream.Context(), info.FullMethod); err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(srv, stream)
}

// authorize verifies the user information given by certificate
// against the mapped roles for a specific method.
func authorize(ctx context.Context, method string) error {
	// reads the peer information from context
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("error to read peer information")
	}
	// reads user tls inforation
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return errors.New("error to get auth information")
	}
	// access the leaf certificate to get user roles
	certs := tlsInfo.State.VerifiedChains
	if len(certs) == 0 || len(certs[0]) == 0 {
		return errors.New("missing certificate chain")
	}
	// find user roles from certificate extensions
	var roles []string
	for _, ext := range certs[0][0].Extensions {
		if oid := OidToString(ext.Id); IsOidRole(oid) {
			roles = ParseRoles(string(ext.Value))
			break
		}
	}
	// check user permissions to execute a specific method
	if !HasPermission(method, roles) {
		return errors.New("unauthorized, user does not have privileges enough")
	}
	return nil
}
