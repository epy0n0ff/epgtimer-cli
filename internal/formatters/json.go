package formatters

import (
	"encoding/json"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// JSONFormatter formats rules as JSON
type JSONFormatter struct{}

// Format converts rules to JSON format
func (j *JSONFormatter) Format(rules []models.AutoAddRule) (string, error) {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
