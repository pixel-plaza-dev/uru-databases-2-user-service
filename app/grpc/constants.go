package grpc

const (
	// UserServiceUriKey is the key of the default URI for the user service
	UserServiceUriKey = "USER_SERVICE_HOST"

	// AuthServiceUriKey is the key of the default URI for the auth service
	AuthServiceUriKey = "AUTH_SERVICE_HOST"

	// KeyFilePath is the path to the key file
	KeyFilePath = "certificates/pixel-plaza/user-service/server-key.pem"

	// CertificateFilePath is the path to the certificate file
	CertificateFilePath = "certificates/pixel-plaza/user-service/server-cert.pem"

	// CACertificatePath is the path to the CA certificate
	CACertificatePath = "certificates/pixel-plaza/ca/ca-cert.pem"
)
