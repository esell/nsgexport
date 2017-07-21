package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	configFile   = flag.String("c", "conf.json", "config file location")
	parsedconfig = conf{}
	NSGList      = make(map[string]TemplateNSG)
)

type conf struct {
	Subscription  string `json:"subscription"`
	ClientID      string `json:"clientid"`
	ClientSecret  string `json:"clientsecret"`
	TenantName    string `json:"tenantname"`
	ResourceGroup string `json:"resourcegroup"`
}

func httpReq(authToken string, method string, URL string, postData []byte, isJson bool) (*http.Response, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(postData))
	if err != nil {
		log.Println("Error with building request for "+URL+": ", err)
		return &http.Response{}, err
	}
	if authToken != "" {
		req.Header.Add("Authorization", "Bearer "+authToken)
	}
	if isJson {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error with request for "+URL+": ", err)
		return &http.Response{}, err
	}

	return resp, nil
}

func getAuthToken(clientID, clientSecret, tenantName string) (AuthTokenResp, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("client_id", parsedconfig.ClientID)
	data.Add("client_secret", parsedconfig.ClientSecret)
	data.Add("resource", "https://management.azure.com/")

	r, _ := httpReq("", "POST", "https://login.microsoftonline.com/"+parsedconfig.TenantName+"/oauth2/token", []byte(data.Encode()), false)
	var authResp AuthTokenResp
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error readying auth response body: ", err)
	}
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		log.Println("Error unmarshaling auth response: ", err)
	}

	return authResp, nil
}

func listNSGs() {
	data := url.Values{}
	data.Set("api-version", "2016-09-01")
	authResp, err := getAuthToken(parsedconfig.ClientID, parsedconfig.ClientSecret, parsedconfig.TenantName)
	if err != nil {
		log.Println("error getting auth token: ", err)
	}

	s, _ := httpReq(authResp.AccessToken, "GET", "https://management.azure.com/subscriptions/"+parsedconfig.Subscription+"/resourceGroups/"+parsedconfig.ResourceGroup+"/providers/Microsoft.Network/networkSecurityGroups?"+data.Encode(), nil, true)
	var nsgResp NSGResponse
	defer s.Body.Close()
	body, err := ioutil.ReadAll(s.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
	}
	err = json.Unmarshal(body, &nsgResp)
	if err != nil {
		log.Println("error converting to JSON: ", err)
	}

	for _, v := range nsgResp.Value {
		NSGList[v.Name] = v.Transform()
	}
}

func buildFinalConfig() {

	finalTemplate := ARMTemplate{
		ContentVersion: "1.0.0.0",
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	}
	for _, v := range NSGList {
		blah, err := json.Marshal(v)
		if err != nil {
			log.Println("error marshalling")
		}
		finalTemplate.Resources = append(finalTemplate.Resources, blah)
	}
	finalBlah, err := json.Marshal(finalTemplate)
	if err != nil {
		log.Println("da error")
	}
	fmt.Printf("%s\n", finalBlah)
}

func main() {
	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal("unable to read config file, exiting...")
	}
	if err := json.Unmarshal(file, &parsedconfig); err != nil {
		log.Fatal("unable to marshal config file, exiting...")
	}

	listNSGs()
	buildFinalConfig()
}
