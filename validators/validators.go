// Package validators provides functions to validate if the rules of the `.editorconfig` are respected
package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker.go/utils"
)

// Indentation validates a files indentation
func Indentation(line string, indentStyle string, indentSize int) error {
	if indentStyle == "space" {
		return Space(line, indentSize)
	} else if indentStyle == "tab" {
		return Tab(line)
	}

	return nil
}

// Space validates if a line is indented correctly respecting the indentSize
func Space(line string, indentSize int) error {
	if len(line) > 0 && indentSize > 0 {
		// match recurring spaces indentSize times - this can be recurring or never
		// match either a space followed by a * and maybe a space (block-comments)
		// or match everything despite a space or tab-character
		regexpPattern := fmt.Sprintf("^( {%d})*( \\* ?|[^ \t])", indentSize)

		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if !matched {
			return fmt.Errorf("Wrong amount of left-padding spaces(want multiple of %d)", indentSize)
		}

	}

	return nil
}

// Tab validates if a line is indented with only tabs
func Tab(line string) error {
	if len(line) > 0 {
		regexpPattern := "^\t*[^ \t]"
		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if !matched {
			return errors.New("Wrong indentation type(spaces instead of tabs)")
		}

	}

	return nil
}

// TrailingWhitespace validates if a line has trailing whitespace
func TrailingWhitespace(line string, trimTrailingWhitespace bool) error {
	if trimTrailingWhitespace {
		regexpPattern := "^.*[ \t]+$"
		matched, err := regexp.MatchString(regexpPattern, line)

		if err != nil {
			panic(err)
		}

		if matched {
			return errors.New("Trailing whitespace")
		}
	}

	return nil
}

// FinalNewline validates if a file has a final and correct newline
func FinalNewline(fileContent string, insertFinalNewline bool, endOfLine string) error {
	if insertFinalNewline {
		regexpPattern := fmt.Sprintf("%s$", utils.GetEolChar(endOfLine))
		matched, err := regexp.MatchString(regexpPattern, fileContent)

		if err != nil {
			panic(err)
		}

		if !matched {
			return errors.New("Wrong line endings or new final newline")
		}
	}

	return nil
}

// LineEnding validates if a file uses the correct line endings
func LineEnding(fileContent string, endOfLine string) error {
	if endOfLine != "" {
		expectedEolChar := utils.GetEolChar(endOfLine)
		expectedEols := len(strings.Split(fileContent, expectedEolChar))
		lfEols := len(strings.Split(fileContent, "\n"))
		crEols := len(strings.Split(fileContent, "\r"))
		crlfEols := len(strings.Split(fileContent, "\r\n"))

		switch endOfLine {
		case "lf":
			if !(expectedEols == lfEols && crEols == 1 && crlfEols == 1) {
				return errors.New("Not all lines have the correct end of line character")
			}
		case "cr":
			if !(expectedEols == crEols && lfEols == 1 && crlfEols == 1) {
				return errors.New("Not all lines have the correct end of line character")
			}
		case "crlf":
			// A bit hacky because \r\n matches \r and \n
			if !(expectedEols == crlfEols && lfEols == expectedEols && crEols == expectedEols) {
				return errors.New("Not all lines have the correct end of line character")
			}
		}
	}

	return nil
}
