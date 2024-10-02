package bufioiter_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jinford/bufioiter"
)

func ExampleNewScanner() {
	const input = `golang
python
java
`
	for text, err := range bufioiter.NewScanner(strings.NewReader(input)) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		fmt.Println(text)
	}

	// Output:
	// golang
	// python
	// java
}

// Use a Scanner to implement a simple word-count utility by scanning the
// input as a sequence of space-delimited tokens.
func ExampleNewScanner_words() {
	// An artificial input source.
	const input = "Now is the winter of our discontent,\nMade glorious summer by this sun of York.\n"

	scanner := bufioiter.NewScanner(strings.NewReader(input),
		// Set the split function for the scanning operation.
		bufioiter.Split(bufio.ScanWords),
	)

	// Count the words.
	count := 0
	for _, err := range scanner {
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
		}

		count++
	}

	fmt.Printf("%d\n", count)

	// Output: 15
}

// Use a Scanner with a custom split function (built by wrapping ScanWords) to validate
// 32-bit decimal input.
func ExampleNewScanner_custom() {
	// An artificial input source.
	const input = "1234 5678 1234567901234567890"

	scanner := bufioiter.NewScanner(strings.NewReader(input),
		// Create a custom split function by wrapping the existing ScanWords function.
		bufioiter.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			advance, token, err = bufio.ScanWords(data, atEOF)
			if err == nil && token != nil {
				_, err = strconv.ParseInt(string(token), 10, 32)
			}
			return
		}),
	)

	for text, err := range scanner {
		if err != nil {
			fmt.Printf("Invalid input: %s", err)
		}

		fmt.Println(text)
	}

	// Output:
	// 1234
	// 5678
	// Invalid input: strconv.ParseInt: parsing "1234567901234567890": value out of range
}

// Use a Scanner with a custom split function to parse a comma-separated
// list with an empty final value.
func ExampleNewScanner_emptyFinalToken() {
	// Comma-separated list; last entry is empty.
	const input = "1,2,3,4,"

	scanner := bufioiter.NewScanner(strings.NewReader(input),
		bufioiter.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			for i := 0; i < len(data); i++ {
				if data[i] == ',' {
					return i + 1, data[:i], nil
				}
			}
			if !atEOF {
				return 0, nil, nil
			}
			// There is one final token to be delivered, which may be the empty string.
			// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
			// but does not trigger an error to be returned from Scan itself.
			return 0, data, bufio.ErrFinalToken
		}),
	)

	for text, err := range scanner {
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
		}

		fmt.Printf("%q ", text)
	}

	// Output: "1" "2" "3" "4" ""
}

func ExampleNewScanner_earlyStop() {
	const input = "1,2,STOP,4,"

	scanner := bufioiter.NewScanner(strings.NewReader(input),
		bufioiter.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			i := bytes.IndexByte(data, ',')
			if i == -1 {
				if !atEOF {
					return 0, nil, nil
				}
				// If we have reached the end, return the last token.
				return 0, data, bufio.ErrFinalToken
			}
			// If the token is "STOP", stop the scanning and ignore the rest.
			if string(data[:i]) == "STOP" {
				return i + 1, nil, bufio.ErrFinalToken
			}
			// Otherwise, return the token before the comma.
			return i + 1, data[:i], nil
		}),
	)

	for text, err := range scanner {
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
		}

		fmt.Printf("Got a token %q\n", text)
	}

	// Output:
	// Got a token "1"
	// Got a token "2"
}
