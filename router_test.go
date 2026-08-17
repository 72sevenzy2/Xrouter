package main

import (
	"encoding/json"
	"net/http"

	"github.com/72sevenzy2/json-parser/helpers"
)

type Entity struct {
	User string `json:"user"`
	Id   int    `json:"id"`
}

func HiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Entity // initilising the entity

		err := json.NewDecoder(r.Body).Decode(&i) // decoding the body to get the data we want
		if err != nil {                           // if there is no data which we needed in the body, throw an json error msg
			helpers.Failed(w)
			return
		}
		// respond with json returning the users user and the users id
		v, err := json.Marshal(i)
		if err != nil {
			helpers.Failed(w)
			return
		}
	
		helpers.Ok(w, v)
	}
}
