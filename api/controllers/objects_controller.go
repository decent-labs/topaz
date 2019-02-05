package controllers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

func Trust(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))

	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	rs, o := services.Trust(uint(aid), rh)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(o)
}

func TrustUpdate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))

	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	uuid := path.Base(r.URL.Path)

	rs, h := services.TrustUpdate(uint(aid), uuid, rh)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(h)
}

func Verify(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))
	o := path.Base(r.URL.Path)

	rs, os := services.Verify(uint(aid), o)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(os)
}

func Report(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))
	s, _ := strconv.Atoi(path.Base(path.Dir(r.URL.Path)))
	e, _ := strconv.Atoi(path.Base(r.URL.Path))

	rs, os := services.Report(uint(aid), s, e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(os)
}
