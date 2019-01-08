package main

import "fmt"
import "net/http"
import "io/ioutil"
//import "net/url"
import "os"
//import "strings"
import "encoding/json"
import "time"

type Wtoken struct {
     token string `json:"token"`

     }
func main() {
     fmt.Printf("WattTime\n")
     timeout := time.Duration(5 * time.Second)
     client := &http.Client{
     	    Timeout: timeout,
	    }
     req,err := http.NewRequest("GET","https://api2.watttime.org/v2/login",nil)
     req.SetBasicAuth("bizzarri","Idontlike2018")
     resp, err := client.Do(req)

     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     defer resp.Body.Close()
     bodyText, err := ioutil.ReadAll(resp.Body)
//     fmt.Printf("body: %s\n",bodyText)
     var wtoken map[string]interface{}
     err = json.Unmarshal(bodyText,&wtoken)
     if err != nil {
     fmt.Printf("Error: %v\n",err)
     os.Exit(-1)
     }
     	
//     fmt.Printf("token: %s\n",wtoken["token"])

     req,err = http.NewRequest("GET","https://api2.watttime.org/v2/index/?ba=CAISO_ZP26&latitude=&longitude=&style=all",nil)

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