package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var DEBUG = 0

func LoadfromFile(fn string) {
	// read file
	f, err := os.Open(fn)

	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("Error:", err)
		return
	}

	scanner := bufio.NewScanner(f)

	lineNum := 0
	area := ""

	comment := regexp.MustCompile(`\s*^[#!]`)

	for scanner.Scan() {
		ln := scanner.Text()
		lineNum += 1

		if comment.MatchString(ln) {
			if DEBUG > 0 {
				fmt.Fprintln(os.Stderr, "comment line", ln)
			}
			continue
		}

		if DEBUG > 0 {
			fmt.Fprintln(os.Stderr, ">>>", lineNum, ln)
		}

		elm := strings.Split(ln, ",")

		if len(elm) == 1 {
			// Area
			area = elm[0]
			continue
		}
		if len(elm) >= 4 {
			out := " common.Test{ " +
				strings.Join([]string{
					fmt.Sprintf("Linenum: %d", lineNum),
					fmt.Sprintf("Area: \x22%s\x22", area),
				}, ", ") + ", "

			out += fmt.Sprintf("Desc: \x22%s\x22", elm[0])
			out += ",\n          "
			out += fmt.Sprintf("FQDN: \x22%s\x22, ", elm[1])
			out += fmt.Sprintf("IP: \x22%s\x22, ", elm[2])

			ports := strings.Split(elm[3], ":")

			if DEBUG > 0 {
				fmt.Fprintln(os.Stderr, ">>>", ports)
			}

			out += fmt.Sprintf("Ports: []int{ %s },", strings.Join(ports, ", "))

			if len(elm) >= 5 {
				out += fmt.Sprintf("DateIntro: \x22%s\x22, ", elm[4])
			}

			if len(elm) >= 6 {
				out += fmt.Sprintf("BizReason: \x22%s\x22, ", elm[5])
			}

			out += " },\n"

			if DEBUG > 0 {
				fmt.Fprintln(os.Stderr, ">>>", out)
			}

			fmt.Printf(out)

			continue
		}
		fmt.Printf(">Bad format>>%q\n", strings.Split(ln, ","))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

}

func main() {
	fmt.Println(`
package data
import "ServScanner/common"
var Tests = []common.Test{`)

	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		fmt.Println("// file:", v)
		LoadfromFile(v)
	}

	fmt.Println("}\n")

}
