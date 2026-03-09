package cmd

import "testing"

// TestIsHighRisk demonstrates standard Go testing, which does not require
// external frameworks like pytest.
func TestIsHighRisk(t *testing.T) {
	// A common Go pattern is table-driven testing using a slice of anonymous structs.
	// This makes it extremely easy to add new test cases compared to writing
	// individual `def test_...` functions for every single scenario.
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "Should flag simple delete",
			command:  "kubectl delete pods --all",
			expected: true,
		},
		{
			name:     "Should flag helm uninstall",
			command:  "helm uninstall my-release",
			expected: true,
		},
		{
			name:     "Should flag terraform destroy",
			command:  "terraform destroy -auto-approve",
			expected: true,
		},
		{
			name:     "Should flag apply",
			command:  "kubectl apply -f deployment.yaml",
			expected: true,
		},
		{
			name:     "Should pass read-only get",
			command:  "kubectl get pods",
			expected: false,
		},
		{
			name:     "Should pass read-only describe",
			command:  "kubectl describe node",
			expected: false,
		},
		{
			name:     "Should pass helm list",
			command:  "helm ls -A",
			expected: false,
		},
	}

	for _, tt := range tests {
		// t.Run executes subtests, outputting the name if it fails
		t.Run(tt.name, func(t *testing.T) {
			result := isHighRisk(tt.command)
			if result != tt.expected {
				// t.Errorf is how you assert failures in Go instead of `assert X == Y`
				t.Errorf("isHighRisk(%q) = %v; want %v", tt.command, result, tt.expected)
			}
		})
	}
}
