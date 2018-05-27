package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"strconv"
	"strings"

	"fmt"

	flag "github.com/ogier/pflag"
)

// state object for the switches and options
type ewState struct {
	quantity        int
	verbose         bool
	wordCount       int
	outputFile      string
	minimumChars    int
	maximumChars    int
	randomCaps      bool
	capsWordIndex   int
	preventCaps     bool
	preventSpecials bool
	preventInteger  bool
	wordList        map[int]string
}

// stuff to make sure this doesn't go on forever
const sanityCheck = 100000

var sanity int

// logger
var o levelOutput

func main() {
	// declarations
	var s ewState
	o.InitColours()

	flag.Usage = func() {
		banner := `
'########:'########:'########:'##:::::'##::'#######::'########::'########:::'######::
 ##.....:: ##.....:: ##.....:: ##:'##: ##:'##.... ##: ##.... ##: ##.... ##:'##... ##:
 ##::::::: ##::::::: ##::::::: ##: ##: ##: ##:::: ##: ##:::: ##: ##:::: ##: ##:::..::
 ######::: ######::: ######::: ##: ##: ##: ##:::: ##: ########:: ##:::: ##:. ######::
 ##...:::: ##...:::: ##...:::: ##: ##: ##: ##:::: ##: ##.. ##::: ##:::: ##::..... ##:
 ##::::::: ##::::::: ##::::::: ##: ##: ##: ##:::: ##: ##::. ##:: ##:::: ##:'##::: ##:
 ########: ##::::::: ##:::::::. ###. ###::. #######:: ##:::. ##: ########::. ######::
........::..::::::::..:::::::::...::...::::.......:::..:::::..::........::::......:::

Author:  Morgaine "sectorsect" Timms
License: MIT
Warning: Some of the following options when used in combination can
         significantly weaken the pass-phrases generated. 
         You probably know what you are doing though, yeah?`

		fmt.Printf("%v\n", strings.Replace(banner, "#", o.magenta.Sprintf("#"), -1))
		fmt.Println("\n\nThings It Does:")
		flag.PrintDefaults()
	}

	// basic features
	flag.IntVarP(&s.quantity, "quantity", "q", 1, "Number of passphrases to generate")
	flag.BoolVarP(&s.verbose, "verbose", "v", false, "Show all the things")
	flag.IntVarP(&s.wordCount, "wordcount", "w", 4, "Number of words per passphrase.")
	flag.StringVarP(&s.outputFile, "output-to-file", "o", "", "Filepath to output to")

	// allow char specs
	flag.IntVarP(&s.minimumChars, "minimum-chars", "m", -1, "Minimum characters for passphrases.")
	flag.IntVarP(&s.maximumChars, "maximum-chars", "M", -1, "Maximum characters for passphrases. Truncates.")

	// allow capsWordIndex of a word or random
	flag.BoolVarP(&s.randomCaps, "random-caps", "R", false, "Capitalise a random character in each passphrase")
	flag.IntVarP(&s.capsWordIndex, "caps-position", "c", -1, "Capitalise the word at this index.")

	// let passwords be terrible against my better judgement
	// you probably know what you are doing, right?
	flag.BoolVarP(&s.preventSpecials, "prevent-special", "S", false, "Prevents special chars in passphrase")
	flag.BoolVarP(&s.preventCaps, "prevent-caps", "C", false, "Prevents capitalisation in passphrase")
	flag.BoolVarP(&s.preventInteger, "prevent-int", "I", false, "Prevents integers in passphrase")
	flag.Parse()

	// assume leftover arguments are meant to be the number of passwords to generate unless there are too many
	if len(flag.Args()) > 1 {
		o.Warn.Println("Too many non-flagged arguments found. Exiting.")
		flag.Usage()
		os.Exit(1)
	}
	if s.quantity == 1 && len(flag.Args()) == 1 {
		argsQuantity, err := strconv.Atoi(flag.Args()[0])
		if err != nil {
			o.Fatality.Fatalln(err)
		}
		s.quantity = argsQuantity
	}

	// init loggers
	if s.verbose {
		o.Init(4, false, false)
	} else {
		o.Init(3, false, false)
	}

	// now we are ready, get the wordlist into state and generate
	s.wordList = getEffDiceMap()
	outputList, err := generatePassphrases(s)
	if err != nil {
		o.Warn.Fatalln(err)
	}

	// if there is a filepath stated, attempt to write there otherwise to console
	if len(s.outputFile) > 0 {
		fileOutput(outputList, s.outputFile)
	} else {
		for _, v := range outputList {
			o.PrintNP.Println(v)
		}
	}

}

// fileOutput outputs to...well... a file
func fileOutput(lines []string, path string) (err error) {
	// get the output
	fileOutputString := "\n"
	for _, l := range lines {
		fileOutputString += l + "\n"
	}
	var outputbytes = []byte(fileOutputString + "\n")

	// open output file
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if closeError := fo.Close(); closeError != nil {
			if err == nil {
				err = closeError
			}
		}
	}()

	// write all lines
	if _, err := fo.Write(outputbytes); err != nil {
		return err
	}
	return nil
}

/*///////////////
  // Generator /
 ////////////*/

func generatePassphrases(s ewState) ([]string, error) {
	// setup the return value
	var outputList []string

	// sanity check the min and max values and wordcount
	if s.minimumChars > s.maximumChars && s.maximumChars != -1 {
		return nil, errors.New("bad character values - minimum was greater than maximum")
	}
	if s.wordCount < 3 {
		return nil, errors.New("will not generate passphrase with less than 3 words")
	}
	if s.capsWordIndex > s.wordCount {
		return nil, errors.New("requested capitalisation index is out of range")
	}

	o.OverShare.Println("Generating " + strconv.Itoa(s.wordCount) + " words per passphrase.")

	// make flow easier if competing flags were set
	if s.preventCaps {
		s.randomCaps = false
	}

	// loop for the number of passwords requested
	for index := 0; index < s.quantity; index++ {
		// generate the requested words
		o.OverShare.Println("Rolling!")
		outputString := ""
		// TODO: implement choice of which position is capsWordIndex
		if s.preventCaps || s.randomCaps {
			outputString = lookupAndConcatEffWords(s.wordCount, s.capsWordIndex, true, s.wordList)
		} else {
			outputString = lookupAndConcatEffWords(s.wordCount, s.capsWordIndex, false, s.wordList)
		}

		// truncate, keeping in mind possible next steps
		if s.maximumChars != -1 {
			slicePos := s.maximumChars
			// figure out how many extra chars to lose
			if !s.preventSpecials {
				slicePos--
			}
			if !s.preventInteger {
				slicePos--
			}

			// just do it already
			outputString = outputString[:slicePos]
		}

		// break dictionary based attacks using the wordlist
		// also comply with common password policies to make them more useful
		if s.randomCaps {
			// TODO: pick a char and uppercase it randomly
			outputString = capitaliseRandomPositionInString(outputString)
		}
		if !s.preventSpecials {
			// TODO: pick random special, insert randomly
			insertionChar, err := getRandomSpecialChar()
			if err != nil {
				return nil, err
			}
			outputString = insertCharAtRandomPos(outputString, insertionChar)
		}
		if !s.preventInteger {
			insertionInt := strconv.Itoa(randomIntSecure(10))
			outputString = insertCharAtRandomPos(outputString, insertionInt)
		}

		o.OverShare.Println("Generated: " + outputString)

		// check if a miniumum password length was requested
		if s.minimumChars != -1 {
			// if this one is too short, reset the run and don't add to the list
			if len(outputString) < s.minimumChars+1 {
				index--

				// sanity check
				err := oneStepCloserToTheEdge()
				if err != nil {
					return nil, err
				}

			} else { // met minimum spec, so add to the list
				outputList = append(outputList, outputString)
			}
		} else { // no min req, add to the list
			outputList = append(outputList, outputString)
		}
	}
	return outputList, nil
}

func lookupAndConcatEffWords(quantity int, upperCasedIndex int, ignoreUpperCaseByWordIndex bool, wordList map[int]string) (output string) {

	// if we are doing case by word index, ensure we know which one to capitalise
	if !ignoreUpperCaseByWordIndex {
		if upperCasedIndex == -1 {
			upperCasedIndex = randomIntSecure(100000) % quantity
		} else {
			upperCasedIndex--
		}
	}

	// loop through generating for the proscribed quantity
	for wordRolls := 0; wordRolls < quantity; wordRolls++ {
		key := rollFiveDice()
		if val, ok := wordList[key]; ok {
			o.OverShare.Println("Rolled [" + strconv.Itoa(key) + "]. Found " + val)

			// capitalise if it is so written
			if !ignoreUpperCaseByWordIndex {
				if upperCasedIndex == wordRolls {
					val = strings.ToUpper(val[:1]) + val[1:]
				}
			}
			output += val
		} else { // failed to find word, so lets just grab another
			wordRolls--
		}
	}
	return output
}

// oneStepCloserToTheEdge
func oneStepCloserToTheEdge() error {
	sanity++
	if sanity > sanityCheck {
		everythingYouSayToMe := []string{
			"SANITY CHECK!!",
			strconv.Itoa(sanityCheck) + " passwords have been generated and discarded due to length minimum.",
			"Please don't do this again?",
			"(try raising wordcount per passphrase)",
		}
		for _, message := range everythingYouSayToMe {
			o.Warn.Println(message)
		}
		return errors.New("sanity check detected possible infinite loop")
	}
	return nil
}

// insertCharAtRandomPos does what it says on the tin.
// except it never allows placement at the beginning or end to make it a little bit harder to crack the passwords
func insertCharAtRandomPos(input string, insertion string) (output string) {
	// pick a spot to break the string. Never allow placement in the first or last positions
	placeToInsert := randomIntSecure(len(input)-1) + 1

	// insert
	output = input[:placeToInsert] + insertion + input[placeToInsert:]
	return output
}

// rollFivedice will generate an int with 5 digits each with a value from 1-6
func rollFiveDice() int {

	var counter int
	index := 1

	// roll a die for each placeValue of a 5 digit number
	for index < 100000 {
		// generate a random int between 1 and 6 (a die)
		n := randomIntSecure(6) + 1

		// add the last die roll to the current placeValue and setup for the next run
		counter += n * index
		index = index * 10
	}

	return counter
}

// capitaliseRandomPositionInString; does what it says on the tin
func capitaliseRandomPositionInString(input string) string {
	// setup to single out a character
	inputBytes := []byte(input)
	pos := randomIntSecure(100000) % len(inputBytes)

	upperCharAtPos := strings.ToUpper(string(inputBytes[pos]))
	inputBytes[pos] = []byte(upperCharAtPos)[0]
	return string(inputBytes)
}

// randomIntSecure generates a random integer in the range [0, max)
// It uses the cryptographically secure random number generator of the OS
func randomIntSecure(max int) int {
	max64 := int64(max)
	random, err := rand.Int(rand.Reader, big.NewInt(max64))
	if err != nil {
		panic(err)
	}
	return int(random.Int64())
}

// getRandomSpecialChar picks a random special char as defined by OWASP
// https://www.owasp.org/index.php/Password_special_characters
func getRandomSpecialChar() (string, error) {

	// according to OWASP, these are all the special chars
	var specialChars = []string{
		" ",
		"!",
		`"`,
		"#",
		"$",
		"%",
		"&",
		`'`,
		"(",
		")",
		"*",
		"+",
		",",
		"-",
		".",
		"/",
		":",
		";",
		"<",
		"=",
		">",
		"?",
		"@",
		"[",
		`\`,
		"]",
		"^",
		"_",
		"`",
		"{",
		"|",
		"}",
		"~",
	}

	// get a random pos in the slice and return
	random := randomIntSecure(100000)
	return specialChars[random%len(specialChars)], nil
}
