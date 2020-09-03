package api

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
)

type network struct {
	ID       int    `json:"id"`
	RemoteIP string `json:"remoteip"`
	Hostname string `json:"hostname"`
	LocalIP  string `json:"localip"`
	PublicIP string `json:"publicip"`
	Nmap     []int  `json:"nmap"`
	OS       string `json:"os"`
}

var countID int
var finalJSON = make(map[string]interface{})
var datas []network

// BasicAuth - basic auth for POST and GET endpoints
func BasicAuth(logger *log.Logger, username, password, realm string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var remoteIP string
		forwarded := r.Header.Get("X-FORWARDED-FOR")
		if forwarded != "" {
			logger.Printf("from 'BasicAuth' - someone IP: %s \n", forwarded)
			remoteIP = forwarded
		}
		logger.Printf("from 'BasicAuth' - someone IP: %s \n", r.RemoteAddr)
		remoteIP = r.RemoteAddr
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			logger.Printf("Unauthorised.\n")
			w.Write([]byte("Doesn't work like that.\n"))
			return
		}
		switch {
		case r.Method == http.MethodGet:
			logger.Printf("MethodGet")
			getIPs(w, r)
		case r.Method == http.MethodPost:
			logger.Printf("MethodPost")
			postIPs(w, r, remoteIP)
		}
	}
}

func postIPs(w http.ResponseWriter, r *http.Request, remoteIP string) {
	var net network
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&net)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	countID++
	datas = append(datas, network{
		ID:       countID,
		RemoteIP: remoteIP,
		Hostname: net.Hostname,
		LocalIP:  net.LocalIP,
		PublicIP: net.PublicIP,
		Nmap:     net.Nmap,
		OS:       net.OS,
	})
}

func getIPs(w http.ResponseWriter, r *http.Request) {
	finalJSON["data"] = datas
	// json.NewEncoder(w).Encode(finalJson)
	b, _ := json.Marshal(finalJSON)
	w.Write([]byte(b))
}
