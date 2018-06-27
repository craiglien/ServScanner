package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"container/list"
	"reflect"
)

type Test struct {
	Linenum, Num         int
	Area, Desc, FQDN, IP string
	Ports                []int
	DateIntro, BizReason string
}

func LoadfromFile(fn string) {
	tsts := list.New()
	// read file
	f, err := os.Open(fn)

	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("Error:", err)
		return
	}

	r := csv.NewReader(f)
	r.Comma = ','
	lineNum := 0

	for {
		rec, err := r.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("rec:", rec)
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("---line num: %d\n", lineNum)
		fmt.Printf("---line num: %d %s %s\n", lineNum, len(rec), rec)
		fmt.Printf("---%s\n", rec[0])
		lineNum += 1
		tsts.PushBack(rec)

	}

	count := tsts.Len()
	fmt.Printf("records: %d\n", count)

	for e := tsts.Front(); e != nil; e = e.Next() {
		fmt.Println(e)
		fmt.Println(e.Value)
		a := e.Value
		fmt.Println("type ----")
		fmt.Println(reflect.ValueOf(e).Kind())
		fmt.Println(reflect.ValueOf(*e).Kind())
		//		fmt.Println( *e.Value )
		fmt.Println(reflect.ValueOf(a).Kind())
		fmt.Println("----")
		//		fmt.Println( a[0:3] )
		fmt.Println("----")

	}

}
