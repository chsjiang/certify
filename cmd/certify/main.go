package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const usage = `Usage of certify:
certify [flag] [ip-or-dns-san] [cn:default certify]

$ certify -init
⚡️ Initialize new CA Certificate and Key

$ certify server.local 172.17.0.1
⚡️ Generate certificate with alt name server.local and 172.17.0.1

Also you can set subject common name by providing cn:yourcn

$ certify cn:web-server
⚡️ Generate certificate with common name web-server

You must create new CA by run -init before you can create certificate.
`

var (
	caPath    = "ca-cert.pem"
	caKeyPath = "ca-key.pem"
)

func main() {
	init := flag.Bool("init", false, "initialize new CA Certificate and Key")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}
	flag.Parse()

	if *init {
		pkey, err := generatePrivateKey(caKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("CA private key file generated", caKeyPath)

		if err := generateCA(pkey.PrivateKey, os.Args[2], caPath); err != nil {
			log.Fatal(err)
		}
		fmt.Println("CA certificate file generated", caPath)
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("you must provide at least two argument.\n\n")
		fmt.Fprint(flag.CommandLine.Output(), usage)
		os.Exit(1)
	}

	if !isExist(caPath) || !isExist(caKeyPath) {
		log.Fatal("error CA Certificate or Key is not exists, run -init to create it.")
	}

	keyPath := getFilename(os.Args, true)

	pkey, err := generatePrivateKey(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Private key file generated", keyPath)

	if err := generateCert(pkey.PrivateKey, os.Args); err != nil {
		log.Fatal(err)
	}
}
