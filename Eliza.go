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
*https://gist.github.com/ianmcloughlin/c4c2b8dc586d06943f54b75d9e2250fe
*https://github.com/data-representation/eliza/blob/master/data/responses.txt
*https://github.com/codeanticode/eliza/blob/master/data/eliza.script
*https://www.smallsurething.com/implementing-the-famous-eliza-chatbot-in-python/
*https://stackoverflow.com/questions/37274282/regex-with-replace-in-golang
*https://regexone.com/lesson/letters_and_digits?
*https://stackoverflow.com/questions/3012788/how-to-check-if-a-line-is-blank-using-regex
*
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
	//Adapted from https://gist.github.com/ianmcloughlin/c4c2b8dc586d06943f54b75d9e2250fe
	boundaries := regexp.MustCompile(`(\b[^\w']|$)`)
	tokens := boundaries.Split(capturedString, -1)

	// List of the reflections in 2d array.
	reflections := [][]string{
		{`\byou're\b`, `I'm`},
		{`\b(?i)I'm\b`, `you're`},
		{`\bare\b`, `am`},
		{`\bam\b`, `are`},
		{`\bwere\b`, `was`},
		{`\bwas\b`, `were`},
		{`\byou\b`, `I`},
		{`\bi\b`, `you`},
		{`\bme\b`, `you`},
		{`\byour\b`, ` my`},
		{`\bmy\b`, `your`},
		{`\byou've\b`, `I've`},
		{`\bI've\b`, `you've`},
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
	/*The below regular expression was adapted from https://github.com/data-representation/eliza/blob/master/data/responses.txt
	 *Compile the regular expression that says to match while ingnoring case and multi line.(?im)
	 *At the start of a line 0 or more whitespace followed by the characters I can with one or zero "'".\s*I can'?t
	 *followed by t.
	 *capture anything that is not (. !) zero or more times.([^\.!]*)
	 *followed by (. !) zero or more times [\.!]* followed by whitespace zero or more times.\s*
	 *Finally end of the line.$
	 */
	r1 := regexp.MustCompile(`(?im)^\s*I can'?t ([^\.!]*)[\.!]*\s*$`)

	//Match the words "I can't" and capture for replacement
	matched := r1.MatchString(str)

	//condition if "I can't" is matched
	if matched {
		/*Change my to your, i to you etc...
		*pass the captured string into ReplaceAllString function
		*against the original string.
		*This then gets passed to the reflections function to change nouns.
		 */
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		//Test output of string from Reflections
		//fmt.Println(reflectString)
		answers := []string{"Perhaps you could " + reflectString + " if you tried.",
			"How do you think that you can't " + reflectString + "?",
			"Have you tried ?",
			"Perhaps you could " + reflectString + " now.",
			"Do you really want to be able to" + reflectString + "?",
		}
		//replace and create the new string with the answer and changed nouns.
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//return final string to the front end.
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I need ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Why do you need " + reflectString + "?",
			"Would it really help you to get " + reflectString + "?",
			"Are you sure you need " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`Why don\'?t you ([^\?]*)\??`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Do you really think I don't " + reflectString + "?",
			"Perhaps eventually I will " + reflectString + ".",
			"Do you really want me to " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*my name is ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"HEllo " + reflectString + " how are you today?",
			"Perhaps eventually I will get to know you better " + reflectString + ".",
			"Do you like your name " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`Why can\'?t I ([^\?]*)\??`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Do you think you should be able to " + reflectString + "?",
			"If you could " + reflectString + ", what would you do?",
			"I don't know -- why can't you " + reflectString + "?",
			"Have you really tried?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I'?\s*a?m ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Did you come to me because you are " + reflectString + "?",
			"How long have you been " + reflectString + "?",
			"How do you feel about being " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}

	r1 = regexp.MustCompile(`(?im)^\s*Are you ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Why does it matter whether I am " + reflectString + "?",
			"Would you prefer it if I were not " + reflectString + "?",
			"Perhaps you believe I am " + reflectString + ".",
			"I may be " + reflectString + " -- what do you think?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*What ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Why do you ask?",
			"How would an answer to that help you?",
			"What do you think?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*How ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"How do you suppose?",
			"Perhaps you can answer your own question.",
			"What is it you're really asking?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Because ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Is that the real reason?",
			"What other reasons come to mind?",
			"Does that reason apply to anything else?",
			"If " + reflectString + ", what else must be true?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Hello ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Hello... I'm glad you could drop by today.",
			"Hi there... how are you today?",
			"Hello, how are you feeling today?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I think ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Do you doubt " + reflectString + "?",
			"Do you really think so?",
			"But you're not sure " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(.*) friend (.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Tell me more about your friends.",
			"When you think of a friend, what comes to mind?",
			"Why don't you tell me about a childhood friend?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?i)Yes`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"You seem quite sure.",
			"OK, but can you elaborate a bit?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?i)No`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Why do you answer no?",
			"OK, but can you elaborate a bit?",
			"Do you consider that a negitave thought?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(.*) sorry (.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"There are many times when no apology is needed.",
			"What feelings do you have when you apologize?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(.*) computer(.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Are you really talking about me?",
			"Does it seem strange to talk to a computer?",
			"How do computers make you feel?",
			"Do you feel threatened by computers?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Is it ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Do you think it is " + reflectString + "?",
			"Perhaps it's " + reflectString + " -- what do you think?",
			"If it were " + reflectString + ", what would you do?",
			"It could well be that " + reflectString + ".",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*It is ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"You seem very certain.",
			"If I told you that it probably isn't " + reflectString + ", what would you feel?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Can you ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"What makes you think I can't " + reflectString + "?",
			"If I could " + reflectString + ", then what?",
			"Why do you ask if I can " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Can I ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Perhaps you don't want to " + reflectString + ".",
			"Do you want to be able to " + reflectString + "?",
			"If you could " + reflectString + ", would you?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*You are ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Why do you think I am " + reflectString + "?",
			"Does it please you to think that I'm " + reflectString + "?",
			"Perhaps you would like me to be " + reflectString + ".",
			"Perhaps you're really talking about yourself?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*You'?\s*re ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Why do you say I am " + reflectString + "?",
			"Why do you think I am " + reflectString + "?",
			"Are we talking about you, or me?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I don'?\s*t ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Don't you really " + reflectString + "?",
			"Why don't you " + reflectString + "?",
			"Do you want to " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I feel ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Good, tell me more about these feelings.",
			"Do you often feel " + reflectString + "?",
			"When do you usually feel " + reflectString + "?",
			"When you feel " + reflectString + ", what do you do?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I have ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Why do you tell me that you've " + reflectString + "?",
			"Have you really " + reflectString + "?",
			"Now that you have " + reflectString + ", what will you do next?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I would ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Could you explain why you would " + reflectString + "?",
			"Why would you " + reflectString + "?",
			"Who else knows that you would " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Is there ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Do you think there is " + reflectString + "?",
			"It's likely that there is " + reflectString + ".",
			"Would you like there to be " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*My ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"I see, your " + reflectString + ".",
			"Why do you say that your " + reflectString + "?",
			"When your " + reflectString + ", how do you feel?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*You ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"We should be discussing you, not me.",
			"Why do you say that about me?",
			"Why do you care whether I " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*Why ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Why don't you tell me the reason why " + reflectString + "?",
			"Why do you think " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*I want ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"What would it mean to you if you got " + reflectString + "?",
			"Why do you want " + reflectString + "?",
			"What would you do if you got " + reflectString + "?",
			"If you got " + reflectString + ", then what would you do?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(.*) mother(.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Tell me more about your mother.",
			"What was your relationship with your mother like?",
			"How do you feel about your mother?",
			"How does this relate to your feelings today?",
			"Good family relations are important.",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(.*) father(.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Tell me more about your father.",
			"How did your father make you feel?",
			"How do you feel about your father?",
			"Does your relationship with your father relate to your feelings today?",
			"Do you have trouble showing affection with your family?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}

	r1 = regexp.MustCompile(`(.*) child(.*)`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Did you have close friends as a child?",
			"What is your favorite childhood memory?",
			"Do you remember any dreams or nightmares from childhood?",
			"Did the other children sometimes tease you?",
			"How do you think your childhood experiences relate to your feelings today?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*hello([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Hi How are you?",
			"I hope you are keeping well.",
			"Hello I am happy you are talking with me.",
			"Hello what are you thinking of asking me?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	//adapted from https://stackoverflow.com/questions/3012788/how-to-check-if-a-line-is-blank-using-regex
	r1 = regexp.MustCompile(`^\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{
			"Can you please fill out the text box so i can speak with you?",
			"How can i talk with you if you just put empty text into the chat box?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}

	//Get random number from the length of the array of random struct
	//an array of strings for the random response
	randomResponse := []string{
		"I’m not sure what you’re trying to say. Could you explain it to me?",
		"How does that make you feel?",
		"Why do you say that?",
		"Please go on",
		"Can you tell me more about it?",
		"Can you understand me?",
		"How can I help?",
		"What does that suggest to you?",
		"Do you feel strongly about discussing such things?",
		"Are you still there?",
		"What have you been up to recently?",
		"How are you finding the conversation we are having?",
	}
	//Return a random index of the array
	return randomResponse[rand.Intn(len(randomResponse))]

}
func requestHandler(w http.ResponseWriter, r *http.Request) {

	//set the header content type to text/html
	w.Header().Set("Content-Type", "text/html")

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
	rand.Seed(time.Now().UTC().UnixNano())
	http.ListenAndServe(":8080", nil)
}
