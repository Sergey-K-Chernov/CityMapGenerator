package main

import (
	"testing"
	"strconv"
	"strings"
	"slices"
)


func TestJsonToCookiesStrings(t *testing.T){
	str1s := "{\"Value1\":250,\"Array\":[1,2,3]}"
	str1r := "{%22Value1%22:250,%22Array%22:[1,2,3]}"

	json1 := []byte(str1s)
	r1 := jsonToCookieStrings(json1)
	if len(r1) != 1 {
		t.Fatalf("Wrong results number in test 1")
	}
	if r1[0] != str1r {
		t.Fatalf("Wrong result in test 1.\nNeed:\n%s\nGot:\n%s\n", str1r, r1[0])
	}

	str2s := ""
	for i := range 9000 {
		if i == 4094 {
			str2s += "r"
		} else {
			str2s += strconv.Itoa(i%10)
		}
	}
	str2s = strings.ReplaceAll(str2s, "r", "🗵")

	json2 := []byte(str2s)
	r2 := jsonToCookieStrings(json2)
	if len(r2) != 3 {
		t.Fatalf("Wrong results number in test 2. Need 3, got %d", len(r2))
	}

	r := r2[0] + r2[1] + r2[2]
	if r != str2s {
		t.Fatalf("Wrong result in test 2.\nNeed:\n%s\nGot:\n%s\n", str1r, r1[0])
	}

}


func TestCookiesStringsToJson(t *testing.T){
	str1r := "{\"Value1\":250,\"Array\":[1,2,3]}"
	str1s := []string{"{%22Value1%22:250,%22Ar", "ray%22:[1,2,3]}"}

	need := []byte(str1r)
	r := cookieStringsToJson(str1s)
	if !slices.Equal(r, need) {
		t.Fatalf("Wrong result")
	}
}