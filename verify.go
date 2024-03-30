// package mytest provides testing utilities for long, repetitive, table driven tests.
package mytest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// max number of characters to display, before and after first discrepency detected.
// adjusting this value has no effect on stored files, only on displayed messages.
var DISPLAY_WINDOW = 160

// Verify provided content against reference file.
// If no reference file found, create it.
// A .want extension and a _ prefix are added to the reference filename base name, and path is removed to force storage in local source folder.
// If content differs from existing reference file, create a xxx.got file for further review and fail the test.
func Verify(t *testing.T, content string, reference string) {

	content = fmt.Sprintf("Test name : %s\nThis file : %s\n%s\n", t.Name(), filepath.Base(wantFile(reference)), content)

	fmt.Println(GREEN+"Verifying test results against file : "+RESET, wantFile(reference))
	check, err := os.ReadFile(wantFile(reference))
	if err != nil {
		fmt.Println(RED + "File not found, create it as a reference for future test. Make sure you manually review it !" + RESET)
		os.WriteFile(wantFile(reference), []byte(content), 0644)
		return
	}
	sc := string(check)
	if sc != content { // we know it is different, lets try to show where ...
		for i, c := range content {
			if i >= len([]rune(sc)) || c != ([]rune(sc)[i]) {
				i1 := i - DISPLAY_WINDOW
				i2 := i + DISPLAY_WINDOW
				if i1 <= 0 {
					i1 = 0
				}
				if i2 > len(content) {
					i2 = len(content)
				}

				fmt.Printf("\n===============================================================\n%s : Results differ from reference file", t.Name())
				fmt.Printf("\n============================ got ==============================\n%s%s%s%s\n",
					content[i1:i], RED, content[i:i2], RESET)
				if i2 >= len([]rune(sc)) {
					i2 = len([]rune(sc))
				}
				fmt.Printf("\n============================ want==============================\n%s%s%s%s\n",
					sc[i1:i], RED, sc[i:i2], RESET)

				_ = os.WriteFile(gotFile(reference), []byte(content), 0644)
				fmt.Printf("\n%s*** FAILING *** Got file saved in : %s%s\n", RED, gotFile(reference), RESET)
				t.Fatalf("Result differs from reference file in %s", t.Name())
			}
		}

		// If we reach here, it means content is matching ref, but both files are different.
		// There must be extra reference ? Let's show it.
		i := len(content)
		i1 := len(content) - DISPLAY_WINDOW
		if i1 < 0 {
			i1 = 0
		}
		i2 := len(content) + 160
		if i2 > len(sc) {
			i2 = len(sc)
		}
		fmt.Printf("\n===============================================================\n%s : Results differ from reference file", t.Name())
		fmt.Printf("\n============================ got ==============================\n%s\n",
			content[i1:i])
		fmt.Println()
		fmt.Printf("\n============================ want==============================\n%s%s%s%s\n",
			sc[i1:i], RED, sc[i:i2], RESET)

		_ = os.WriteFile(gotFile(reference), []byte(content), 0644)
		fmt.Printf("\n%s*** FAILING *** Got file saved in : %s%s\n", RED, gotFile(reference), RESET)
		t.Fatalf("Result differs from reference file in %s", wantFile(reference))

	}
}

// Apply Verify to the content of the specified file.
// Not suitable for large files, because file will be loaded entirely in memory.
func VerifyFile(t *testing.T, pathToFile string) {
	f, err := os.Open(pathToFile)
	if err != nil {
		t.Fatalf("%sError accessing file %s %s: %v", RED, pathToFile, RESET, err)
	}
	defer f.Close()
	bb, err := io.ReadAll(f)
	patt := regexp.MustCompile("[^A-Za-z0-9]")
	reference := patt.ReplaceAllString(pathToFile, ".")
	Verify(t, string(bb), reference)
}

// construct the reference file name
func wantFile(filename string) string {
	filename = filepath.Base(filename)
	filename, _ = filepath.Abs("_" + filename)
	return filename + ".want"
}

// construct the result (.got) file name if different from reference
func gotFile(filename string) string {
	filename = filepath.Base(filename)
	filename, _ = filepath.Abs("_" + filename)
	return filename + ".got"
}
