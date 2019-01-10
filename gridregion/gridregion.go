//******************************************************
//* gridregion - return the name of the grid region (balancing authority)
//* given a lat/long
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes three (one optional) arguments on command line:
//* -debug  displays debug information
//* -lat latitude
//* -long longitude to use

//* defaults to lat=42.372 and long=-72.519 (examples on WattTime
//* web site)
//*
//* see watttime.org for an interactive map to figure out
//* your grid designation interactively
//*
//*****************************************************


package main

import "fmt"
import "net/http"
import "io/ioutil"
import "flag"
import "os"

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
	Abbrev string `json:"abbrev"`
	Id    int `json:"id"`
        Name  string `json:"name"`
}
	

func Check(err error, msg string)  {
     if err != nil {
     fmt.Printf("Error: %s\n%v\n",msg,err)
     os.Exit(-1)
     }
}
func main() {

	//
	// lat - latitude to use
	// long - longitude to use
	// debug - debug flag
	//
        version := 0.0
	var lat float64
	var long float64
	boolPtr := flag.Bool("debug", false, "Debug flag")

        flag.Float64Var(&lat,"lat",42.372,"Latitude to use")
        flag.Float64Var(&long,"long",-72.519,"Longitude to use")
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


	fmt.Printf("Grid Region Name for Latitude/Longitude: %f, %f\n",lat,long)
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
        slat := fmt.Sprintf("%08.3f",lat)
        slong := fmt.Sprintf("%08.3f",long)

	gridstr := "https://api2.watttime.org/v2/ba-from-loc/?latitude="
	gridstr = gridstr + slat
	gridstr = gridstr + "&longitude=" + slong
        if debug {
		fmt.Printf("request: %s\n",gridstr)
	}
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

        var unwrap Datadef
	err = json.Unmarshal(response,&unwrap)
        Check(err,"Error unmarshalling response")
//
	fmt.Printf("Abbreviation: %s\n",unwrap.Abbrev)
	fmt.Printf("ID: %d\n",unwrap.Id)
	fmt.Printf("Name: %s\n",unwrap.Name)


		

}
