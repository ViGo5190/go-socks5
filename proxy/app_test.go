package proxy

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxyOK(t *testing.T) {

	hw := []byte{72, 101, 108, 108, 111, 44, 32, 99, 108, 105, 101, 110, 116, 10} // Hello, client

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(hw)
	}))

	defer ts.Close()

	l, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		t.Errorf("Error on creating listener: %v", err)
	}

	lg := log.New()
	s := &Server{Logger: lg}
	go s.Serve(l)
	defer s.Stop()

	proxyURL, err := url.Parse("socks5://" + l.Addr().String())
	if err != nil {
		t.Errorf("Error on parse proxy url: %v", err)
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	r, err := myClient.Get(ts.URL)
	if err != nil {
		t.Errorf("Error while get response from server: %v", err)
	}

	defer r.Body.Close()

	contents, err := ioutil.ReadAll(r.Body)
	//_, err = ioutil.ReadAll(r.Body)

	if err != nil {
		t.Errorf("Error on read response body : %v", err)
	}

	if !bytes.Equal(contents, hw) {
		t.Errorf("Expected response body %v, got %v", hw, contents)
	}

}

func BenchmarkProxy_Start(b *testing.B) {
	hw := []byte{72, 101, 108, 108, 111, 44, 32, 99, 108, 105, 101, 110, 116, 10} // Hello, client

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(hw)
	}))

	defer ts.Close()

	l, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		b.Errorf("Error on creating listener: %v", err)
	}

	lg := log.New()
	s := &Server{Logger: lg}
	go s.Serve(l)

	proxyURL, err := url.Parse("socks5://" + l.Addr().String())
	if err != nil {
		b.Errorf("Error on parse proxy url: %v", err)
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	for n := 0; n < b.N; n++ {
		r, err := myClient.Get(ts.URL)
		if err != nil {
			b.Errorf("Error while get response from server: %v", err)
		}

		contents, err := ioutil.ReadAll(r.Body)

		if err != nil {
			r.Body.Close()
			b.Errorf("Error on read response body : %v", err)
		}

		if !bytes.Equal(contents, hw) {
			r.Body.Close()
			b.Errorf("Expected response body %v, got %v", hw, contents)
		}
		r.Body.Close()
	}

}
