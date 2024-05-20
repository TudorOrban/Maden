package main

import "testing"

func TestHasSufficientResources(t *testing.T) {
	tests := []struct {
		name string
		node Node
		req Resources
		expected bool
	}{
		{
			name: "sufficient resources",
			node: Node{
				Capacity: Resources{CPU: 10, Memory: 1000},
				Used: Resources{CPU: 3, Memory: 500},
			},
			req: Resources{CPU: 6, Memory: 400},
			expected: true,
		},
		{
            name: "insufficient CPU",
            node: Node{
                Capacity: Resources{CPU: 10, Memory: 1000},
                Used:     Resources{CPU: 7, Memory: 200},
            },
            req:      Resources{CPU: 4, Memory: 100},
            expected: false,
        },
        {
            name: "insufficient Memory",
            node: Node{
                Capacity: Resources{CPU: 10, Memory: 1000},
                Used:     Resources{CPU: 2, Memory: 800},
            },
            req:      Resources{CPU: 2, Memory: 300},
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
		node Node
		pod Pod
		expected bool
	}{
		{
			name: "affinity matched",
			node: Node{
				Labels: map[string]string{"region": "us-west"},
			},
			pod: Pod{
				Affinity: map[string]string{"region": "us-west"},
			},
			expected: true,
		},
		
        {
            name: "affinity not matched",
            node: Node{
                Labels: map[string]string{"region": "us-east"},
            },
            pod: Pod{
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
        name     string
        node     Node
        pod      Pod
        expected bool
    }{
        {
            name: "anti-affinity matched",
            node: Node{
                Labels: map[string]string{"zone": "zone1"},
            },
            pod: Pod{
                AntiAffinity: map[string]string{"zone": "zone1"},
            },
            expected: false,  // Should return false as the anti-affinity condition matches
        },
        {
            name: "anti-affinity not matched",
            node: Node{
                Labels: map[string]string{"zone": "zone2"},
            },
            pod: Pod{
                AntiAffinity: map[string]string{"zone": "zone1"},
            },
            expected: true,  // Should return true as the anti-affinity condition does not match
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
        name     string
        node     Node
        pod      Pod
        expected bool
    }{
        {
            name: "tolerations matched",
            node: Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: Pod{
                Tolerations: map[string]string{"key1": "value1"},
            },
            expected: true,  // Should return true as tolerations match the taints
        },
        {
            name: "tolerations not matched",
            node: Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: Pod{
                Tolerations: map[string]string{"key1": "value2"},
            },
            expected: false,  // Should return false as tolerations do not match the taints
        },
        {
            name: "no tolerations for taints",
            node: Node{
                Taints: map[string]string{"key1": "value1"},
            },
            pod: Pod{
                Tolerations: map[string]string{},
            },
            expected: false,  // Should return false as no tolerations are provided for existing taints
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
