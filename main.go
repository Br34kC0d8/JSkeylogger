package main


import (
	"flag"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool{return true},
	}
	listenAddr string
	wsAddr string
	jsTemplate *template.Template
)

func init()  {
	flag.StringVar(&listenAddr, "listen-addr", "", "Address to listen on")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for Websocket connection")
	flag.Parse()
	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err!= nil{
		panic(err)
	}
}
func serveWS(res http.ResponseWriter, req *http.Request){
	conn, err := upgrader.Upgrade(res,req,nil)
	if err != nil{
		http.Error(res, "server error", 500)
		return
	}
	defer conn.Close()
	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())
	for{
		_,msg,err:= conn.ReadMessage()
		if err != nil{
			return
		}
		fmt.Printf("From %s: %s\n", conn.RemoteAddr().String(), string(msg))
	}
}

func serveFile(res http.ResponseWriter, r *http.Request)  {
	res.Header().Set("Content-Type", "application/javascript")
	jsTemplate.Execute(res,wsAddr)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws",serveWS)
	r.HandleFunc("/k.js", serveFile)
	log.Fatal(http.ListenAndServe(":8080",r))
}