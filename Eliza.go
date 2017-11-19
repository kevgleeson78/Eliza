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
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Reflections(capturedString string) string {
	//adapted from https://stackoverflow.com/questions/10196462/regex-word-boundary-excluding-the-hyphen
	//To prevent "you're"  or any word with a "'" from getting split into three tokens
	boundaries := regexp.MustCompile(`(\b[^\w']|$)`)
	tokens := boundaries.Split(capturedString, -1)

	// List the reflections.
	reflections := [][]string{
		{`\bare\b`, `am`},
		{`\bam/b`, `are`},
		{`\bwere\b`, `was`},
		{`\bwas\b`, `were`},
		{`\byou\b`, `I`},
		{`\bI\b`, `you`},
		{`\bme\b`, `you`},
		{`\byour\b`, ` my`},
		{`\bmy\b`, `your`},
		{`\byou've\b`, `I've`},
		{`\bI've\b`, `you've`},
		{`\bI'm\b`, `you're`},
		{`\byou're\b`, `I'm`},
		{`\bme\b`, `you`},
		{`\byou\b`, `me`},
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

//Function ElizaResponse to take in and return a string
func ElizaResponse(str string) string {
	//Adapted from https://gist.github.com/ianmcloughlin/c4c2b8dc586d06943f54b75d9e2250fe

	/*Regex MatchString function with isolation of the word "father"
	*with a boundry ignore case regex command.
	 */
	if matched, _ := regexp.MatchString(`(?i)\bfather\b`, str);
	//Condition to replace the original string if it has the word "father"
	matched {
		return "Why don’t you tell me more about your father?"
	}
	r1 := regexp.MustCompile(`(?i)^\bI'?\s*a?m\b`)

	//Match the words "I am" and capture for replacement
	matched := r1.MatchString(str)

	//condition if "I am" is matched
	if matched {
		replacementString := "How do you know that you are"
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := Reflections(r1.ReplaceAllString(str, "$1"))
		//Concat the new opening line at the end of the function
		return replacementString + response
	}
	//Capture and replace How are you
	r1 = regexp.MustCompile(`\bI have\b`)
	matched = r1.MatchString(str)
	if matched {
		replacementString := "How long have you had "
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := Reflections(r1.ReplaceAllString(str, "$1"))
		//Concat the new opening line at the end of the function
		return replacementString + response
	}
	//Capture and replace How are you
	r1 = regexp.MustCompile(`\bMy name is\b`)
	matched = r1.MatchString(str)
	if matched {
		replacementString := "Hello "
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := Reflections(r1.ReplaceAllString(str, "$1"))
		if response == "" {
			return "I didn't get your name can you please tell me again"
		}
		//Concat the new opening line at the end of the function
		return replacementString + response
	}
	//Get random number from the length of the array of random struct
	//an array of strings for the random response
	randomResponse := []string{
		"I’m not sure what you’re trying to say. Could you explain it to me?",
		"How does that make you feel?",
		"Why do you say that?",
		"Can you tell me more about it?",
		"Can you understand me?",
		"How can I help?",
		"Are you still there?",
	}
	//Return a random index of the array
	return randomResponse[rand.Intn(len(randomResponse))]

}
func requestHandler(w http.ResponseWriter, r *http.Request) {

	//set the header content type to text/html
	w.Header().Set("Content-Type", "text/html")
	rand.Seed(time.Now().UTC().UnixNano())
	//fmt.Fprintf(w, "Hello, %s! \n", r.URL.Query().Get("value"))
	fmt.Fprintf(w, ElizaResponse(r.URL.Query().Get("value")))
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
