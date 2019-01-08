package main

import "fmt"
import "net/http"
import "io/ioutil"
//import "net/url"
import "os"
//import "strings"
import "encoding/json"

type Wtoken struct {
     token string `json:"token"`

     }
func main() {
     fmt.Printf("WattTime\n")
     client := &http.Client{}
     req,err := http.NewRequest("GET","https://api2.watttime.org/v2/login",nil)
     req.SetBasicAuth("bizzarri","Idontlike2018")
     resp, err := client.Do(req)
     defer resp.Body.Close()
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     bodyText, err := ioutil.ReadAll(resp.Body)
//     fmt.Printf("body: %s\n",bodyText)
     var wtoken map[string]interface{}
     err = json.Unmarshal(bodyText,&wtoken)
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     	
//     fmt.Printf("token: %s\n",wtoken["token"])

     req,err = http.NewRequest("GET","https://api2.watttime.org/v2/data/?ba=CAISO_ZP26&latitude=&longitude=&starttime=2018-12-31&endtime=2019-01-05",nil)
     //req,err = http.NewRequest("GET","https://api2.watttime.org/v2/ba-from-loc?latitude=34.57&logitude=-121.10",nil)
//     req,err = http.NewRequest("GET","https://api2.watttime.org/v2/index?ba=CAISO_ZP26",nil)
     bearer := "Bearer " + wtoken["token"].(string)
     req.Header.Add("Authorization",bearer)
     resp, err = client.Do(req)
//     defer resp.Body.Close()
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     response, err := ioutil.ReadAll(resp.Body)
     fmt.Printf("Response: %s\n",response)
}