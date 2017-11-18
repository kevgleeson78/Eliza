/*
*App-Name: Eliza
*@Author:  Kevin Gleeson
*Date:     18/11/2017
*Version:  1.0
*Sources:
*http://api.jquery.com/jquery.ajax/
*https://stackoverflow.com/questions/9372033/how-do-i-pass-parameters-that-is-input-textbox-value-to-ajax-in-jquery
*https://github.com/data-representation/go-ajax/blob/master/static/index.html
*https://stackoverflow.com/questions/23805443/remove-the-form-input-fields-data-after-click-on-submit
 */

package main

import (
	"fmt"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {

	//set the header content type to text/html
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "Hello, %s! \n", r.URL.Query().Get("value"))

}
func main() {
	//store the directory where the html and template files are held
	fs := http.FileServer(http.Dir("src"))
	//Start at the root directory
	http.Handle("/", fs)
	//select the index.html file
	http.HandleFunc("/index", requestHandler)
	http.HandleFunc("/input-text", requestHandler)
	//Listen out for requests to the server
	http.ListenAndServe(":8080", nil)
}
