package rebalancer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func processSSELine(line string) (interface{}, error) {
	if strings.HasPrefix(line, "data:") {
		// Extract the JSON part of the line (strip "data:")
		jsonData := strings.TrimPrefix(line, "data:")
		jsonData = strings.TrimSpace(jsonData) // Remove any extra spaces or newlines

		var payload Response
		err := json.Unmarshal([]byte(jsonData), &payload)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON payload: %v", err)
		}

		return payload, nil
	} else {
		logrus.Info(fmt.Sprintf("Ignoring unexpected line: %s\n", line))
		return nil, nil
	}
}
