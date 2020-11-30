package main

import (
        "bufio"
        "encoding/csv"
        "fmt"
        "golang.org/x/sys/unix"
        "log"
        "math/rand"
        "os"
        "strings"
        "time"
)

const defaultRuns = 10

func main() {
        f, err := os.Open(".pronouns")
        if err != nil {
                log.Fatal("Error reading pronoun file", err)
        }
        defer f.Close()

        csvReader := csv.NewReader(f)
        // pronouns file format is NAME, SUBJECT PRONOUN, OBJECT PRONOUN
        // e.g. Mark he him
        pronouns, err := csvReader.ReadAll()
        if err != nil {
                log.Fatal("Error reading csv data", err)
        }

        runPractice(pronouns)
}

func runPractice(pronouns [][]string) {
        resetToMiddle()

        lastIsCorrect := true
        isExiting := false
        numCorrect := 0
        start := time.Now()
        var randomPronoun []string
        for i := 0; i < defaultRuns; i++ {
                if lastIsCorrect {
                        randomPronoun = pronouns[rand.Intn(len(pronouns))]
                }
                isExiting, lastIsCorrect = runTest(randomPronoun)
                resetToMiddle()

                if isExiting {
                        break
                }

                if lastIsCorrect {
                        numCorrect += 1
                }
        }
        end := time.Now()
        elapsed := end.Sub(start)
        printStats(elapsed, numCorrect, defaultRuns)
}

func printStats(elapsed time.Duration, numCorrect int, runCount int) {
        elapsedSeconds := elapsed.Seconds()
        fmt.Printf("Time: %.2fs\n", elapsedSeconds)
        fmt.Printf("With %d/%d (%.2f%%) correct\n", numCorrect, runCount, float64(100 * numCorrect) / float64(runCount))

        if numCorrect != runCount {
                adjustAmount := float64(5 * (runCount - numCorrect))
                fmt.Printf("Adjusted Time: %.2fs\n", adjustAmount + elapsedSeconds)
        }
}

func getWindowSize() (*unix.Winsize, error) {
	size, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	}

	return size, nil
}

// See ANSI escape codes here
// https://en.wikipedia.org/wiki/ANSI_escape_code
func runCSI(command string) {
        fmt.Printf("\033[" + command)
}

func clearScreen() {
        runCSI("2J")
}

func resetScreen() {
        clearScreen()
        runCSI("1;1H")
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func resetToMiddle () {
        size, err := getWindowSize()

        if err != nil {
                resetScreen()
                return
        }

        clearScreen()
        // if less than 10, no point in vertical centering
        if size.Row > 10 {
                command := fmt.Sprintf("%d;1H", (size.Row - 10) / 2)
                runCSI(command)
        }
}

// returns (exit status, is correct or not)
func runTest(pronoun []string) (bool, bool) {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println(pronoun[0])
        fmt.Println()
        text, _ := reader.ReadString('\n')
        formattedText := strings.ToLower(strings.TrimSuffix(text, "\n"))

        if formattedText == "exit" {
                return true, true
        }

        isCorrect := strings.ToLower(pronoun[1]) == formattedText

        if !isCorrect {
                showIncorrect(pronoun)
                reader.ReadString('\n')
        }

        return false, isCorrect
}

func printInCenter(text string) {
        size, err := getWindowSize()
        textLength := len(text)
        if err != nil {
                fmt.Println(text)
        }

        var halfCol int
        if int(size.Col) > textLength {
                halfCol = int(float64(int(size.Col) - textLength) / float64(2))
        } else {
                halfCol = 0
        }
        startPadding := strings.Repeat(" ", halfCol)
        // we want to add end padding so background highlighting works
        endPadding := strings.Repeat(" ", int(size.Col) - halfCol - textLength)
        fmt.Printf("%s%s%s\n", startPadding, text, endPadding)
}

func getColorCommand(code string) (string) {
        return fmt.Sprintf("\u001B[%sm", code)
}

func setColor(code string) {
        fmt.Printf(getColorCommand(code))
}

func resetColor() {
        setColor("0")
}

func getColoredTextString(text string, color string) (string) {
        return fmt.Sprintf("%s%s%s", getColorCommand(color), text, getColorCommand("0"))
}

func showIncorrect(pronoun []string) {
        fmt.Println()
        // red text
        setColor("1;33")
        // grey background
        setColor("41")
        printInCenter("Incorrect")
        resetColor()
        fmt.Println(pronoun[0] + "'s correct pronouns are " + getColoredTextString(pronoun[1], "36") + " and " + getColoredTextString(pronoun[2], "36"))
        fmt.Print("\nPress any key to continue")
}

