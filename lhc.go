/*
SPDX-License-Identifier: MIT

MIT License

Copyright (c) 2017 Thanh Ha

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var VERSION = "0.1.0"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func fetchLicense(filename string) string {
    file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

    code, commentSection := false, false
	licenseText := ""
	scanner := bufio.NewScanner(file)

    b1 := make([]byte, 2)
    n1, err := file.Read(b1)
    check(err)
    if isComment(string(b1)) {
        fmt.Printf("Comment string detected. %d bytes: %s\n", n1, string(b1))
        code = true
    }
    file.Seek(0, 0)  // Reset so we can read the full file

    i := 0
	for scanner.Scan() {
		s := scanner.Text()

        if strings.Contains(s, "Copyright") {
            continue
        } else if strings.Contains(s, "SPDX-License-Identifier") {
            continue
        }

        if code == true {
    		if strings.HasPrefix(s, "/*") {
    			commentSection = true
    		} else if commentSection && strings.Contains(s, "*/") {
    			commentSection = false
    		}

    		if !commentSection &&
    			!isComment(s) {
    			break
    		}

    		s = strings.TrimPrefix(s, "#")
    		s = strings.TrimPrefix(s, "//")
    		s = strings.TrimPrefix(s, "/*")
    		s = strings.Split(s, "*/")[0]
        }

		licenseText += s

        // Limit to reading only the first few lines to not read entire code file
        i++
        if i > 100 {
            break
        }
	}

    return stripSpaces(licenseText)
}

func isComment(str string) bool {
    if !strings.HasPrefix(str, "#") &&
        !strings.HasPrefix(str, "//") &&
        !strings.HasPrefix(str, "/*") {
        return false
    }

    return true
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] [PATTERN]...\n", os.Args[0])
	fmt.Printf("Scans a directory for files matching PATTERN and compares them with an expected license header.\n")
	fmt.Printf("\nPATTERN is a space separated list of regex patterns to search for files.\n")
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	licensePtr := flag.String("license", "license.txt", "Comma-separated list of license files to compare against.")
	versionPtr := flag.Bool("version", false, "Print version")
	// directoryPtr := flag.String("directory", ".", "Directory to search for files.")

	flag.Usage = usage
	flag.Parse()

	if *versionPtr {
		fmt.Println("License Checker version", VERSION)
		os.Exit(0)
	}

	fmt.Println("Search Patterns:", flag.Args())

	licenseText := fetchLicense(*licensePtr)
	fmt.Println("License Text")
	fmt.Println(licenseText)

	headerText := fetchLicense("lhc.go")
	fmt.Println("Header Text")
	fmt.Println(headerText)

	if licenseText != headerText {
		fmt.Println("WARNING: License header does not match.", "lhc.go")
	}
}
