package main

import (
	serv "caffeine/serv"
	cl "caffeine/client"
)

func main() {
	go serv.RunServ()
	cl.RunClient()
	// quic.RunClient()()
}
