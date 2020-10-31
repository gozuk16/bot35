package main

import (
	"fmt"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
)

func TestRedmineSuccess(t *testing.T) {
	config.Redmine.APIToken = "aaaa"
	result, err := redmine("https://my.redmine.jp/demo/issues/23182")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if result == "" {
		t.Fatal("failed test")
	}
}

func TestRedmineFail(t *testing.T) {
	config.Redmine.APIToken = "bbbb"
	_, err := redmine("https://my.redmine.jp/demo/issues/1")
	if err == nil {
		t.Fatal("failed test")
	}
	fmt.Println(err)
}
