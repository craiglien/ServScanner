package nettests

import (
	"fmt"
	"net"
	"testing"
)

func TestDNSLookup(t *testing.T) {
	var testCases = []struct {
		host     string
		srcIP       string
		expectOk bool
	}{
		// no reply
		{"127.0.0.1", "127.0.0.1", true},
		// many A records
		{"www.google.com", "127.0.0.1", true},
		// nothing returned
		{"com.caife9ieBiequeotheri", "127.0.0.1", false},
	}

	for tn, tc := range testCases {
		DEBUG := 0
		ip := net.ParseIP(tc.srcIP)
		
		t.Logf("#%d%v\n", tn, tc)
		str, err := DNSLookup(&net.IPAddr{ip,""}, tc.host)
		if DEBUG > 0 {
			fmt.Printf(">>>>%v<>%v<<<<\n", str, err)
		}

		t.Logf("#%d%v(%v)\n", tn, str, err)

		if tc.expectOk && err != nil {
			t.Errorf("#%d: got %v; want nil", tn, err)
		} else if !tc.expectOk && err == nil {
			t.Errorf("#%d: got nil; want non-nil", tn)
		}
	}
}
