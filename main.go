package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gomicro/ledger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/gomicro/steward"
)

type configuration struct {
	RegistrationTokens string `envconfig:"REGISTRATION_TOKENS" default:"validtoken"`
	Hostname           string `envconfig:"HOST" default:"0.0.0.0"`
	Port               string `envconfig:"PORT" default:"8080"`

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

type UsageRecord struct {
	CustomerIdentifier string `json:"CustomerIdentifier"`
	Dimension          string `json:"Dimension"`
	Quantity           int    `json:"Quantity"`
	Timestamp          int    `json:"Timestamp"`
}

type marketplacemeteringBatchMeterUsageInput struct {
	ProductCode  string        `json:"ProductCode"`
	UsageRecords []UsageRecord `json:"UsageRecords"`
}

type UsageRecordResult struct {
	MeteringRecordID string      `json:"MeteringRecordId"`
	Status           string      `json:"Status"`
	UsageRecord      UsageRecord `json:"UsageRecord"`
}

type marketplacemeteringBatchMeterUsageOutput struct {
	Results            []UsageRecordResult `json:"Results"`
	UnprocessedRecords []UsageRecord       `json:"UnprocessedRecords"`
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

	rand.Seed(time.Now().UnixNano())
	o := marketplacemeteringBatchMeterUsageOutput{}
	for idx := range i.UsageRecords {
		r := i.UsageRecords[idx]
		switch rand.Intn(2) {
		case 0:
			o.UnprocessedRecords = append(o.UnprocessedRecords, r)
		case 1:
			rando, _ := uuid.NewRandom()
			o.Results = append(o.Results, UsageRecordResult{
				UsageRecord:      r,
				Status:           "Success",
				MeteringRecordID: rando.String(),
			})
		}
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
