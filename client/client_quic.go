package quic

import (
	"context"
	"crypto/tls"
	"time"

	// "fmt"
	"log"

	quic "github.com/quic-go/quic-go"
)

const address = "localhost:4242"

func RunClient() {
	session, err := quic.DialAddr(context.Background(), address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for i:=0; i<100; i++{
		time.Sleep(1 * time.Second)
	_, err = stream.Write([]byte("Hello, server!!!"))
	if err != nil {
		log.Fatal(err)
	}
}

	// buf := make([]byte, 100)
	// n, err := stream.Read(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Received:", string(buf[:n]))
}
