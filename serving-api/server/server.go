package server

import (
	"time"
	"crypto/tls"
	"github.com/gorilla/handlers"
	"net/http"
	"github.com/gorilla/mux"
)

//New Server
func New(smux *mux.Router, bindAddress string, hts bool) *http.Server{
	//cors
	corsH := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	// tls Configuration
	tlsConfig := & tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, 
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},

	}
	//basic server
	var server *http.Server
	if hts==true {
		server = &http.Server{

			Addr:         bindAddress,
			Handler:      corsH(smux),
			TLSConfig:    tlsConfig,
			ReadTimeout:  8 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  300 * time.Second,
		}
	}else{
		server = &http.Server{

			Addr:         bindAddress,
			Handler:      corsH(smux),
			ReadTimeout:  8 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  300 * time.Second,
		}
	}


	return server
}