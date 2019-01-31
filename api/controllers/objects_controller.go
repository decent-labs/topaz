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

	rs, h := services.Trust(uint(aid), rh)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(h)
}

func Verify(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid, _ := strconv.Atoi(r.Header.Get("appId"))
	h := path.Base(r.URL.Path)

	rs, hs := services.Verify(uint(aid), h)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(hs)
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
