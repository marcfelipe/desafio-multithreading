package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Address      string `json:"address"`
	Neighborhood string `json:"neighborhood"`
	Service      string `json:"service"`
}

func main() {
	searchCep := "69054734"

	channelViaCEP := make(chan CEP)
	channelBrAPI := make(chan CEP)

	go BuscaCep(searchCep, channelBrAPI, BrasilCEP{})
	//time.Sleep(time.Second * 2)				//ViaCEP faster, than to test BrasilCEP as first this delay is added
	go BuscaCep(searchCep, channelViaCEP, ViaCEP{})

	select {
	case ret := <-channelViaCEP:
		fmt.Printf("Retornado por: %s \n dados:\n %+v", ret.Service, ret)

	case ret := <-channelBrAPI:
		fmt.Printf("Retornado por: %s \n dados:\n %+v", ret.Service, ret)

	case <-time.After(time.Second):
		println("Timeout. No CEP Found.")
	}

}

func BuscaCep[V ViaCEP | BrasilCEP](cep string, ch chan CEP, serviceSearch V) {
	var service string
	if reflect.TypeOf(serviceSearch) == reflect.TypeOf(ViaCEP{}) {
		service = "ViaCEP"
		req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		var res ViaCEP
		err = json.Unmarshal(body, &res)
		if err != nil {
			panic(err)
		}
		ch <- CEP{
			Cep:          strings.ReplaceAll(res.Cep, "-", ""),
			State:        res.Uf,
			City:         res.Localidade,
			Address:      res.Logradouro,
			Neighborhood: res.Bairro,
			Service:      service,
		}
	}
	if reflect.TypeOf(serviceSearch) == reflect.TypeOf(BrasilCEP{}) {
		service = "BrasilCEP"
		req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		var res BrasilCEP
		err = json.Unmarshal(body, &res)
		if err != nil {
			panic(err)
		}
		// time.Sleep(time.Second * 3)
		ch <- CEP{
			City:         res.City,
			State:        res.State,
			Cep:          res.Cep,
			Address:      res.Street,
			Neighborhood: res.Neighborhood,
			Service:      service,
		}

	}
}
