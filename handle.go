package main

import (
	"encoding/pem"
	"crypto/x509"
	"crypto/md5"
	"crypto/rsa"
	"crypto"
	"encoding/hex"
	"strconv"
	"net/http"
	"log"
	"crypto/rand"
)

type ticketMux struct {
	mux *http.ServeMux
}

func (p ticketMux) preHandle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
func (p ticketMux) RegisterRouter() {
	p.mux = http.NewServeMux()
	http.Handle("/", p.preHandle(p.mux))
	p.mux.HandleFunc("/index", p.index)
	p.mux.HandleFunc("/rpc/ping.action", p.pingAction)
	p.mux.HandleFunc("/rpc/obtainTicket.action", p.obtainTicket)
}

func (p ticketMux) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("JetBrains license server is running!"))
}

func (p ticketMux) pingAction(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	salt := r.URL.Query().Get("salt")
	xmlResponse := "<PingResponse><message></message><responseCode>OK</responseCode><salt>" +
		salt + "</salt></PingResponse>"
	xmlSignature, _ := p.signature(xmlResponse)
	w.Header().Add("Content-Type", "text/xml")
	w.Write([]byte("<!-- " + xmlSignature + " -->\n" + xmlResponse))
}

func (p ticketMux) obtainTicket(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	salt := r.URL.Query().Get("salt")
	username := r.URL.Query().Get("userName")

	if salt == "" || username == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	prolongationPeriod := 607875500

	xmlResponse := "<ObtainTicketResponse><message></message><prolongationPeriod>" +
		strconv.Itoa(prolongationPeriod) + "</prolongationPeriod><responseCode>OK</responseCode><salt>" +
		salt + "</salt><ticketId>1</ticketId><ticketProperties>licensee=" + username +
		"\tlicenseType=0\t</ticketProperties></ObtainTicketResponse>"
	xmlSignature, _ := p.signature(xmlResponse)
	w.Header().Add("Content-Type", "text/xml")
	w.Write([]byte("<!-- " + xmlSignature + " -->\n" + xmlResponse))
}

func (p ticketMux) signature(message string) (string, error) {
	pemData, _ := pem.Decode(gPrivateKey)
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(pemData.Bytes)
	hashedMessage := md5.Sum([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.MD5, hashedMessage[:])
	if err != nil {
		return "", err
	}

	hexSignature := hex.EncodeToString(signature)
	return hexSignature, nil
}
