package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"my.localhost/funny/confusion_oauth2_httprouter/src/config"
	"my.localhost/funny/confusion_oauth2_httprouter/src/misc"
	"my.localhost/funny/confusion_oauth2_httprouter/src/oauth2"
	"my.localhost/funny/confusion_oauth2_httprouter/src/upload"
	"my.localhost/funny/confusion_oauth2_httprouter/src/auth"
	"my.localhost/funny/confusion_oauth2_httprouter/src/database"
	"my.localhost/funny/confusion_oauth2_httprouter/src/dishes"

	"github.com/julienschmidt/httprouter"

	_ "github.com/go-sql-driver/mysql"
)

const serverListenPort = "0.0.0.0:3000"
const sslPort = ":3000"
const serverListenSslPort = "0.0.0.0" + sslPort
const certPath = "../certs/www.confusion.com.crt"
const keyPath = "../certs/www.confusion.com.key"

func getIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Welcome to Fusion!\n")
}

func listenOnInsecurePortAndRedirect() {

	log.Println("Server starting. Listening on " + serverListenPort)
	// redirect every http request to https
	go http.ListenAndServe(serverListenPort, http.HandlerFunc(redirectToSecurePort))
}

func listenOnSecurePort(router *httprouter.Router) {

	log.Println("Server starting. Listening on " + serverListenSslPort)

	// start server on https port
	server := http.Server{
		Addr:    serverListenSslPort,
		Handler: router,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	certFilePath, _ := filepath.Abs(certPath)
	keyFilePath, _ := filepath.Abs(keyPath)
	err := server.ListenAndServeTLS(certFilePath, keyFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	configFilePath := misc.GetConfigFilePath()
	config := config.ReadDbConfig(configFilePath)

	database.SetupDatabase(config)

	router := httprouter.New()
	router.GET("/", getIndex)
	
	dishes.SetupRoutes(router)
	auth.SetupRoutes(router)
	upload.SetupRoutes(router, config)
	oauth2.SetupRoutes(router, config)

	setupDefaultRoutes(router)

	listenOnSecurePort(router)
}
