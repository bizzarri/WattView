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
	Val float64 `json:"value"`
	Fuel string `json:"fuel"`
}
	

func Check(err error, msg string)  {
     if err != nil {
     fmt.Printf("Error: %s\n%v\n",msg,err)
     os.Exit(-1)
     }
}
func main() {
     debug := false
     location := "CAISO_ZP26"
     fmt.Printf("Grid Data for WattTime Zone %s\n",location)
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

	gridstr := "https://api2.watttime.org/v2/data/?ba="
	gridstr = gridstr + location
	gridstr = gridstr + "&latitude=&longitude=&starttime=2019-01-05T09:00:00-00:00&endtime=2019-01-05T09:05:00-00:00"
	req,err = http.NewRequest("GET",gridstr,nil)

     bearer := "Bearer " + wtoken["token"].(string)
     req.Header.Add("Authorization",bearer)
     resp, err = client.Do(req)
	Check (err,"Error retrieving data")
     response, err := ioutil.ReadAll(resp.Body)
	Check (err,"Error reading data from GET")
     if debug {
	     fmt.Printf("Response: %s\n",response)
     }

        var unwrap []Datadef
	err = json.Unmarshal(response,&unwrap)
        Check(err,"Error unmarshalling response")
//
//        for idx := 0; idx < 5; idx++ {
        for idx := range unwrap {
	fmt.Printf("Zone: %s\n",unwrap[idx].Ba)
	fmt.Printf("Type: %s\n",unwrap[idx].Dtype)
	fmt.Printf("Point Time: %s\n",unwrap[idx].Point_time)
	fmt.Printf("Frequency: %f\n",unwrap[idx].Frequency)
	fmt.Printf("Value: %f\n",unwrap[idx].Val)
	fmt.Printf("Market: %s\n",unwrap[idx].Market)
	fmt.Printf("Fuel: %s\n",unwrap[idx].Fuel)
	}		

		

}
