package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gomicro/ledger"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	RegistrationTokens string `envconfig:"REGISTRATION_TOKENS" default:"validtoken"`
	Hostname           string `envconfig:"HOSTNAME" default:"0.0.0.0"`
	Port               string `envconfig:"HOSTNAME" default:"8080"`

	//Fake response data
	CustomerID  string `envconfig:"CUSTOMER_ID" default:"somecustomer"`
	ProductCode string `envconfig:"PRODUCT_CODE" default:"someproduct"`
}

var (
	log    *ledger.Ledger
	config configuration
)

type marketplacemeteringResolveCustomerInput struct {
	RegistrationToken string `json:"RegistrationToken"`
}

type marketplacemeteringResolveCustomerOutput struct {
	CustomerIdentifier string `json:"CustomerIdentifier"`
	ProductCode        string `json:"ProductCode"`
}

type marketplacemeteringBatchMeterUsageInput struct {
	ProductCode  string `json:"ProductCode"`
	UsageRecords []struct {
		CustomerIdentifier string `json:"CustomerIdentifier"`
		Dimension          string `json:"Dimension"`
		Quantity           string `json:"Quantity"`
		Timestamp          string `json:"Timestamp"`
	} `json:"UsageRecords"`
}

type marketplacemeteringBatchMeterUsageOutput struct {
	Results []struct {
		MeteringRecordID string `json:"MeteringRecordId"`
		Status           string `json:"Status"`
		UsageRecord      struct {
			CustomerIdentifier string `json:"CustomerIdentifier"`
			Dimension          string `json:"Dimension"`
			Quantity           string `json:"Quantity"`
			Timestamp          string `json:"Timestamp"`
		} `json:"UsageRecord"`
	} `json:"Results"`
	UnprocessedRecords []struct {
		CustomerIdentifier string `json:"CustomerIdentifier"`
		Dimension          string `json:"Dimension"`
		Quantity           string `json:"Quantity"`
		Timestamp          string `json:"Timestamp"`
	} `json:"UnprocessedRecords"`
}

func main() {
	configure()

	http.Handle("/", registerEndpoints())

	startService()
}

func configure() {
	log = ledger.New(os.Stdout, ledger.ParseLevel("debug"))
	log.Debug("Logger initialized")

	err := envconfig.Process("", &config)
	if err != nil {
		ledger.Errorf("Failed to process from the environment: %v", err.Error())
		os.Exit(1)
	}

	log.Debugf("Configured with (%v) for valid tokens", config.RegistrationTokens)
}

func registerEndpoints() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/ResolveCustomer/", handleResolveCustomer).Methods("POST")
	r.HandleFunc("/BatchMeterUsage/", handleBatchMeterUsage).Methods("POST")
	return r
}

func handleBatchMeterUsage(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	i := marketplacemeteringBatchMeterUsageInput{}
	json.Unmarshal(b, &i)

	log.Debugf("Received a BatchMeterUsage request %v", i.ProductCode)

	o := marketplacemeteringBatchMeterUsageOutput{}

	b, err = json.Marshal(o)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(b)
	return
}

func handleResolveCustomer(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	i := marketplacemeteringResolveCustomerInput{}
	json.Unmarshal(b, &i)

	log.Debugf("Received a ResolveCustomer request %v", i.RegistrationToken)
	found := false
	tokens := strings.Split(config.RegistrationTokens, ",")
	for _, t := range tokens {
		if i.RegistrationToken == t {
			found = true
		}
	}

	if found {
		log.Debug("Token is valid")
		o := marketplacemeteringResolveCustomerOutput{
			CustomerIdentifier: config.CustomerID,
			ProductCode:        config.ProductCode,
		}

		b, err = json.Marshal(o)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write(b)
		return
	}

	log.Debug("Token is invalid")
	w.WriteHeader(401)
	w.Write(b)
}

func startService() error {
	log.Debugf("Listening on %v:%v", config.Hostname, config.Port)
	return http.ListenAndServe(net.JoinHostPort(config.Hostname, config.Port), nil)
}
