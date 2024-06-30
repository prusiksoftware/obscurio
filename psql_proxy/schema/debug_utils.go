package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func writetoFile(str string, filename string) {
	os.Remove(filename)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	n, err := f.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
	if n != len(str) {
		log.Fatal("failed to write all bytes")
	}
}

func readFromFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b := make([]byte, 1024)
	s := ""
	more := true
	for more {
		n, err := f.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		s += string(b[:n])
		if n < 1024 {
			more = false
		}
	}
	return s
}

func PrettyPrint(i interface{}) {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func PrettyPrints(i interface{}) string {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
