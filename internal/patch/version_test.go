package patch_test

import (
	"github.com/justjack1521/mevpatch/internal/patch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion_Equal(t *testing.T) {

	t.Parallel()

	type test struct {
		name   string
		a      patch.Version
		b      patch.Version
		result bool
	}

	var tests = []test{
		{
			name: "true",
			a: patch.Version{
				Major: 1,
				Minor: 2,
				Patch: 1,
			},
			b: patch.Version{
				Major: 1,
				Minor: 2,
				Patch: 1,
			},
			result: true,
		},
		{
			name: "false",
			a: patch.Version{
				Major: 1,
				Minor: 2,
				Patch: 1,
			},
			b: patch.Version{
				Major: 1,
				Minor: 3,
				Patch: 1,
			},
			result: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var result = tc.a.Equal(tc.b)
			assert.Equal(t, tc.result, result)
		})
	}

}

func TestVersion_GeneratePreviousVersions(t *testing.T) {

	t.Parallel()

	type test struct {
		name     string
		current  patch.Version
		steps    int
		expected []patch.Version
	}

	var tests = []test{
		{
			name:     "standard_test",
			current:  patch.Version{Major: 2, Minor: 1, Patch: 2},
			steps:    5,
			expected: []patch.Version{{2, 1, 1}, {2, 1, 0}, {2, 0, 9}, {2, 0, 8}, {2, 0, 7}},
		},
		{
			name:     "end_early_test",
			current:  patch.Version{Patch: 3},
			steps:    5,
			expected: []patch.Version{{0, 0, 2}, {0, 0, 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var results = tc.current.GeneratePreviousVersions(tc.steps)
			assert.Equal(t, tc.expected, results)
		})
	}

}

func TestNewVersion(t *testing.T) {

	t.Parallel()

	type test struct {
		name     string
		value    string
		expected patch.Version
		fail     bool
	}

	var tests = []test{
		{
			name:  "valid_version",
			value: "1.2.1",
			expected: patch.Version{
				Major: 1,
				Minor: 2,
				Patch: 1,
			},
		},
		{
			name:  "invalid_version",
			value: "1.2",
			fail:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := patch.NewVersion(tc.value)
			if tc.fail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tc.expected.Equal(result))
			}
		})
	}

}
