package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
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

func getResponseWithBasicAuth(m map[string]string, user string, passwd string) (*http.Response, error) {
	pt, _, line, _ := runtime.Caller(0) // debug用に現在のスタックから情報を取得

	req, err := newRequest(m)

	proxyUrl, _ := url.Parse(os.Getenv("HTTP_PROXY"))
	proxy := "no"
	if proxyURL, _ := http.ProxyFromEnvironment(req); proxyURL != nil {
		proxy = proxyURL.String()
	}
	log.Printf("%s, %d: [DEBUG] request proxy: %s\n", runtime.FuncForPC(pt).Name(), line, proxy)


	var tr *http.Transport
	if proxy == "no" {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyUrl),
		}
	}
	client := &http.Client{
		Transport: tr,
	}

	req.SetBasicAuth(user, passwd)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%d] Unable to get this url", resp.StatusCode)
	}

	return resp, err
}

func getResponse(m map[string]string) (*http.Response, error) {
	req, err := newRequest(m)

	//proxyUrl, _ := url.Parse(os.Getenv("HTTP_PROXY"))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: tr,
	}

	proxy := "no"
	if proxyURL, _ := http.ProxyFromEnvironment(req); proxyURL != nil {
		proxy = proxyURL.String()
	}
	log.Printf("[DEBUG] request proxy: %s\n", proxy)

	//resp, err := http.DefaultClient.Do(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%d] Unable to get this url", resp.StatusCode)
	}

	return resp, err
}
