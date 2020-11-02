package main

import (
        "bufio"
        "encoding/csv"
        "fmt"
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

        for {
                runTest(pronouns)
        }
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func runTest(pronouns [][]string) {
        randomPronoun := pronouns[rand.Intn(len(pronouns))]
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(randomPronoun[0] + "\n")
        text, _ := reader.ReadString('\n')
        formattedText := strings.ToLower(strings.TrimSuffix(text, "\n"))

        fmt.Println(strings.ToLower(randomPronoun[1]) == formattedText, "\n")
}
