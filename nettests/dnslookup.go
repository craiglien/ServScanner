package nettests

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

func DNSLookup(srcIP net.Addr, host string) ([]string, error) {
	DEBUG := 0

	if DEBUG > 0 {
		fmt.Fprintf(os.Stderr, "DNSLookup>>>%s<<<\n", host)
	}

	match, err := regexp.MatchString("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}", host)

	if match {
		// Ok IPv4
		if DEBUG > 0 {
			fmt.Fprintf(os.Stderr, ">>>already an IP (%s)\n", host)
		}
		return []string{host}, nil
	}

	info, err := net.LookupHost(host)

	if DEBUG > 0 {
		fmt.Fprintf(os.Stderr, "HostLookupReply>>>%v<<>%v<<<\n", info, err)
	}

	if err != nil {
		return info, err
	} else {
		var ret []string
		for hn, hc := range info {
			if DEBUG > 0 {
				fmt.Fprintf(os.Stderr, "HostLookupReply>>>%v<<>>%v<<\n", hn, hc)
			}
			match, _ := regexp.MatchString("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}", hc)
			if match {
				ret = append(ret, hc)
			}

		}

		return ret, err
	}
}
