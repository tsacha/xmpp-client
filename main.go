package main

import (
	"encoding/json"
	"framagit.org/tsacha-xmpp/xmpp"
	"io/ioutil"
	"os/user"
)

type Configuration struct {
	Account  string
	Password string
}

func main() {

	usr, _ := user.Current()

	file, _ := ioutil.ReadFile(
		usr.HomeDir + "/.config/xmpp-client/config.json")

	var c Configuration
	json.Unmarshal(file, &c)

	xmpp.Connect(c.Account, c.Password)

	select {}
}
