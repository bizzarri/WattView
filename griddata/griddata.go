//******************************************************
//* griddata - read detailed grid data on a balancing authority
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
import "flag"
import "os"
//import "strings"
import "encoding/json"

type Wtoken struct {
     token string `json:"token"`

     }
	
type MakeAcct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Org      string `json:"org"`
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

	//
	// location - balancing authority for parameter
	// debug - debug flag
	//
        version := 0.0
	var location string 
	boolPtr := flag.Bool("debug", false, "Debug flag")
	flag.StringVar(&location, "l", "CAISO_ZP26", "Balancing Authority abreviation")
	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

	//
	// get account and password from $HOME/.WattTime/account
	// should be set by makeacct
	//
	homedir := os.Getenv("HOME")
	acctfile := homedir + "/.WattTime/account"
	accts, err := ioutil.ReadFile(acctfile)
	Check(err,"Accounts file not found or other read error")
	var macct MakeAcct
	err = json.Unmarshal(accts,&macct)
	Check(err,"Error unmarshalling accounts files")
	account := macct.Username
	password := macct.Password
        if debug {
		fmt.Printf("Account Name: %s\n",account)
		fmt.Printf("Password: %s\n",password)
	}


     fmt.Printf("Grid Data for Balancing Authority  %s\n",location)
     client := &http.Client{}
     req,err := http.NewRequest("GET","https://api2.watttime.org/v2/login",nil)
	req.SetBasicAuth(account,password)
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
	Check(err,"Error getting request")
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
	fmt.Printf("Balancing Authority: %s\n",unwrap[idx].Ba)
	fmt.Printf("Data Type: %s\n",unwrap[idx].Dtype)
	fmt.Printf("Time Stamp: %s\n",unwrap[idx].Point_time)
	fmt.Printf("Frequency: %f seconds\n",unwrap[idx].Frequency)
	fmt.Printf("Value: %f\n",unwrap[idx].Val)
	fmt.Printf("Electricity Grid Region Market: %s\n",unwrap[idx].Market)
	fmt.Printf("Type of Fuel: %s\n",unwrap[idx].Fuel)
	}		

		

}
