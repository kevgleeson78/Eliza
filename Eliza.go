/*
*App-Name: Eliza
*@Author:  Kevin Gleeson
*Date:     18/11/2017
*Version:  1.0
*Sources:
*http://api.jquery.com/jquery.ajax/
*
*
 */

package main

import (
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {

	//set the header content type to text/html
	w.Header().Set("Content-Type", "text/html")

}
func main() {
	//store the directory where the html and template files are held
	fs := http.FileServer(http.Dir("src"))
	//Start at the root directory
	http.Handle("/", fs)
	//select the index.html file
	http.HandleFunc("/index", requestHandler)

	//Listen out for requests to the server
	http.ListenAndServe(":8080", nil)
}
