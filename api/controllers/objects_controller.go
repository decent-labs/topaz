package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/models"
	"github.com/decentorganization/topaz/api/services"
)

func Trust(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	bytes, _ := ioutil.ReadAll(r.Body)

	requestObject := new(models.Object)
	requestObject.AppID, _ = strconv.Atoi(r.Header.Get("appId"))
	requestObject.DataBlob = bytes

	responseStatus, app := services.Trust(requestObject)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(app)
}
