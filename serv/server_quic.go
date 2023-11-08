package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"sync"

	quic "github.com/quic-go/quic-go"
)

const addr = "localhost:4242"

var bufPool = sync.Pool{
	New: func() interface{} {
		// Create a new buffer with the desired size.
		return make([]byte, 1000)
	},
}

func RunServ() {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go handleSession(sess)
	}
}

func handleSession(sess quic.Connection) {
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go handleStream(stream)
	}
}

func handleStream(stream quic.Stream) {
	for {
		// Get a buffer from the pool.
		buf := bufPool.Get().([]byte)

		n, err := stream.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Received:", string(buf[:n]))

		// Put the buffer back into the pool when you're done with it.
		bufPool.Put(buf)
	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
