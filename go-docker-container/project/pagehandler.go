package main

import "net/http"

type ViewData struct {
	EnvVariable1 string `json:"EnvVariable1"`
}

//Handles the Index requests
func index(w http.ResponseWriter, r *http.Request) {
	vd := ViewData{
		EnvVariable1: TEST_VAR_ENV42069,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", vd)
	HandleError(w, err1)
}
