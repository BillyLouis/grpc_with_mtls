package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/security/advancedtls"

	userpb "github.com/BilyLouis/grpc_with_mtls/cmd/grpc_with_mtls/proto"
)

const (
	serverCertPath    = "/etc/all_mycerts/server/server-cert.pem"
	serverKeyPath     = "/etc/all_mycerts/server/server-key.pem"
	serverCAPath      = "/etc/all_mycerts/server/ca.pem"
	serverCredRefresh = 1 * time.Minute
	serverListenAddr  = ":50051"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *userServer) SyncUsers(stream userpb.UserService_SyncUsersServer) error {
	// Send mock user to client
	go func() {
		for {
			mock := &userpb.User{
				Id:        1,
				DbName:    "server",
				DbEmail:   "server@demo.com",
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			if err := stream.Send(mock); err != nil {
				log.Println("Send error:", err)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		user, err := stream.Recv()
		if err != nil {
			log.Println("Receive error:", err)
			return err
		}
		log.Printf("Received from client: %s (%s)", user.DbName, user.DbEmail)
	}
}

func createMTLSServerCreds() (grpc.ServerOption, error) {
	identity, err := pemfile.NewProvider(pemfile.Options{
		CertFile:        serverCertPath,
		KeyFile:         serverKeyPath,
		RefreshDuration: serverCredRefresh,
	})
	if err != nil {
		return nil, err
	}

	root, err := pemfile.NewProvider(pemfile.Options{
		RootFile:        serverCAPath,
		RefreshDuration: serverCredRefresh,
	})
	if err != nil {
		return nil, err
	}

	creds, err := advancedtls.NewServerCreds(&advancedtls.Options{
		IdentityOptions:   advancedtls.IdentityCertificateOptions{IdentityProvider: identity},
		RootOptions:       advancedtls.RootCertificateOptions{RootProvider: root},
		RequireClientCert: true,
		VerificationType:  advancedtls.CertVerification,
	})
	if err != nil {
		return nil, err
	}

	return grpc.Creds(creds), nil
}

func startServer() error {
	listener, err := net.Listen("tcp", serverListenAddr)
	if err != nil {
		return err
	}

	tlsOption, err := createMTLSServerCreds()
	if err != nil {
		return err
	}

	server := grpc.NewServer(tlsOption)
	userpb.RegisterUserServiceServer(server, &userServer{})

	log.Println("Secure gRPC server listening on", serverListenAddr)
	return server.Serve(listener)
}

func main() {
	if err := startServer(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
