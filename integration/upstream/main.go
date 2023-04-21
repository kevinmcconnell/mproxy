package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DemoHandler struct {
	http.Handler
}

func NewDemoHandler() *DemoHandler {
	mux := http.NewServeMux()
	h := &DemoHandler{
		mux,
	}

	mux.HandleFunc("/cable", h.WebsocketHandler)
	mux.HandleFunc("/", h.IndexHandler)

	return h
}

func (h *DemoHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	fmt.Printf("host: %s request: %s\n", host, r.URL)
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `
<html>
  <head>
		<title>Hello world</title>
	</head>
	<body>
		<h1>Hello from %s</h1>
		<p id="now"></p>

		<script>
			function connect() {
				const socket = new WebSocket("ws://" + location.host + "/cable")
				socket.onclose = function() {
					setTimeout(connect, 1000)
				}
				socket.onmessage = function(e) {
					document.getElementById("now").textContent = JSON.parse(e.data).now
				}
			}

			connect()
		</script>
	</body>
</html>
		`, host)
}

func (h *DemoHandler) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Now string `json:"now"`
	}
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		err := conn.WriteJSON(Response{time.Now().Local().Format("03:04:05")})
		if err != nil {
			log.Err(err).Msg("Websocket error; closing")
			return
		}
		time.Sleep(time.Second)
	}
}

func main() {
	addr := ":3000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	http.Handle("/", NewDemoHandler())
	panic(http.ListenAndServe(addr, nil))
}
