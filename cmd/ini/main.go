package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pedroalbanese/simpleini"
)

var (
	f = flag.String("f", "", "Target INI File ('-' for stdin/stdout)")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s set|get|del [section] [parameter] [value]\n", os.Args[0])
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

	if flag.Arg(0) == "del" && flag.Arg(1) != "" && flag.Arg(2) == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		simpleini.DeleteSection(ini, flag.Arg(1))

		if *f == "-" || *f == "" {
			fmt.Print(ini)
			os.Exit(0)
		} else {
			f, err := os.OpenFile(*f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = ini.Write(f, true)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		os.Exit(0)
	}

	if flag.Arg(0) == "del" && flag.Arg(1) != "" && flag.Arg(2) != "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		simpleini.DeleteProperty(ini, flag.Arg(1), flag.Arg(2))

		if *f == "-" || *f == "" {
			fmt.Print(ini)
			os.Exit(0)
		} else {
			f, err := os.OpenFile(*f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = ini.Write(f, true)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		os.Exit(0)
	}

	if flag.Arg(0) == "get" && flag.Arg(1) == "" {
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

	if flag.Arg(0) == "get" && flag.Arg(2) == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		str, err := ini.Properties(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(str); i++ {
			fmt.Printf("%s\n", str[i])
		}
		os.Exit(0)
	}

	if flag.Arg(0) == "get" && flag.Arg(3) == "" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		str, err := ini.GetString(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", str)
		os.Exit(0)
	}

	if flag.Arg(0) == "set" {
		ini, err := simpleini.Parse(file)
		if err != nil {
			log.Fatal(err)
		}
		ini.SetString(flag.Arg(1), flag.Arg(2), flag.Arg(3))
		val, err := ini.GetString(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		if val != flag.Arg(3) {
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
}
