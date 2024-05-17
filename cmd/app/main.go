package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/alexey-savchuk/computer-club/internal/system"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	lines := []string{}
	for i := 0; i < 3; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read line: %s\n", err)
			os.Exit(1)
		}
		lines = append(lines, line)
	}

	preamble, err := system.ParsePreamble(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse preamble: %s\n", err)
	}

	club := system.NewComputerClubSystem(
		preamble.NumTables,
		preamble.OpenTime,
		preamble.CloseTime,
		preamble.Cost,
	)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read line: %s\n", err)
			os.Exit(1)
		}

		inEvent, err := system.ParseEvent(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse event: %s\n", err)
			os.Exit(1)
		}

		outEvents, err := club.Process(inEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to process event: %s\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(writer, inEvent)
		for _, outEvent := range outEvents {
			fmt.Fprintln(writer, outEvent)
		}
	}

	outEvents, _ := club.CloseClub()
	for _, outEvent := range outEvents {
		fmt.Fprintln(writer, outEvent)
	}

	profit := club.GetStatistics()
	for _, p := range profit {
		hours := int(p.TimeOccupied.Hours())
		minutes := int(p.TimeOccupied.Minutes()) % 60

		fmt.Fprintf(
			writer, "%d %d %02d:%02d\n",
			p.TableNum, p.Profit,
			hours, minutes,
		)
	}
}
