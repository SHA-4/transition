package main

import (
        "encoding/csv"
        "fmt"
        "log"
        "os"
)

func main() {
        f, err := os.Open(".pronouns")
        if err != nil {
                log.Fatal("Error reading pronoun file", err)
        }
        defer f.Close()

        csvReader := csv.NewReader(f)
        pronouns, err := csvReader.ReadAll()
        if err != nil {
                log.Fatal("Error reading csv data", err)
        }

        fmt.Println("Data", pronouns)
}

