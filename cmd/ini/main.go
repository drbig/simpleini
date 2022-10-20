package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"bytes"
	"strings"

	"github.com/pedroalbanese/simpleini"
)

var (
	p = flag.String("p", "", "Parameter (imperative, omit to list parameters)")
	s = flag.String("s", "", "Section (imperative, omit to list sections)")
	v = flag.String("v", "", "Value ('-' or '?' for delete entry)")
	f = flag.String("f", "", "Target INI File ('-' for stdin/stdout)")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage of "+os.Args[0]+": ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var file *os.File
	var err error
	if *f == "-" || *f == "" {
		file = os.Stdin
	} else if *f != "-" {
		file, err = os.OpenFile(*f, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *s == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		str := ini.Sections()
		for i := 0; i < len(str); i++ {
			fmt.Printf("%s\n", str[i])
		}
		os.Exit(0)
	}
	if *p == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		str, err := ini.Properties(*s)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(str); i++ {
			fmt.Printf("%s\n", str[i])
		}
		os.Exit(0)
	} 
	if *v == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		str, err := ini.GetString(*s, *p)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", str)
		os.Exit(0)
	} 
	if *v == "?" || *v == "-" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		ini.SetString(*s, *p, "delete")
		var buf bytes.Buffer
		err = ini.Write(&buf, true)
		if err != nil {
			log.Fatal("Write error: %s", err)
			return
		}
		lines := strings.Replace(buf.String(), *p+" = delete\r\n", "", 1)
		lines = strings.Replace(lines, *p+" = delete\n", "", 1)
		if *f == "-" || *f == "" {
			fmt.Print(lines)
			os.Exit(0)
		} else {
			f, err := os.OpenFile(*f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			_, err = f.WriteString(lines)
			if err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	} 
	ini, err := simpleini.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	ini.SetString(*s, *p, *v)
	val, err := ini.GetString(*s, *p)
	if err != nil {
		log.Fatal(err)
	}
	if val != *v {
		log.Fatal("Bad posterior value")
		return
	}
	var buf bytes.Buffer
	err = ini.Write(&buf, true)
	if err != nil {
		log.Fatal("Write error: %s", err)
		return
	}
	_, err = file.Seek(0, 0)
	if *f == "-" || *f == "" {
		ini.Write(os.Stdout, true)
	} else {
		ini.Write(file, true)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
