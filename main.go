package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file")
	timeLimit := flag.Int("limit", 2, "time limit for quiz")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided CSV file")
	}

	score := 0
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemsLoop:
	for _, problem := range problems {

		fmt.Printf("%s=", problem.q)
		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemsLoop
		case answer := <-answerChannel:
			if answer == problem.a {
				score++
			}
		}
	}

	fmt.Printf("Score is %d out of %d\n", score, len(lines))
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}
