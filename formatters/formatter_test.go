package formatters

import (
	"testing"
)

func TestTelegrafSummary(t *testing.T) {
	tests := []struct {
		name     string
		summary  TelegrafSummary
		expected map[string]string
	}{
		{
			name: "Install Test",
			summary: TelegrafSummary{
				ActionSummary: ActionSummary{
					Action:   "Install",
					StartCmd: "sudo service telegraf start",
					Config:   "/etc/telegraf/telegraf.conf",
				},
				Plugins: []string{"cpu", "mem", "disk"},
			},
			expected: map[string]string{
				"StartCmd": "sudo service telegraf start",
				"Config":   "/etc/telegraf/telegraf.conf",
				"Plugins":  "cpu, mem, disk",
			},
		},
		{
			name: "Update Api Key Test",
			summary: TelegrafSummary{
				ActionSummary: ActionSummary{
					Action:     "Update Api Key",
					RestartCmd: "sudo service telegraf restart",
					Config:     "/etc/telegraf/telegraf.conf",
				},
			},
			expected: map[string]string{
				"RestartCmd": "sudo service telegraf restart",
				"Config":     "/etc/telegraf/telegraf.conf",
			},
		},
		{
			name: "Uninstall Test",
			summary: TelegrafSummary{
				ActionSummary: ActionSummary{
					Action: "Uninstall",
				},
			},
			expected: map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.summary.GenerateContent()
			for key, expectedValue := range test.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected %s for key %s, got %s", expectedValue, key, result[key])
				}
			}
		})
	}

}

func TestOtelContribSummary(t *testing.T) {
	tests := []struct {
		name     string
		summary  OtelContribSummary
		expected map[string]string
	}{
		{
			name: "Install Test",
			summary: OtelContribSummary{
				ActionSummary: ActionSummary{
					Action:   "Install",
					StartCmd: "sudo service otelcontribcol start",
					Config:   "/etc/otelcontribcol/config.yaml",
				},
				Receiver: "hostmetrics",
				Exporter: "carbon",
			},
			expected: map[string]string{
				"StartCmd": "sudo service otelcontribcol start",
				"Config":   "/etc/otelcontribcol/config.yaml",
				"Receiver": "hostmetrics",
				"Exporter": "carbon",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.summary.GenerateContent()
			for key, expectedValue := range test.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected %s for key %s, got %s", expectedValue, key, result[key])
				}
			}
		})
	}

}
