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

	// List the reflections.
	reflections := [][]string{
		{`\b(?i)I'm\b`, `you're`},
		{`\bare\b`, `am`},
		{`\bam\b`, `are`},
		{`\bwere\b`, `was`},
		{`\bwas\b`, `were`},
		{`\byou\b`, `I`},
		{`\b(?i)I\b`, `you`},
		{`\bme\b`, `you`},
		{`\byour\b`, ` my`},
		{`\bmy\b`, `your`},
		{`\byou've\b`, `I've`},
		{`\bI've\b`, `you've`},

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
		//Test output of string commimg out of Reflections
		//fmt.Println(reflectString)
		answers := []string{"Perhaps you could " + reflectString + " if you tried.",
			"How do you think that you can't " + reflectString + "?",
			"Have you tried ?",
			"Perhaps you could " + reflectString + "now.",
			"Do you really want to be able to" + reflectString + "?",
		}
		//replace and create the new string with the answer and changed nouns.
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//return final string to the front end.
		return response
	}
	//Capture and replace I have
	//The below regular expression was adapted from https://github.com/data-representation/eliza/blob/master/data/responses.txt
	r1 = regexp.MustCompile(`(?im)^\s*I have ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))

		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, "How long have you "+reflectString+"?")
		//Concat the new opening line at the end of the function
		return response
	}
	//Capture and replace My name is
	r1 = regexp.MustCompile(`(?im)^\s*My name is ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		//Responses adapated from https://github.com/codeanticode/eliza/blob/master/data/eliza.script
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Hello " + reflectString + " how have you been?",
			reflectString + " I dont think I recognise that name.",
			reflectString + " is the same name as my Father.",
			"Where does " + reflectString + " originate from?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i remember ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{" Do you often think of " + reflectString + "?",
			" Does thinking of " + reflectString + " bring anything else to mind ?",
			" What else do you recollect ?",
			" Why do you recollect " + reflectString + " just now ?",
			" What in the present situation reminds you of " + reflectString + " ?",
			" What is the connection between me and " + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*sorry ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		answers := []string{"Please don't apologise.",
			"Apologies are not necessary.",
			"I've told you that apologies are not required.",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*do you remember ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Did you think I would forget " + reflectString + "?",
			"Why do you think I should recall " + reflectString + " now ?",
			"What about " + reflectString + "?",
			"You mentioned " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*if ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Do you think its likely that " + reflectString + "?",
			"Do you wish that " + reflectString + "?",
			"What do you know about" + reflectString + "?",
			"Really, if " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*am i ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Do you believe you are " + reflectString + "?",
			"Would you want to be " + reflectString + "?",
			"Do you wish I would tell you you are" + reflectString + "?",
			"What would it mean if you were" + reflectString + "?",
		}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}

	r1 = regexp.MustCompile(`(?im)^\s*are you ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Why are you interested in whether I am " + reflectString + " or not ?",
			"Would you prefer if I weren't " + reflectString + "?",

			"Do you sometimes think I am " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*are ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Did you think they might not be " + reflectString + "?",
			"Would you like it if they were not " + reflectString + "?",
			"What if they were not " + reflectString + "?",
			"Possibly they are " + reflectString + "."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*your ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Why are you concerned over my " + reflectString + "?",
			"What about your own " + reflectString + "?",
			"Are you worried about someone else's " + reflectString + "?",
			"Really, my " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*was i ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"What if you were " + reflectString + "?",
			"Do you think you were " + reflectString + "?",
			"Were you " + reflectString + "?",
			"What would it mean if you were " + reflectString + "?",
			"What does " + reflectString + " suggest to you ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i was ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Were you really ?",
			"Why do you tell me you were " + reflectString + " now ?",
			"Perhaps I already know you were" + reflectString + "."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*were you ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Would you like to believe I was " + reflectString + "?",
			"What suggests that I was " + reflectString + "?",
			"What do you think ?",
			"Perhaps I was " + reflectString + ".",
			"What if I had been " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i am sad ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"I am sorry to hear that you are " + reflectString + ".",
			"Do you think that coming here will help you not to be " + reflectString + "?",
			"I'm sure it's not pleasant to be " + reflectString + ".",
			"Can you explain what made you " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i want ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"What would it mean to you if you got " + reflectString + "?",
			"Why do you want " + reflectString + "?",
			"Suppose you got " + reflectString + " soon ?",
			"What if you never got " + reflectString + "?",
			"What would getting " + reflectString + " mean to you ?",
			"What does wanting" + reflectString + "have to do with this discussion ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\si am happy ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"How have I helped you to be" + reflectString + "?",
			"Has your treatment made you " + reflectString + "?",
			"What makes you " + reflectString + " just now ?",
			"Can you explan why you are suddenly " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i believe ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Do you really think so ?",
			"But you are not sure you " + reflectString + ".",
			"Do you really doubt you " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i am ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Is it because you are " + reflectString + "that you came to me ?",
			"How long have you been " + reflectString + "?",
			"Do you believe it is normal to be " + reflectString + "?",
			"Do you enjoy being " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i don't ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Don't you really " + reflectString + "?",
			"Why don't you " + reflectString + "?",
			"Do you wish to be able to" + reflectString + "?",
			"Does that trouble you ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*i feel ([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Tell me more about such feelings.",
			"Do you often feel " + reflectString + "?",
			"Do you enjoy feeling " + reflectString + "?",
			"Of what does feeling " + reflectString + " remind you ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*you are([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"What makes you think I am " + reflectString + "?",
			"Does it please you to believe I am " + reflectString + "?",
			"Do you sometimes wish you were " + reflectString + "?",
			"Perhaps you would like to be " + reflectString + "."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*yes([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"You seem to be quite positive.",
			"You are sure.",
			"I see.",
			"I understand."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*no([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"Are you saying no just to be negative?",
			"You are being a bit negative.",
			"Why not ?",
			"Why 'no' ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*my([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{
			"Your " + reflectString + "?",
			"What else comes to mind when you think of your" + reflectString + "?",
			"Why do you say your " + reflectString + "?",
			"Does that suggest anything else which belongs to you ?",
			"Is it important that your " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*can you([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"You believe I can " + reflectString + "don't you ?",
			"You want me to be able to " + reflectString + ".",
			"Perhaps you would like to be able to " + reflectString + "yourself."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*can i([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Whether or not you can " + reflectString + " depends on you more than me.",
			"Do you want to be able to " + reflectString + "?",
			"Perhaps you don't want to " + reflectString + "."}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*what([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"Why do you ask ?",
			"Does that question interest you ?",
			"What is it you really wanted to know ?",
			"Are such questions much on your mind ?",
			"What answer would please you most ?",
			"What do you think ?",
			"What comes to mind when you ask that ?",
			"Have you asked such questions before ?",
			"Have you asked anyone else ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*because([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {

		answers := []string{"Is that the real reason ?",
			"Don't any other reasons come to mind ?",
			"Does that reason seem to explain anything else ?",
			"What other reasons might there be ?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*why don't you([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Do you believe I don't " + reflectString + "?",
			"Perhaps I will " + reflectString + " in good time.",
			"Should you " + reflectString + " yourself ?",
			"You want me to " + reflectString + "?"}
		//Only keep the captured part of the string
		//Pass in everything after the captured part of the statement to the function Reflections
		response := r1.ReplaceAllString(str, answers[rand.Intn(len(answers))])
		//Concat the new opening line at the end of the function
		return response
	}
	r1 = regexp.MustCompile(`(?im)^\s*why can't i([^\.!]*)[\.!]*\s*$`)
	matched = r1.MatchString(str)
	if matched {
		reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
		answers := []string{"Do you think you should be able to " + reflectString + "?",
			"Do you want to be able to " + reflectString + "?",
			"Do you believe this will help you to " + reflectString + "?",
			"Have you any idea why you can't " + reflectString + "?"}
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
