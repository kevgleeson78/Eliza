# Eliza

This is a repository is an example of the Eliza chatbot with the [Go](https://golang.org/) programming language.

I have created this application as a project for the subject [Data Representation and Querying](https://data-representation.github.io/)
in the third year of software development at [GMIT](http://gmit.ie) Galway.

Author: [Kevin Gleeson](https://github.com/kevgleeson78)

Third year student at:[GMIT](http://gmit.ie) Galway

## Cloning, compiling and running the application.

1: Download [git](https://git-scm.com/downloads) to your machine if not already installed.

1.1: Download [go](https://golang.org/dl/) if not already installed.

2: Open git bash and cd to the folder you wish to hold the repository.
Alternatively you can right click on the folder and select git bash here.
This will open the git command prompt in the folder selected.
 
 3: To clone the repository type the following command in the terminal making sure you are in the folder needed for the repository.
```bash
>git clone https://github.com/kevgleeson78/Eliza
```
4: To compile the application cd to the folder and type the following 
```bash
> go build 
```
This will compile and create an executable file from the .go file and give it the name of the folder.

5: To run the application ensure you cd to folder the application is held.
Type the following command
```bash
>./Eliza
```
6: This will run the application As a server listening out for requests from the browser on you machine on port :8080

7: Once the application is running open a browser of your choice and type into the address bar http://localhost:8080/ this will bring you to the Eliza home page.

8: From there just type into the form and submit to get a response from Eliza. 

# Application Functionality 
This application is influenced by the Eliza chatbot created by Joseph Weizenbaum in 1966.

It takes a text input from a user and mimics a human response by analyzing the string that has been taken into the system. It then returns an answer related to the original text. For example, the user types into the terminal “I don’t like Mondays” Eliza could respond with “Can you tell me why you don’t like Mondays?” etc.
Within this application the go programming language was used to serve a web page using http request (http package) on port 8080. 

A html form is used to communicate with the server via AJAX on the browser (client side).
An initial message is displayed to the user to prompt a text input.
Once the user types in a message it is posted on the GUI Display screen inside a dynamically created div via JQuery. This div has its own class for individual styling. 
It is also retrieved and sent to a function on the backend to begin manipulating the string for Eliza’s response.
The input string is retrieved by the requestHandler function within the go application and then passed to the ElizaResponse function.
## Request Handler
```Go
func requestHandler(w http.ResponseWriter, r *http.Request) {

    //set the header content type to text/html
    w.Header().Set("Content-Type", "text/html")

    //fmt.Fprintf(w, "Hello, %s! \n", r.URL.Query().Get("value"))
    fmt.Fprintf(w, ElizaResponse(r.URL.Query().Get("value")))
}
```
## Eliza response function
The Eliza response function then takes the user string and performs a regular expression on the string in this example it is checking it the string begins with the words I can’t. 
A conditional is used to check if the phrase or word has been matched. 
If it has been matched, everything after “I can’t” is captured and passed on to the ReplaceAllStrings function which in turn is passed to the Reflections Function. 
```Go
//Function ElizaResponse to take in and return a string
func ElizaResponse(str string) string {
    r1 := regexp.MustCompile(`(?im)^\s*I can'?t ([^\.!]*)[\.!]*\s*$`)

    //Match the words "I can't" and capture for replacement
    matched := r1.MatchString(str)

    //condition if "I can't" is matched
    if matched {
```
## Replace all strings function

The replaceAllStrings Function takes the original string and replaces is with only the captured portion of the string denoted by $1 for first captured group. 
For example, “I can’t go to work today” would become “go to work today”.
```Go
reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
```
## Reflections function
The Reflections function Takes each word in the captured string and changes the pronouns for example:
You are not happy with my questions would become, I am not happy with your questions.

This is achieved by splitting every word in the string on its boundary (space) and storing it in a variable called tokens. A loop is then used to check if each word matches the first row of the reflection array. I if does it will then change it to it’s mapped pronoun and the loop is ended to stop the word from being matched again.
 The string is then retuned and Joined back together using strings.join() which takes each token a re constructs the string.
 ```GO
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
```

The result of both functions is then stored in a variable called reflect string.
```Go
reflectString := Reflections(r1.ReplaceAllString(str, "$1"))
```

An array of answers can be used to randomly return an answer with the reflected string concatenated to the answer if needed.
```Go
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
```

The final response string is then returned to the application GUI and is held within another dynamically created div via JQuery.
```javaScript
$('#main-view').append('<div class="eliza-response" id="eliza-output-area"><p>'+(data)+'</p></div>').html();
```
A setTimeout function was used set to the ajax response to create the illusion of Eliza Typing the return message. 

Resources used to create this app:

https://startbootstrap.com/template-overviews/bare

http://api.jquery.com/jquery.ajax

https://stackoverflow.com/questions/9372033/how-do-i-pass-parameters-that-is-input-textbox-value-to-ajax-in-jquery

https://github.com/data-representation/go-ajax/blob/master/static/index.html

https://stackoverflow.com/questions/23805443/remove-the-form-input-fields-data-after-click-on-submit

https://stackoverflow.com/questions/12430907/create-div-using-form-data-with-ajax-jquery

https://stackoverflow.com/questions/42018775/pattern-matching-and-regular-expression-in-perl

https://gist.github.com/ianmcloughlin/c4c2b8dc586d06943f54b75d9e2250fe

https://github.com/data-representation/eliza/blob/master/data/responses.txt

https://github.com/codeanticode/eliza/blob/master/data/eliza.script

https://www.smallsurething.com/implementing-the-famous-eliza-chatbot-in-python

https://stackoverflow.com/questions/37274282/regex-with-replace-in-golang

https://regexone.com/lesson/letters_and_digits

https://stackoverflow.com/questions/3012788/how-to-check-if-a-line-is-blank-using-regex





