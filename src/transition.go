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

        resetToMiddle()
        lastIsCorrect := true
        var randomPronoun []string
        for {
                if lastIsCorrect {
                        randomPronoun = pronouns[rand.Intn(len(pronouns))]
                }
                lastIsCorrect = runTest(randomPronoun)
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

func runTest(pronoun []string) (bool) {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println(pronoun[0])
        fmt.Println()
        text, _ := reader.ReadString('\n')
        formattedText := strings.ToLower(strings.TrimSuffix(text, "\n"))

        isCorrect := strings.ToLower(pronoun[1]) == formattedText

        if !isCorrect {
                showIncorrect(pronoun)
                reader.ReadString('\n')
        }
        resetToMiddle()

        return isCorrect
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

