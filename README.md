# grpc_with_mtls
Client and server gRPC endpoints with mutual TLS.  
Specifically, a sample of grpc protocol to use mutual TLS (mTLS) to authenticate & secure both end of the communication.


# Certificate Creation
```shell
#### Client Certificate creation:
 To do ...
 - Create Docker Network (client container)
 - Create  client x509 TLS certificates that run in a Docker container

#### Server Certificate Creation:
 To do ...
 - Create Docker Network (server container)
 - Create  Server x509 TLS certificates that run in a Docker container
```

# Proto Buff Creation
```shell
 To do ...
 - Install libs/packages to generate protobuff

```

# Database Setup and Configuration
```shell
 To do ...
 - Create Docker Network (database container)
 - MariaDB configuration


```

### Output of a client & server mTLS exchange:
```shell
When the certificates of verified on both client & server, the connection is established and data exchanges can be done.
If even one character is changed from the certificates the connection cannot be establish to exchange the secured endpoints.
```


# Authors
- [Billy Louis](): TCP/IP connection between Client and Server using Golang RPC (GRPC)


# Badges
Hardware Team: [NSAL.com](https://NSAL.com/)

[![NSA License](https://img.shields.io/badge/License-NSAL-green.svg)](https://choosealicense.com/licenses/nsal/)
