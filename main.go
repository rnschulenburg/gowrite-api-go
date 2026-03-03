package main

import (
	"github.com/joho/godotenv"
	"github.com/rnschulenburg/gowrite-api-go/App/Services/AiService"
	"github.com/rnschulenburg/gowrite-api-go/Core/Ws"
	"github.com/rnschulenburg/gowrite-api-go/Package/DbConnection"
	"github.com/rnschulenburg/gowrite-api-go/routers"
	"github.com/rnschulenburg/gowrite-api-go/routers/auth"
	"log"
	"net/http"
	"os"
)

func main() {

	setEnvironment()
	auth.InitAuth()
	AiService.InitAi()
	DbConnection.InitDB()

	defer DbConnection.CloseDB()

	apiPort := os.Getenv("ApiPort")
	routerApi := routers.InitRoutes()
	handlerApi := auth.CorsHandler(routerApi)
	// httpProtocol := os.Getenv("HttpProtocol")

	mux := http.NewServeMux()
	mux.Handle("/", handlerApi)
	mux.HandleFunc("/ws", Ws.WebSocketHandler)

	log.Println("gowrite listening on Port: " + apiPort)

	//if httpProtocol == "http" {

	log.Fatal(http.ListenAndServe("0.0.0.0:"+apiPort, mux))

	//} else if httpProtocol == "https" {
	//
	//	certFile, keyFile := getCert()
	//	log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+apiPort, certFile, keyFile, mux))
	//
	//} else {
	//	panic("enter http or https for .env:HttpProtocol")
	//}
}

func setEnvironment() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}

//func getCert() (certFile string, keyFile string) {
//	certFile = os.Getenv("CertFile")
//	keyFile = os.Getenv("KeyFile")
//
//	_, err2 := os.Open(certFile)
//	if err2 != nil {
//		log.Println("certFile.csr not found in path: " + certFile)
//		panic(err2)
//	}
//
//	_, err3 := os.Open(keyFile)
//	if err3 != nil {
//		log.Println("keyFile.key not found in path: " + keyFile)
//		panic(err3)
//	}
//
//	return certFile, keyFile
//}
