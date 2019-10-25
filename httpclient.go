package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"crypto/tls"
)

func newRequest(m map[string]string) (*http.Request, error) {
	url := m["url"]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// ここのRedmine依存処理はいずれ共通化する
	req.Header.Set("Content-Type", "application/json")
	if m["headerKey"] == "X-Redmine-API-Key" {
		req.Header.Set(m["headerKey"], m["headerValue"])
	}

	dump, err := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)
	if err != nil {
		log.Fatal("Error request dump")
	}

	return req, err
}

func getResponse(m map[string]string) (*http.Response, error) {
	req, err := newRequest(m)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%d] Unable to get this url", resp.StatusCode)
	}

	/*
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		} else if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[%d] Unable to get this url", resp.StatusCode)
		}
	*/

	return resp, err
}
