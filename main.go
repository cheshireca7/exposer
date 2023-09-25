package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"os"
	"os/signal"
	"syscall"
	"time"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"encoding/json"
	"encoding/base64"
	"sort"

	"gopkg.in/yaml.v3"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/uncover"
	"github.com/projectdiscovery/uncover/sources"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/fatih/color"
)

// STRUCTS
type options struct {
	query string
	engine string
	nobanner bool 
	configFile string
}

type ElasticConfig struct {
    URL      string `yaml:"URL"`
    PORT     string `yaml:"PORT"`
    USERNAME string `yaml:"USERNAME"`
    PASSWORD string `yaml:"PASSWORD"`
}

type Entry struct {
	Timestamp float64 `json:"timestamp"`
	Source string `json:"source"`
	Ports []float64 `json:"ports"`
	Host string `json:"host"`
	URL string `json:"url"`
}

// GLOBALS
var (
	green = color.New(color.FgGreen).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue = color.New(color.FgBlue).SprintFunc()
)

// FUNCTIONS

// Log messages with colors
func showMessage(risk string, m []string) {
	message := strings.Join(m, ",")
	switch risk {
	case "ok":
		fmt.Println("[" + green("OK") + "] " + message)
	case "error":
		gologger.Error().Msg(message)
	case "info":
		gologger.Info().Msg(message)
	case "warning":
	}
}

// Trap SIGINT and show custom message
func trapSigs() {
	sigChannel := make(chan os.Signal, 1)
  	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
	    	sig := <-sigChannel
    
		switch sig {
		case os.Interrupt:
			log.Println("[" + yellow("WRNG") + "] SIGINT detected, exiting ...")
			os.Exit(1)
	    	case syscall.SIGTERM:
		        fmt.Println("Killing ...")
		}
	}()
}

// Beautiful banner :)
func banner() {
	b64_banner := "ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICANCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgDQogLGFkUFBZYmEsICA4YiwgICAgICxkOCAgOGIsZFBQWWJhLCAgICAsYWRQUFliYSwgICAsYWRQUFliYSwgICAsYWRQUFliYSwgIDhiLGRQUFliYSwgIA0KYThQX19fX184OCAgIGBZOCwgLDhQJyAgIDg4UCcgICAgIjhhICBhOCIgICAgICI4YSAgSThbICAgICIiICBhOFBfX19fXzg4ICA4OFAnICAgIlk4ICANCjhQUCIiIiIiIiIgICAgICk4ODgoICAgICA4OCAgICAgICBkOCAgOGIgICAgICAgZDggICBgIlk4YmEsICAgOFBQIiIiIiIiIiAgODggICAgICAgICAgDQoiOGIsICAgLGFhICAgLGQ4IiAiOGIsICAgODhiLCAgICxhOCIgICI4YSwgICAsYTgiICBhYSAgICBdOEkgICI4YiwgICAsYWEgIDg4ICAgICAgICAgIA0KIGAiWWJiZDgiJyAgOFAnICAgICBgWTggIDg4YFliYmRQIicgICAgYCJZYmJkUCInICAgYCJZYmJkUCInICAgYCJZYmJkOCInICA4OCAgICAgICAgICANCiAgICAgICAgICAgICAgICAgICAgICAgICA4OCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgDQogICAgICAgICAgICAgICAgICAgICAgICAgODggICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIA=="
	decodedBytes, _ := base64.StdEncoding.DecodeString(b64_banner)
	fmt.Println(string(decodedBytes) + cyan("\n\n-- Monitor your favorite services exposed to the Internet 👀\n\n"))
}

// Check response body for error resturned from Elasticsearch
func checkResponse(response *esapi.Response) map[string]interface{} {

	// Parse response body
	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		showMessage("error", []string{"Unable to read JSON body from the response"})
		log.Fatalf("%s\n", err)
	}

	// Serialize response body
	var resData map[string]interface{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		showMessage("error", []string{"Unable to serialize JSON body from response"})
		log.Fatalf("%s\n", err)
	}

	if response.IsError() {
		showMessage("error", []string{"Something went wrong with the ElasticSearch server communication"})
		err := resData["error"].(map[string]interface{})["root_cause"].([]interface{})[0].(map[string]interface{})["reason"].(string)
		log.Fatalf("%s\n", err)
	}

	// Return map type if response is ok
	return resData

}

// Initialize uncover service
func uncoverInit(query string, engine string) *uncover.Service {
	opts := uncover.Options{
		Agents:   strings.Split(engine,","),
		Queries:  strings.Split(query,","),
		Limit:    50,
		MaxRetry: 2,
		Timeout:  20,
	}

	u, err := uncover.New(&opts)
	if err != nil {
		panic(err)
	}

	return u
}

// Initialize Elasticsearch connection and create new index
func elasticSearchInit(cf string) (*elasticsearch.Client, context.Context, string) {
	
	// Read configuration file
	yamlData, err := ioutil.ReadFile(cf)
	if err != nil {
		showMessage("error", []string{"Unable to read YAML configuration file"})
		log.Fatalf("%v", err)
	}

	// Parse configuration file
	var elasticConfig ElasticConfig
	if err := yaml.Unmarshal(yamlData, &elasticConfig); err != nil {
		showMessage("error", []string{"Unable to serialize YAML configuration file"})
		log.Fatalf("%v", err)
	}

	// Configure Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{"https://" + elasticConfig.URL + ":" + elasticConfig.PORT},
		Username:  elasticConfig.USERNAME,
		Password:  elasticConfig.PASSWORD,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Initiallize client
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		showMessage("error", []string{"Unable to create Elasticsearch client"})
		log.Fatalf("%v", err)
	}

	// Configure index
	ctx := context.Background()
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	indexName := timestamp + "_uncover_results"
	indexSettings := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 2
		}
	}`

	// Perforn request to create index
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(indexSettings),
	}

	if _, err := req.Do(ctx, client); err != nil {
		showMessage("error", []string{"Unable to create index"})
		log.Fatalf("%v\n", err)
	}

	// Return client, context and index created if everything went right
	return client, ctx, indexName
}


// Store the new IP 
func storeOutput(jsonData map[string]interface{}, esc *elasticsearch.Client, ctx context.Context, index string) {
	
	Source, _ := jsonData["source"].(string)
	Port, _ := jsonData["port"].(float64)
	Ports := []float64{Port}
	Host, _ := jsonData["host"].(string)
	URL, _ := jsonData["url"].(string)
	Timestamp, _ := jsonData["timestamp"].(float64)
	IP, _ := jsonData["ip"].(string)

	entry := Entry{
		Timestamp: Timestamp,
		Source: Source,
		Ports: Ports,
		Host: Host,
		URL: URL,
	}

	// Generate entry JSON body
	entryJson, err := json.Marshal(entry)
	if err != nil {
		showMessage("error", []string{"Unable to serialize JSON data for store request"})
		log.Fatalf("%s\n", err)
	}
	
	// Configure store request
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: IP,
		Body:       bytes.NewReader(entryJson),
		Refresh:    "true",
	}

	// Perform store request
	res, err := req.Do(ctx, esc)
	if err != nil {
		showMessage("error", []string{"Unable to store new IP address"})
		log.Fatalf("%s\n", err)
	}
	defer res.Body.Close()

	// Check elasticsearch response
	checkResponse(res)

}

// Check if IP already exists in the index
func checkIPAddress(ip string, esc *elasticsearch.Client, ctx context.Context, index string) (map[string]interface{}, bool){

	// Configure search request 
	query := `{
		"query": {
			"term": {
				"_id":"`+ip+`"
			}
		}
	}`
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}

	// Perform search request
	res, err := req.Do(ctx, esc)
	if err != nil {
		showMessage("error", []string{"Unable to check if IP address already exists"})
		log.Fatalf("%s\n", err)
	}
	defer res.Body.Close()

	// Check Elasticsearch response
	resData := checkResponse(res)

	// If hits = 0, it means IP does not exists
	hits := resData["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)

	if hits == 0 {
		return nil, false
	}

	// Retrun IP data if it exists within the index
	return resData, true

}

// Update IP with the new ports discovered
func storePorts(data map[string]interface{}, esc *elasticsearch.Client, ctx context.Context, index string) {
	
	// Get IP data from index
	IP, _ := data["ip"].(string)
	IPData, _ := checkIPAddress(IP, esc, ctx, index)

	// Get IP ports already stored
	hits := IPData["hits"].(map[string]interface{})
	hitsArray := hits["hits"].([]interface{})
	oldPorts := hitsArray[0].(map[string]interface{})["_source"].(map[string]interface{})["ports"].([]interface{})

	// Append new port to already stored ones
	newPort, _ := data["port"].(float64)
	for _, num := range oldPorts {
	        if num == newPort {
			return
		}
	}

	newPorts := append(oldPorts, newPort)
	intPorts := make([]float64, len(newPorts))
	for i, v := range newPorts {
		intPorts[i] = v.(float64)
	}	
	sort.Float64s(intPorts)

	// Configure update request
	updateRequest := map[string]interface{}{
		"doc": map[string]interface{}{
			"ports": intPorts,
		},
	}
	updateRequestBody, err := json.Marshal(updateRequest)
	if err != nil {
		showMessage("error", []string{"Error serializing update request body"})
		log.Fatalf("%s\n", err)
	}
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: IP,
		Body:       strings.NewReader(string(updateRequestBody)),
	}

	// Perform update request
	res, err := req.Do(ctx, esc)
	if err != nil {
		showMessage("error", []string{"Error performing update request"})
		log.Fatalf("%s\n", err)
	}

	// Check Elasticsearch response
	checkResponse(res)

}


func main() {

	trapSigs()

	// Define user flags
	opt := &options{}
	set := goflags.NewFlagSet()
	set.SetDescription(`Monitor services on the Internet in real time and store results to Elasticsearch`)
	set.StringVarP(&opt.query, "query", "q", "", "Search query")
	set.BoolVarP(&opt.nobanner, "no-banner", "nb", false, "Hide the beautiful banner")
	set.StringVarP(&opt.configFile, "configfile", "cf", "", "Specify the config file for Elasticsearch")
	set.StringVarP(&opt.engine, "engine", "e", "", "Search engine (shodan,shodan-idb,fofa,censys,quake,hunter,zoomeye,netlas,publicwww,criminalip,hunterhow,all) (default shodan)")

	// Parse user falgs
	if err := set.Parse(); err != nil {
		showMessage("error", []string{"Could not parse flags!"})
		log.Fatalf("%s\n", err)
	}

	if !opt.nobanner { 
		banner()
	}	

	if opt.engine == "all" {
		opt.engine = "shodan,shodan-idb,fofa,censys,quake,hunter,zoomeye,netlas,criminalip,publicwww,hunterhow"
	}

	// Initialize Elasticsearch communication
	esc, ctx, index := elasticSearchInit(opt.configFile)
	showMessage("info", []string{"Creating new index: " + index})

	// Initiate uncover
	u := uncoverInit(opt.query, opt.engine)	
	showMessage("info", []string{"Monitoring query: '" + opt.query + "'"})

	// Handle uncover results
	var count = 0
	result := func(result sources.Result) {

		var data map[string]interface{}
		dataJson := []byte(result.JSON())
		if err := json.Unmarshal(dataJson, &data); err != nil {
			showMessage("error", []string{"Unable to parse JSON data from uncover resutls"})
			log.Fatalf("%s\n", err)
		}
		
		// Check if IP was already stored
		ip, _ := data["ip"].(string)
		_, ok := checkIPAddress(ip, esc, ctx, index)
		if !ok {
			storeOutput(data, esc, ctx, index)
			count++
		}else{
			storePorts(data, esc, ctx, index)
		}

		fmt.Printf("\r[" + blue("INF") + "] Number of entries stored: %d", count)
	}

	// Infinite loop executing uncover query 
	for {
		if err := u.ExecuteWithCallback(context.TODO(), result); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)
	}

}
