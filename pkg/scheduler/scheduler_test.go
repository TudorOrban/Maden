package scheduler

import (
	"maden/pkg/shared"

	"testing"
)

func TestHasSufficientResources(t *testing.T) {
	tests := []struct {
		name string
		node shared.Node
		req shared.Resources
		expected bool
	}{
		{
			name: "sufficient resources",
			node: shared.Node{
				Capacity: shared.Resources{CPU: 10, Memory: 1000},
				Used: shared.Resources{CPU: 3, Memory: 500},
			},
			req: shared.Resources{CPU: 6, Memory: 400},
			expected: true,
		},
		{
            name: "insufficient CPU",
            node: shared.Node{
                Capacity: shared.Resources{CPU: 10, Memory: 1000},
                Used:     shared.Resources{CPU: 7, Memory: 200},
            },
            req: shared.Resources{CPU: 4, Memory: 100},
            expected: false,
        },
        {
            name: "insufficient Memory",
            node: shared.Node{
                Capacity: shared.Resources{CPU: 10, Memory: 1000},
                Used: shared.Resources{CPU: 2, Memory: 800},
            },
            req: shared.Resources{CPU: 2, Memory: 300},
            expected: false,
        },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hasSufficentResources(&test.node, &test.req); got != test.expected {
				t.Errorf("hasSufficientResources() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestMatchesAffinity(t *testing.T) {
	tests := []struct {
		name string
		node shared.Node
		pod shared.Pod
		expected bool
	}{
		{
			name: "affinity matched",
			node: shared.Node{
				Labels: map[string]string{"region": "us-west"},
			},
			pod: shared.Pod{
				Affinity: map[string]string{"region": "us-west"},
			},
			expected: true,
		},
		
        {
            name: "affinity not matched",
            node: shared.Node{
                Labels: map[string]string{"region": "us-east"},
            },
            pod: shared.Pod{
                Affinity: map[string]string{"region": "us-west"},
            },
            expected: false,
        },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := matchesAffinity(&test.node, &test.pod); got != test.expected {
				t.Errorf("matchesAffinity() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestMatchesAntiAffinity(t *testing.T) {
    tests := []struct {
        name string
        node shared.Node
        pod shared.Pod
        expected bool
    }{
        {
            name: "anti-affinity matched",
            node: shared.Node{
                Labels: map[string]string{"zone": "zone1"},
            },
            pod: shared.Pod{
                AntiAffinity: map[string]string{"zone": "zone1"},
            },
            expected: false,
        },
        {
            name: "anti-affinity not matched",
            node: shared.Node{
                Labels: map[string]string{"zone": "zone2"},
            },
            pod: shared.Pod{
                AntiAffinity: map[string]string{"zone": "zone1"},
            },
            expected: true,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            if got := matchesAntiAffinity(&test.node, &test.pod); got != test.expected {
                t.Errorf("matchesAntiAffinity() = %v, want %v", got, test.expected)
            }
        })
    }
}

func TestMatchesTolerations(t *testing.T) {
    tests := []struct {
        name string
        node shared.Node
        pod shared.Pod
        expected bool
    }{
        {
            name: "tolerations matched",
            node: shared.Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: shared.Pod{
                Tolerations: map[string]string{"key1": "value1"},
            },
            expected: true, 
        },
        {
            name: "tolerations not matched",
            node: shared.Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: shared.Pod{
                Tolerations: map[string]string{"key1": "value2"},
            },
            expected: false,
        },
        {
            name: "no tolerations for taints",
            node: shared.Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: shared.Pod{
                Tolerations: map[string]string{},
            },
            expected: false, 
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            if got := matchesTolerations(&test.node, &test.pod); got != test.expected {
                t.Errorf("matchesTolerations() = %v, want %v", got, test.expected)
            }
        })
    }
}
