package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"ServScanner/common"
	"ServScanner/data"
	"ServScanner/nettests"
)

var (
	Version, buildTime string
)

func getVersion() string {
	return Version
}

func getBuildtime() string {
	return buildTime
}

func showInterfaces() {
	ifaces, err := net.Interfaces()

	// handle err
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
			return
		}

		for _, addr := range addrs {
			//			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				fmt.Printf("%v: %s (%s) ", i.Name, v, v.IP.DefaultMask())
				//				ip = v.IP
			case *net.IPAddr:
				fmt.Printf("%v : %s (%s) ", i.Name, v, v.IP.DefaultMask())
				//				ip = v.IP
			}
		}
	}

	fmt.Print("\n")
}

func showTests() {
	testNum := 1

	for t, v := range data.Tests {
		tstid := fmt.Sprintf("%d:%s-%d-%d", t, v.Area, v.Linenum, testNum)

		fmt.Printf("%sDESC:%s FQDN:%s IP:%s PORTS:%v\n",
			tstid, v.Desc, v.FQDN, v.IP, v.Ports)

		testNum += 1

	}
}

func showAreas() {
	m := make(map[string]bool)

	for _, v := range data.Tests {
		//		fmt.Printf("map: %v\n", m)
		m[v.Area] = true
	}

	//	fmt.Printf("map: %v\n", m)

	sk := make([]string, 0, len(m))

	for k := range m {
		sk = append(sk, k)
	}

	sort.Strings(sk)

	fmt.Printf("Areas: %v\n", strings.Join(sk, ", "))

}

func runTests(srcIP net.Addr, areaList *string, ipList *string) {
	DEBUG := 0
	testNum := 1

	if DEBUG > 0 {
		fmt.Println("Area List " + *areaList)
	}
	if DEBUG > 0 {
		fmt.Println("IP List " + *ipList)
	}
	lOa := strings.Split(*areaList, ":")
	lOi := strings.Split(*ipList, ":")
	if DEBUG > 0 {
		fmt.Println("IP List " + lOi[0])
	}

	for t, v := range data.Tests {

		doTest := false

		aMatch, err := regexp.MatchString("all", lOa[0])

		if err != nil {
			fmt.Printf("error:%s\n", err)
		}

		//		iMatch, err := regexp.MatchString(v.IP, lOi[0] )

		if err != nil {
			fmt.Printf("error:%s\n", err)
		}

		if aMatch {
			doTest = true
		} else {
			for kee, _ := range lOa {
				mat, err := regexp.MatchString("(?i)"+lOa[kee], v.Area)
				if err != nil {
					fmt.Printf("error:%s\n", err)
				}
				if mat {
					doTest = true
				}
			}
		}

		if doTest {
			// Check on IP flag

		}

		if doTest {
			tstid := fmt.Sprintf("%d:%s-%d-%d", t, v.Area, v.Linenum, testNum)

			c := make(chan string)

			go runTest(srcIP, tstid, v, c)
			msg := <-c
			fmt.Print(msg)

			testNum += 1
		}

	}
}

func runTest(srcIP net.Addr, id string, t common.Test, c chan string) {
	DEBUG := 0
	var iplist []string

	out := fmt.Sprintf("%sDESC:%s FQDN:%s IP:%s PORTS:%v", id, t.Desc, t.FQDN, t.IP, t.Ports)
	out += "\n"

	if t.FQDN != "" {
		out += fmt.Sprintf("%sDNSLookup", id)
		val, err := nettests.DNSLookup(srcIP, t.FQDN)
		if err != nil {
			out += fmt.Sprintf("err(%s)", err.Error()) //, val[0])
			if t.IP != "" {
				iplist = []string{t.IP}
			}
		} else {
			out += fmt.Sprintf("%s", val)
			iplist = append(iplist, val...)
		}
		out += fmt.Sprintf("\n")
	} else {
		iplist = append(iplist, t.IP)
	}

	if DEBUG > 0 {
		fmt.Printf("Testing list of ips: %v\n", iplist)
	}

	for ip := range iplist {
		match, err := regexp.MatchString("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}", iplist[ip])
		if match {
			if len(t.Ports) > 0 {
				for port := range t.Ports {
					out += fmt.Sprintf("%sSocketConnect", id)

					if DEBUG > 0 {
						fmt.Printf("Testing ip: %s port: %v\n",
							iplist[ip], t.Ports[port])
					}

					val, err := nettests.SocketConnect(srcIP, iplist[ip], t.Ports[port])

					if err != nil {
						out += fmt.Sprintf("err%s", val)
					} else {
						out += fmt.Sprintf("%s", val)
					}
					out += "\n"
				}
			} else {
				out += fmt.Sprintf("%sSocketConnect", id)
				out += "0 length port list"
				out += "\n"

			}
		} else {
			if DEBUG > 0 {
				out += fmt.Sprintf("Skipping testing ip: %s error: %v\n", iplist[ip], err)
			}
		}
	}

	c <- out
}

func main() {
	var srcIPaddr net.Addr

	shTests := flag.Bool("v", false, "show tests")

	shVersion := flag.Bool("V", false, "show Version")

	shAreas := flag.Bool("l", false, "show areas")

	areaList := flag.String("a", "all", "What areas to test")

	ipList := flag.String("i", "*", "What IPs to test")

	srcIP := flag.String("S", "", "What IPs to test")

	flag.Parse()

	if *shAreas {
		showAreas()
		os.Exit(0)
	}

	if *shTests {
		showTests()
		os.Exit(0)
	}

	if *shVersion {
		fmt.Printf("Version:%s Build time:%s\n", getVersion(), getBuildtime())
		os.Exit(0)
	}

	if len(*srcIP) > 0 {
		ip := net.ParseIP(*srcIP)
		if ip == nil {
			fmt.Println("Error in IP address ", *srcIP)
			os.Exit(1)
		}
		srcIPaddr = &net.IPAddr{ip, ""}
	}

	defer func(st time.Time) {
		fmt.Printf("End time: %s runtime: %s\n",
			time.Now().Format(time.RFC3339), time.Since(st))
		fmt.Println("+----------------------------+")
		fmt.Println("| COPY EVERYTHING ABOVE THIS |")
		fmt.Println("+----------------------------+")
	}(time.Now())

	fmt.Println("+----------------------------+")
	fmt.Println("| COPY EVERYTHING BELOW THIS |")
	fmt.Println("+----------------------------+")

	fmt.Printf("Start time: %s %s %s\n",
		time.Now().Format(time.RFC3339), getVersion(), buildTime)
	fmt.Printf("IP Interfaces:")
	showInterfaces()

	runTests(srcIPaddr, areaList, ipList)
}
