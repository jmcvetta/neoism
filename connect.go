// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

// +build !appengine

package neoism

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"crypto/x509"
	"errors"
	"github.com/jmcvetta/napping"
	"io/ioutil"
	"log"
	"os"
)

// Connect setups parameters for the Neo4j server
// and calls ConnectWithRetry()
func Connect(uri string) (*Database, error) {
	h := http.Header{}
	h.Add("User-Agent", "neoism")

	tcc := &tls.Config{}

	caCertFile := os.Getenv("CACERTSFILE")
	if caCertFile != "" {
		log.Println("CA Cert file was specified, appending it into the Cert pool...")
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, errors.New("Unable to append the certificate to the CA certs pool.")
		}

		tcc.RootCAs = caCertPool
	} else {
		tcc.InsecureSkipVerify = true
	}

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tcc}}

	db := &Database{
		Session: &napping.Session{
			Client: client,
			Header: &h,
		},
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if parsedUrl.User != nil {
		db.Session.Userinfo = parsedUrl.User
	}
	return connectWithRetry(db, parsedUrl, 0)
}
