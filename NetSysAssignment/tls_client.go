package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	conn, err := tls.Dial("tcp", "google.com:443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	state := conn.ConnectionState()

	fmt.Printf("TLS version: %s\n", tlsVersion(state.Version))
	fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(state.CipherSuite))

	cert := state.PeerCertificates[0]
	fmt.Printf("Issuer Organization: %s\n", cert.Issuer.Organization)
}

func tlsVersion(version uint16) string {
	switch version {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "Unknown"
	}
}
