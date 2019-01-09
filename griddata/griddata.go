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

type Datadef struct {
	Ba string `json:"ba"`
	Dtype string `json:"datatype"`
	Point_time string `json:"point_time"`
	Frequency float64 `json:"frequency"`
	Market string `json:"market"`
	Val float64 `json"value"`
}
	

func Check(err error, msg string)  {
     if err != nil {
     fmt.Printf("Error: %s\n%v\n",msg,err)
     os.Exit(-1)
     }
}
func main() {
     debug := false
     fmt.Printf("WattTime\n")
     client := &http.Client{}
     req,err := http.NewRequest("GET","https://api2.watttime.org/v2/login",nil)
     req.SetBasicAuth("bizzarri","Idontlike2018")
     resp, err := client.Do(req)
	Check(err,"Error WattTime login request")
     defer resp.Body.Close()
     bodyText, err := ioutil.ReadAll(resp.Body)
	Check(err,"Error reading body")
     if debug {
	fmt.Printf("body: %s\n",bodyText)
     }
     var wtoken map[string]interface{}
     err = json.Unmarshal(bodyText,&wtoken)
	Check(err,"Error unmarshalling body text from login")
     if debug {
	     fmt.Printf("token: %s\n",wtoken["token"])
     }

     req,err = http.NewRequest("GET","https://api2.watttime.org/v2/data/?ba=CAISO_ZP26&latitude=&longitude=&starttime=2019-01-05T09:00:00-00:00&endtime=2019-01-05T09:05:00-00:00",nil)

     bearer := "Bearer " + wtoken["token"].(string)
     req.Header.Add("Authorization",bearer)
     resp, err = client.Do(req)
	Check (err,"Error retrieving data")
     response, err := ioutil.ReadAll(resp.Body)
	Check (err,"Error reading data from GET")
     if debug {
	     fmt.Printf("Response: %s\n",response)
     }


	var unwrap []interface{}
	err = json.Unmarshal(response,&unwrap)
        Check(err,"Error unmarshalling response")
//
	fmt.Printf("data: %s\n",unwrap[0])
        var datawrap Datadef
//

	datawrap = unwrap[0].(Datadef)
	fmt.Printf("point_time: %s\n",datawrap.Point_time)
		

}
