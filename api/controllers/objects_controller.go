package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/services"
)

func Trust(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	appID, _ := strconv.Atoi(r.Header.Get("appId"))
	body, _ := ioutil.ReadAll(r.Body)

	responseStatus, app := services.Trust(uint(appID), body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(app)
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
