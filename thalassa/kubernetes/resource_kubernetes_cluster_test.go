package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
)

// defaultAutoscalerConfig returns the default autoscaler configuration
func defaultAutoscalerConfig() *kubernetes.AutoscalerConfig {
	return &kubernetes.AutoscalerConfig{
		ScaleDownDisabled:             false,
		ScaleDownDelayAfterAdd:        "10m",
		Estimator:                     "binpacking",
		Expander:                      "binpacking",
		IgnoreDaemonsetsUtilization:   false,
		BalanceSimilarNodeGroups:      false,
		ExpendablePodsPriorityCutoff:  -10,
		ScaleDownUnneededTime:         "10m",
		ScaleDownUtilizationThreshold: 0.5,
		MaxGracefulTerminationSec:     600,
		EnableProactiveScaleUp:        false,
	}
}

func TestConvertAutoscalerConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   interface{}
		expected *kubernetes.AutoscalerConfig
	}{
		{
			name:     "nil config returns default config",
			config:   nil,
			expected: defaultAutoscalerConfig(),
		},
		{
			name:     "empty list returns default config",
			config:   []interface{}{},
			expected: defaultAutoscalerConfig(),
		},
		{
			name:     "nil first element returns default config",
			config:   []interface{}{nil},
			expected: defaultAutoscalerConfig(),
		},
		{
			name:     "invalid type (not a list) returns default config",
			config:   "not a list",
			expected: defaultAutoscalerConfig(),
		},
		{
			name:     "invalid first element (not a map) returns default config",
			config:   []interface{}{"not a map"},
			expected: defaultAutoscalerConfig(),
		},
		{
			name:     "empty map returns default config",
			config:   []interface{}{map[string]interface{}{}},
			expected: defaultAutoscalerConfig(),
		},
		{
			name: "all fields set with custom values",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_disabled":              true,
					"scale_down_delay_after_add":       "15m",
					"estimator":                        "least-waste",
					"expander":                         "priority",
					"ignore_daemonsets_utilization":    true,
					"balance_similar_node_groups":      true,
					"expendable_pods_priority_cutoff":  -5,
					"scale_down_unneeded_time":         "20m",
					"scale_down_utilization_threshold": 0.6,
					"max_graceful_termination_sec":     900,
					"enable_proactive_scale_up":        true,
				},
			},
			expected: &kubernetes.AutoscalerConfig{
				ScaleDownDisabled:             true,
				ScaleDownDelayAfterAdd:        "15m",
				Estimator:                     "least-waste",
				Expander:                      "priority",
				IgnoreDaemonsetsUtilization:   true,
				BalanceSimilarNodeGroups:      true,
				ExpendablePodsPriorityCutoff:  -5,
				ScaleDownUnneededTime:         "20m",
				ScaleDownUtilizationThreshold: 0.6,
				MaxGracefulTerminationSec:     900,
				EnableProactiveScaleUp:        true,
			},
		},
		{
			name: "partial fields set, others use defaults",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_disabled":              true,
					"estimator":                        "most-pods",
					"scale_down_utilization_threshold": 0.75,
				},
			},
			expected: func() *kubernetes.AutoscalerConfig {
				cfg := defaultAutoscalerConfig()
				cfg.ScaleDownDisabled = true
				cfg.Estimator = "most-pods"
				cfg.ScaleDownUtilizationThreshold = 0.75
				return cfg
			}(),
		},
		{
			name: "fields with nil values use defaults",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_disabled":              nil,
					"scale_down_delay_after_add":       nil,
					"estimator":                        nil,
					"expander":                         nil,
					"ignore_daemonsets_utilization":    nil,
					"balance_similar_node_groups":      nil,
					"expendable_pods_priority_cutoff":  nil,
					"scale_down_unneeded_time":         nil,
					"scale_down_utilization_threshold": nil,
					"max_graceful_termination_sec":     nil,
					"enable_proactive_scale_up":        nil,
				},
			},
			expected: defaultAutoscalerConfig(),
		},
		{
			name: "boolean fields with false values",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_disabled":           false,
					"ignore_daemonsets_utilization": false,
					"balance_similar_node_groups":   false,
					"enable_proactive_scale_up":     false,
				},
			},
			expected: defaultAutoscalerConfig(),
		},
		{
			name: "integer fields with zero and negative values",
			config: []interface{}{
				map[string]interface{}{
					"expendable_pods_priority_cutoff": 0,
					"max_graceful_termination_sec":    0,
				},
			},
			expected: func() *kubernetes.AutoscalerConfig {
				cfg := defaultAutoscalerConfig()
				cfg.ExpendablePodsPriorityCutoff = 0
				cfg.MaxGracefulTerminationSec = 0
				return cfg
			}(),
		},
		{
			name: "float field with various threshold values",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_utilization_threshold": 0.3,
				},
			},
			expected: func() *kubernetes.AutoscalerConfig {
				cfg := defaultAutoscalerConfig()
				cfg.ScaleDownUtilizationThreshold = 0.3
				return cfg
			}(),
		},
		{
			name: "string fields with custom time and estimator values",
			config: []interface{}{
				map[string]interface{}{
					"scale_down_delay_after_add": "30m",
					"scale_down_unneeded_time":   "1h",
					"estimator":                  "least-waste",
					"expander":                   "random",
				},
			},
			expected: func() *kubernetes.AutoscalerConfig {
				cfg := defaultAutoscalerConfig()
				cfg.ScaleDownDelayAfterAdd = "30m"
				cfg.ScaleDownUnneededTime = "1h"
				cfg.Estimator = "least-waste"
				cfg.Expander = "random"
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertAutoscalerConfig(tt.config)

			assert.Equal(t, tt.expected.ScaleDownDisabled, result.ScaleDownDisabled, "ScaleDownDisabled should match")
			assert.Equal(t, tt.expected.ScaleDownDelayAfterAdd, result.ScaleDownDelayAfterAdd, "ScaleDownDelayAfterAdd should match")
			assert.Equal(t, tt.expected.Estimator, result.Estimator, "Estimator should match")
			assert.Equal(t, tt.expected.Expander, result.Expander, "Expander should match")
			assert.Equal(t, tt.expected.IgnoreDaemonsetsUtilization, result.IgnoreDaemonsetsUtilization, "IgnoreDaemonsetsUtilization should match")
			assert.Equal(t, tt.expected.BalanceSimilarNodeGroups, result.BalanceSimilarNodeGroups, "BalanceSimilarNodeGroups should match")
			assert.Equal(t, tt.expected.ExpendablePodsPriorityCutoff, result.ExpendablePodsPriorityCutoff, "ExpendablePodsPriorityCutoff should match")
			assert.Equal(t, tt.expected.ScaleDownUnneededTime, result.ScaleDownUnneededTime, "ScaleDownUnneededTime should match")
			assert.Equal(t, tt.expected.ScaleDownUtilizationThreshold, result.ScaleDownUtilizationThreshold, "ScaleDownUtilizationThreshold should match")
			assert.Equal(t, tt.expected.MaxGracefulTerminationSec, result.MaxGracefulTerminationSec, "MaxGracefulTerminationSec should match")
			assert.Equal(t, tt.expected.EnableProactiveScaleUp, result.EnableProactiveScaleUp, "EnableProactiveScaleUp should match")
		})
	}
}
