package collectors

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// SampleData represents the parsed information for each label
type Stat struct {
	Syscall                   string
	NumSamples                int
	Unit                      string
	Min, Max, Sum, SumSquared int
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseStats(input string) ([]Stat, error) {
	var result []Stat

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "snapshot_time") {
			// Skip the snapshot_time line
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		syscall := fields[0]
		
		numSamples, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of samples: %v", err)
		}

		unit := fields[3][1 : len(fields[3])-1] // Extracting unit from [usecs]
		min, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, fmt.Errorf("failed to parse min value: %v", err)
		}

		max, err := strconv.Atoi(fields[5])
		if err != nil {
			return nil, fmt.Errorf("failed to parse max value: %v", err)
		}

		sum, err := strconv.Atoi(fields[6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum value: %v", err)
		}

		sumSquared, err := strconv.Atoi(fields[7])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum squared value: %v", err)
		}

		result = append(result, Stat{
			Syscall:    syscall,
			NumSamples: numSamples,
			Unit:       unit,
			Min:        min,
			Max:        max,
			Sum:        sum,
			SumSquared: sumSquared,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}
