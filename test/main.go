package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/txtar"
)

func main() {
	err := filepath.Walk("/app/test/suite", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while walking files: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".txtar") {
			return nil
		}

		testName := path[15:]
		details = ""

		testData, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read '%s': %w", path, err)
		}

		testData = []byte(strings.ReplaceAll(string(testData), "\r\n", "\n"))

		testArchive := txtar.Parse(testData)
		if err != nil {
			return err
		}

		var testCase TestCase

		// fmt.Println(string(testArchive.Comment))
		for _, file := range testArchive.Files {
			// fmt.Println(string(file.Name))
			switch file.Name {
			case "compiler.exitstatus":
				testCase.compilerExitCodeGiven = true
				for _, code := range strings.Split(string(file.Data), ",") {

					parsedCode, err := strconv.Atoi(strings.TrimSpace(code))
					if err != nil {
						return fmt.Errorf("error while parsing expected compiler exit code in test case %s: %w", path, err)
					}

					testCase.compilerExitCodes = append(testCase.compilerExitCodes, parsedCode)
				}
			case "src.lx":
				testCase.srcGiven = true
				testCase.src = string(file.Data)
			case "program.stdout":
				testCase.programStdoutGiven = true
				testCase.programStdout = strings.TrimSpace(string(file.Data))
			default:
				return fmt.Errorf("unexpected section '%s' in test case %s", file.Name, path)
			}
		}

		if !testCase.srcGiven {
			return fmt.Errorf("no source code defined in test case %s", path)
		}

		if !testCase.compilerExitCodeGiven {
			return fmt.Errorf("no expected compiler exit code defined in test case %s", path)
		}

		src, err := os.Create("./src.lx")
		if err != nil {
			return fmt.Errorf("error while creating or truncating src file: %w", err)
		}

		_, err = src.WriteString(testCase.src)
		if err != nil {
			return fmt.Errorf("error while writing to src file: %w", err)
		}

		details += fmt.Sprintf("\nSOURCE:\n%s", testCase.src)

		laxc := exec.Command("./laxc.exe", "src.lx")

		laxcOutput, err := laxc.CombinedOutput()
		if err != nil && laxc.ProcessState == nil {
			return fmt.Errorf("error while starting laxc in '%s': %w", testName, err)
		}

		if err != nil {
			details += fmt.Sprintf("\nCOMPILER ERROR:\n%s\n", err.Error())

			if strings.Contains(err.Error(), "signal") {
				gdb := exec.Command("gdb", "-q", "laxc.exe")
				gdb.Stdin = bytes.NewBufferString(`
run src.lx
where
q
				`)

				debugInfo, err := gdb.CombinedOutput()
				if err != nil {
					return fmt.Errorf("error while debugging compiler segfault: %w", err)
				}

				details += fmt.Sprintf("\nSTACK TRACE:\n%s\n", debugInfo)
			}
		}

		details += fmt.Sprintf("\nCOMPILER OUTPUT:\n%s", laxcOutput)

		compilerExitCode := laxc.ProcessState.ExitCode()

		if !slices.Contains(testCase.compilerExitCodes, compilerExitCode) {
			fail(
				testName,
				fmt.Sprintf("compiler exited with status code %d, expected %v", compilerExitCode, testCase.compilerExitCodes),
			)

			return nil
		}

		if compilerExitCode != 0 {
			succeed(testName)

			return nil
		}

		ilFile, err := os.Open("./src.lx.il")
		if err != nil {
			return fmt.Errorf("error while opening il-file: %w", err)
		}

		ilData, err := io.ReadAll(ilFile)
		if err != nil {
			return fmt.Errorf("error while reading il-file: %w", err)
		}

		details += fmt.Sprintf("\nIL:%s", string(ilData))

		asmFile, err := os.Open("./src.lx.s")
		if err != nil {
			return fmt.Errorf("error while opening assembly-file: %w", err)
		}

		asmData, err := io.ReadAll(asmFile)
		if err != nil {
			return fmt.Errorf("error while reading il-file: %w", err)
		}

		details += fmt.Sprintf("\nASSEMBLY:%s", string(asmData))

		spim := exec.Command("/app/spim/spim/spim", "-f", "./src.lx.s")

		spimStdout, err := spim.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error while running spim: %w", err)
		}

		spimExitCode := spim.ProcessState.ExitCode()
		if spimExitCode != 0 {
			return fmt.Errorf("spim exited with exit code %d", spimExitCode)
		}

		cleanedSpimStdout := strings.ReplaceAll(string(spimStdout), "Loaded: /usr/share/spim/exceptions.s", "")
		cleanedSpimStdout = strings.TrimSpace(cleanedSpimStdout)

		matches, err := regexp.Match("^"+testCase.programStdout+"$", []byte(cleanedSpimStdout))
		if spimExitCode != 0 {
			return fmt.Errorf("error while matching program stdout: %w", err)
		}

		if !matches {
			fail(
				testName,
				fmt.Sprintf("Actual program output ('%s') did not match expected output ('%s')", cleanedSpimStdout, testCase.programStdout),
			)

			return nil
		}

		succeed(testName)

		return nil
	})

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	if failed != 0 {
		fmt.Printf("FAIL: %d/%d test cases failed", failed, total)

		os.Exit(1)
	} else {
		fmt.Printf("SUCCESS: %d test cases executed successfully", total)
	}
}

type TestCase struct {
	srcGiven              bool
	src                   string
	compilerExitCodeGiven bool
	compilerExitCodes     []int
	programStdoutGiven    bool
	programStdout         string
}

var (
	total   uint
	failed  uint
	details string
)

func succeed(testName string) {
	fmt.Printf("RUN '%s' SUCCEEDED\n", testName)
	total++
}

func fail(testName string, reason string) {
	fmt.Printf("RUN '%s' FAILED: %s\n", testName, reason)
	total++
	failed++
	fmt.Printf("%s\n", details)
}
