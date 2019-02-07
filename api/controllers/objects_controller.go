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
	aid := r.Context().Value(models.AppID).(uint)

	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	rs, o := services.Trust(aid, rh)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(o)
}

func TrustUpdate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid := r.Context().Value(models.AppID).(uint)
	uuid := path.Base(r.URL.Path)

	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	rs, h := services.TrustUpdate(aid, uuid, rh)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(h)
}

func Verify(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid := r.Context().Value(models.AppID).(uint)
	o := path.Base(r.URL.Path)

	rs, os := services.Verify(aid, o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(os)
}

func Report(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid := r.Context().Value(models.AppID).(uint)
	s, _ := strconv.Atoi(path.Base(path.Dir(r.URL.Path)))
	e, _ := strconv.Atoi(path.Base(r.URL.Path))

	rs, os := services.Report(aid, s, e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(os)
}
