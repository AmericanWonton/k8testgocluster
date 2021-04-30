package main

import "net/http"

//Handles the Index requests
func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("DEBUG: here we are in index: \n")
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err1)
}
