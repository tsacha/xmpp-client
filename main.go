package main

import (
	//	"github.com/skratchdot/open-golang/open"

	"encoding/json"
	libxmpp "framagit.org/tsacha-xmpp/xmpp"
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"io/ioutil"
	"net/http"
	"os/user"

	"fmt"
)

type Configuration struct {
	Account  string
	Password string
}

type Connection struct {
	Jid      string `json:"jid"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
	Resource string `json:"resource"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetCredentials() (string, string) {
	usr, _ := user.Current()

	file, _ := ioutil.ReadFile(
		usr.HomeDir + "/.config/xmpp-client/config.json")

	var c Configuration
	json.Unmarshal(file, &c)

	return c.Account, c.Password
}

func testRoster(account string, password string, domain string, resource string) []byte {

	xmpp := libxmpp.Connect(account, password, domain, resource)

	go xmpp.InfinitePing()

	//	xmpp.Disco(domain)
	//	xmpp.Disco(account)

	xmpp.GetRoster()
	contacts, _ := json.Marshal(xmpp.State.Roster.Contacts)

	xmpp.Close()

	return contacts
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Client subscribed")
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		data := Connection{}
		json.Unmarshal(msg, &data)

		fmt.Println(data.Jid)

		conn.WriteMessage(mt, testRoster(data.Jid, data.Password, data.Domain, data.Resource))
	}
}

func main() {
	fmt.Println("Hello world")

	// mux handler
	router := mux.NewRouter()

	router.HandleFunc("/websocket", WsHandler)

	// Serve static assets via the "static" directory
	fs := rice.MustFindBox("static").HTTPBox()

	staticFileServer := http.StripPrefix("/", http.FileServer(fs))
	router.Handle("/{path:.*}", staticFileServer)

	// Open web browser
	//	go open.Start("http://localhost:5282")

	// Serve this program forever
	go http.ListenAndServe(":5282", router)

	select {}
}
