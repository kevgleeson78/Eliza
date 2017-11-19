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
*https://stackoverflow.com/questions/12430907/create-div-using-form-data-with-ajax-jquery
*https://stackoverflow.com/questions/42018775/pattern-matching-and-regular-expression-in-perl
 */

package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func Reflections(capturedString string) string {
	//adapted from https://stackoverflow.com/questions/10196462/regex-word-boundary-excluding-the-hyphen
	//To prevent "you're"  or any word with a "'" from getting split into three tokens
	boundaries := regexp.MustCompile(`(\b[^\w']|$)`)
	tokens := boundaries.Split(capturedString, -1)

	// List the reflections.
	reflections := [][]string{
		{`am`, `are`},
		{`was`, `were`},
		{`i`, `you`},
		{`i'd`, `you would`},
		{`i've`, `you have`},
		{`i'll`, `you will`},
		{`my`, `your`},
		{`are`, `am`},
		{`you've`, `I have`},
		{`you'll`, `I will`},
		{`your`, `my`},
		{`yours`, `mine`},
		{`you`, `me`},
		{`me`, `you`},
	}

	// Loop through each token, reflecting it if there's a match.
	for i, token := range tokens {
		for _, reflection := range reflections {
			if matched, _ := regexp.MatchString(reflection[0], token); matched {
				tokens[i] = reflection[1]
				break
			}
		}
	}

	// Put the tokens back together.
	//A space is need for teh regular expression (\b[^\w']|$)
	//as it dosent allow the word you're to be split into three parts.
	//If the space is not put in as the second argument it will return
	//one continuous string.
	return strings.Join(tokens, " ")
}
func requestHandler(w http.ResponseWriter, r *http.Request) {

	//set the header content type to text/html
	w.Header().Set("Content-Type", "text/html")

	//fmt.Fprintf(w, "Hello, %s! \n", r.URL.Query().Get("value"))
	fmt.Fprintf(w, Reflections(r.URL.Query().Get("value")))
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
