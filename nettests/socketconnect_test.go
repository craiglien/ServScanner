package nettests

import (
	"fmt"
	"net"
	"testing"
)

func TestSocketConnect(t *testing.T) {
	var testCases = []struct {
		host     string
		port     int
		srcIP    string
		ret      string
		expectOk bool
	}{
		// connect to ssh port with echo
		{"127.0.0.1", 22, "127.0.0.1", "", true},
		// connect to http port without echo
		{"127.0.0.1", 80, "127.0.0.1", "", true},
		// connection refused
		{"127.0.0.1", 23, "127.0.0.1", "", false},
		// connection timedout
		{"1.1.1.1", 22, "127.0.0.1", "", false},
	}

	for tn, tc := range testCases {
		DEBUG := 0
		ip := net.ParseIP(tc.srcIP)

		t.Logf("#%d%v\n", tn, tc)
		str, err := SocketConnect(&net.IPAddr{ip,""},tc.host, tc.port)
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
