package main

import (
	"context"
	"log"
	"time"

	userpb "github.com/BillyLouis/grpc_with_mtls/cmd/grpc_with_mtls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/tls/certprovider/pemfile"
	"google.golang.org/grpc/security/advancedtls"
)

const (
	clientCredRefreshInterval = 1 * time.Minute
	clientCAPath              = "/etc/all_mycerts/mysql/certs/ca/ca.pem"
	clientCertPath            = "/etc/all_mycerts/mysql/certs/client/client-cert.pem"
	clientKeyPath             = "/etc/all_mycerts/mysql/certs/client/client-key.pem"
	grpcServerAddr            = "10.39.8.55:50051"
)

func createMTLSClientConn() (*grpc.ClientConn, error) {
	identityProvider, err := pemfile.NewProvider(pemfile.Options{
		CertFile:        clientCertPath,
		KeyFile:         clientKeyPath,
		RefreshDuration: clientCredRefreshInterval,
	})
	if err != nil {
		return nil, err
	}

	rootProvider, err := pemfile.NewProvider(pemfile.Options{
		RootFile:        clientCAPath,
		RefreshDuration: clientCredRefreshInterval,
	})
	if err != nil {
		return nil, err
	}

	creds, err := advancedtls.NewClientCreds(&advancedtls.Options{
		IdentityOptions: advancedtls.IdentityCertificateOptions{
			IdentityProvider: identityProvider,
		},
		RootOptions: advancedtls.RootCertificateOptions{
			RootProvider: rootProvider,
		},
		VerificationType: advancedtls.CertVerification,
	})
	if err != nil {
		return nil, err
	}

	return grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(creds))
}

func startClientStream(client userpb.UserServiceClient) error {
	stream, err := client.SyncUsers(context.Background())
	if err != nil {
		return err
	}

	// Send mock data
	go func() {
		for {
			//This should override by server:
			mockUser := &userpb.User{
				Id:           1001,
				DbName:       "client_user",
				DbPinHash:    "SomeBcryptPasswordHash_48734743",
				DbUnshashPin: "SomeBcryptReturnCheckHash_7447644",
				DbEmail:      "user@example.com",
				CardAtr:      "ATR_32-32-32-32",
				SerialNumber: "SN:123456",
				DbUserSign:   false,
				CreatedAt:    time.Now().Format(time.RFC3339),
			}
			if err := stream.Send(mockUser); err != nil {
				log.Printf("Send error: %v", err)
				return
			}
			log.Println("Sent user:", mockUser.DbName)
			time.Sleep(2 * time.Second)
		}
	}()

	// Receive from server
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("Receive error: %v", err)
			break
		}
		log.Printf("Received from server: %s (%s)", in.DbName, in.DbEmail)
	}

	return nil
}

func main() {
	conn, err := createMTLSClientConn()
	if err != nil {
		log.Fatalf("failed to create secure client connection: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	if err := startClientStream(client); err != nil {
		log.Fatalf("stream error: %v", err)
	}
}
