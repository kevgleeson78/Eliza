/*
*App-Name: Go-guessing-game
*@Author:  Kevin Gleeson
*Date:     15/10/2017
*Version:  1.0
*Sources:
*https://github.com/data-representation/go-echo
*https://golang.org/pkg/net/http/#SetCookie
*https://stackoverflow.com/questions/12130582/setting-cookies-in-golang-net-http
*https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
*https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html
*https://stackoverflow.com/questions/22593259/check-if-string-is-int-golang
*https://stackoverflow.com/questions/28159520/passing-a-query-parameter-to-the-go-http-request-handler-using-the-mux-package
*https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
*https://github.com/gowww/client/blob/master/response.go
*https://godoc.org/hkjn.me/googleauth
*https://golang.org/pkg/strconv/
*https://stackoverflow.com/questions/26189523/go-represent-path-without-query-string
*https://stackoverflow.com/questions/20320549/how-can-you-delete-a-cookie-in-an-http-response
*https://github.com/github/gitignore/blob/master/Go.gitignore
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
