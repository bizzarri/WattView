package main

import "fmt"
import "net/http"
import "io/ioutil"
import "net/url"
import "os"

func main() {
     fmt.Printf("WattTime\n")
     resp, err := http.PostForm("https://api2.watttime.org/v2/register",
     url.Values{"username":{"bizzarri"},"password":{"Idontlike2018"},
     "email":{"maurice@bizzarrisoftware.com"},"org":{"Bizzarri Software"}})
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     defer resp.Body.Close()
     body, err := ioutil.ReadAll(resp.Body)
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }

     fmt.Printf("body: %s\n",body)
}