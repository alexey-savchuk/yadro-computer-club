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

	err := doWork(reader, writer)
	writer.Flush()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func doWork(reader *bufio.Reader, writer *bufio.Writer) error {
	lines := []string{}
	for i := 0; i < 3; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read line: %w", err)
		}
		lines = append(lines, line)
	}

	preamble, err := system.ParsePreamble(lines)
	if err != nil {
		return fmt.Errorf("failed to parse preamble: %w", err)
	}

	club := system.NewComputerClubSystem(
		preamble.NumTables,
		preamble.OpenTime,
		preamble.CloseTime,
		preamble.Cost,
	)

	fmt.Fprintln(writer, preamble.OpenTime.Format("15:04"))
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read line: %w", err)
		}

		inEvent, err := system.ParseEvent(line)
		if err != nil {
			return fmt.Errorf("failed to parse event: %w", err)
		}

		outEvents, err := club.Process(inEvent)
		if err != nil {
			return fmt.Errorf("failed to process event: %w", err)
		}

		for _, outEvent := range outEvents {
			fmt.Fprintln(writer, outEvent)
		}
	}

	if !club.IsClubClose() {
		outEvents, err := club.CloseClub()
		if err != nil {
			return fmt.Errorf("failed to close club: %w", err)
		}
		for _, outEvent := range outEvents {
			fmt.Fprintln(writer, outEvent)
		}
	}
	fmt.Fprintln(writer, preamble.CloseTime.Format("15:04"))

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

	return nil
}
