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

        resetScreen()
        for {
                runTest(pronouns)
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

func runTest(pronouns [][]string) {
        randomPronoun := pronouns[rand.Intn(len(pronouns))]
        reader := bufio.NewReader(os.Stdin)
        fmt.Println(randomPronoun[0])
        text, _ := reader.ReadString('\n')
        formattedText := strings.ToLower(strings.TrimSuffix(text, "\n"))

        isCorrect := strings.ToLower(randomPronoun[1]) == formattedText

        if !isCorrect {
                showIncorrect(randomPronoun)
                reader.ReadString('\n')
        }
        fmt.Print("\n")
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

func setColor(code string) {
        fmt.Printf("\u001B[%sm", code)
}

func resetColor() {
        setColor("0")
}

func showIncorrect(pronoun []string) {
        // red text
        setColor("31")
        // grey background
        setColor("47")
        printInCenter("Incorrect")
        resetColor()
        fmt.Println(pronoun[0] + "'s correct pronouns are " + pronoun[1] + " and " + pronoun[2])
        fmt.Print("\nPress any key to continue")
}

