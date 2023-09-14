// package mytest provides testing utilities for long, repetitive, table driven tests.
package mytest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// max number of characters to display, before and after first discrepency detected.
var DISPLAY_WINDOW = 160

// Verify provided content against reference file.
// If no reference file found, create it.
// A .want extension and a _ prefix are added to the filename.
// If content differs from existing reference file, create a xxx.got file for further review and fail the test.
func Verify(t *testing.T, content string, filename string) {

	content = fmt.Sprintf("Test name : %s\nThis file : %s\n%s\n", t.Name(), filepath.Base(wantFile(filename)), content)

	fmt.Println(GREEN+"Verifying test results against file : "+RESET, wantFile(filename))
	check, err := os.ReadFile(wantFile(filename))
	if err != nil {
		fmt.Println(RED + "File not found, create it as a reference for future test. Make sure you manually review it !" + RESET)
		os.WriteFile(wantFile(filename), []byte(content), 0644)
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

				_ = os.WriteFile(gotFile(filename), []byte(content), 0644)
				fmt.Printf("\n%s*** FAILING *** Got file saved in : %s%s\n", RED, gotFile(filename), RESET)
				t.Fatalf("Result differs from reference file in %s", t.Name())
			}
		}

		// If we reach here, it means conet is matching ref, but both files are different.
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

		_ = os.WriteFile(gotFile(filename), []byte(content), 0644)
		fmt.Printf("\n%s*** FAILING *** Got file saved in : %s%s\n", RED, gotFile(filename), RESET)
		t.Fatalf("Result differs from reference file in %s", wantFile(filename))

	}
}

func wantFile(filename string) string {
	filename = filepath.Base(filename)
	filename, _ = filepath.Abs("_" + filename)
	return filename + ".want"
}

func gotFile(filename string) string {
	filename = filepath.Base(filename)
	filename, _ = filepath.Abs("_" + filename)
	return filename + ".got"
}
