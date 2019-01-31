package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

func Trust(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))

	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	rs, h := services.Trust(uint(aid), rh)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(h)
}

func Verify(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))
	b, _ := ioutil.ReadAll(r.Body)

	rs, hs := services.Verify(uint(aid), b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(hs)
}

func Report(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))
	b, _ := ioutil.ReadAll(r.Body)

	rs, os := services.Report(uint(aid), b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(os)
}
