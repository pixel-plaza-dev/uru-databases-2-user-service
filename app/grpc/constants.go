package grpc

const (
	// AuthServiceUriKey is the key of the default URI for the application
	AuthServiceUriKey = "AUTH_SERVICE_HOST"

	// KeyFilePath is the path to the key file
	KeyFilePath = "certificates/server-key.pem"

	// CertificateFilePath is the path to the certificate file
	CertificateFilePath = "certificates/server-cert.pem"

	// AuthServiceCaPath is the path to the CA certificate for the Auth service
	AuthServiceCaPath = "certificates/ca-cert.pem"
)
