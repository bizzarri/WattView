//******************************************************
//* emissions - read the emissions information from the watttime.org API
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two arguments on command line:
//* -debug  displays debug information
//* -l location uses that location to get info
//* defaults to CAISO_ZP26
//* see watttime.org for an interactive map to figure out
//* your grid designation
//*
//*****************************************************

package main

import "fmt"
import "net/http"
import "io/ioutil"
//import "net/url"
import "os"
//import "strings"
import "encoding/json"
import "time"
import "flag"


type Wtoken struct {
     token string `json:"token"`

}

type Response struct {
	Barea string `json:"ba"`
	Valid int `json:"validFor"`
	Validuntil string `json:"validUntil"`
        Rating string `json:"rating"`
	Green string `json:"switch"`
	Percent string `json:"percent"`
	Freq string `json:"freq"`
	Market string `json:"market"`
}
	
func Check(val error, explain string)  {
	if val != nil {
		fmt.Printf("Error: %s\n%v\n",explain,val)
		os.Exit(-1)
	}
}

func main() {
	//
	// check for debug option
	// check for balancing authority parameter
	//
	// version number for program
	//
	version := 0.0
	//
	// loc - balancing authority
	// debug - debug flag
	//
	var loc string 
	boolPtr := flag.Bool("debug", false, "Debug flag")
	flag.StringVar(&loc, "l", "CAISO_ZP26", "Balancing Authority")
	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

     fmt.Printf("WattTime Emissions Real Time Analysis\n")
     timeout := time.Duration(5 * time.Second)
     client := &http.Client{
     	    Timeout: timeout,
	    }
     req,err := http.NewRequest("GET","https://api2.watttime.org/v2/login",nil)
     req.SetBasicAuth("bizzarri","Idontlike2018")
     resp, err := client.Do(req)
	Check (err,"Error login request call")

     defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if debug {
	   fmt.Printf("body: %s\n",bodyText)
	}
     var wtoken map[string]interface{}
     err = json.Unmarshal(bodyText,&wtoken)
	Check (err,"Error unmarshalling first call for token")
     if debug {
	fmt.Printf("token: %s\n",wtoken["token"])
	}
     request := "https://api2.watttime.org/v2/index/?ba="+loc+"&latitude=&longitude=&style=all"
     req,err = http.NewRequest("GET",request,nil)
	Check (err,"Error creating NewRequest")
     bearer := "Bearer " + wtoken["token"].(string)
     req.Header.Add("Authorization",bearer)
     resp, err = client.Do(req)
	Check(err,"Error getting NewRequest")

     response, err := ioutil.ReadAll(resp.Body)
	Check (err,"Error reading response")
     if debug {
	     fmt.Printf("Response: %s\n",response)
     }
//	var emisres map[string]interface{}
	var emisres Response
     err = json.Unmarshal(response,&emisres)
	Check (err,"Error unmarshalling response")
	fmt.Printf("\nReport for balancing authority: %s\n",emisres.Barea)
        if emisres.Green == "0" {
		fmt.Printf("Don't switch (not green)\n")
	} else {
		fmt.Printf("Switch! (green grid)\n")
	}
       
        timed, err := time.Parse(time.RFC3339,emisres.Validuntil)
	Check(err,"Error parsing Valid Until time")
	fmt.Printf("Valid Until: %02d:%02d:%02d UT\n",timed.Hour(),timed.Minute(),timed.Second())
	fmt.Printf("Rating (0=Extremely Clean, 5=Harmful): %s\n",emisres.Rating)
        fmt.Printf("Percent Dirty (0-100): %s\n",emisres.Percent)
}
