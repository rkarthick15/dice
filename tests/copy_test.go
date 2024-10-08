package tests

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCopy(t *testing.T) {
	conn := getLocalConnection()
	defer conn.Close()

	simpleJSON := `{"name":"John","age":30}`

	testCases := []struct {
		name     string
		commands []string
		expected []interface{}
	}{
		{
			name:     "COPY when source key doesn't exist",
			commands: []string{"COPY k1 k2"},
			expected: []interface{}{int64(0)},
		},
		{
			name:     "COPY with no REPLACE",
			commands: []string{"SET k1 v1", "COPY k1 k2", "GET k1", "GET k2"},
			expected: []interface{}{"OK", int64(1), "v1", "v1"},
		},
		{
			name:     "COPY with REPLACE",
			commands: []string{"SET k1 v1", "SET k2 v2", "GET k2", "COPY k1 k2 REPLACE", "GET k2"},
			expected: []interface{}{"OK", "OK", "v2", int64(1), "v1"},
		},
		{
			name:     "COPY with JSON integer",
			commands: []string{"JSON.SET k1 $ 2", "COPY k1 k2", "JSON.GET k2"},
			expected: []interface{}{"OK", int64(1), "2"},
		},
		{
			name:     "COPY with JSON boolean",
			commands: []string{"JSON.SET k1 $ true", "COPY k1 k2", "JSON.GET k2"},
			expected: []interface{}{"OK", int64(1), "true"},
		},
		{
			name:     "COPY with JSON array",
			commands: []string{`JSON.SET k1 $ [1,2,3]`, "COPY k1 k2", "JSON.GET k2"},
			expected: []interface{}{"OK", int64(1), `[1,2,3]`},
		},
		{
			name:     "COPY with JSON simple JSON",
			commands: []string{`JSON.SET k1 $ ` + simpleJSON, "COPY k1 k2", "JSON.GET k2"},
			expected: []interface{}{"OK", int64(1), simpleJSON},
		},
        {
            name: "COPY with no expiry",
            commands: []string{"SET k1 v1", "COPY k1 k2", "TTL k1", "TTL k2"},
            expected: []interface{}{"OK", int64(1), int64(-1), int64(-1)},
        },
        {
            name: "COPY with expiry making sure copy expires",
            commands: []string{"SET k1 v1 EX 5", "COPY k1 k2", "GET k1", "GET k2", "SLEEP 7", "GET k1", "GET k2"},
            expected: []interface{}{"OK", int64(1), "v1", "v1", "OK", "(nil)", "(nil)"},
        },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteTestKeys([]string{"k1", "k2"})
			for i, cmd := range tc.commands {
				result := fireCommand(conn, cmd)
				assert.Equal(t, tc.expected[i], result, "Value mismatch for cmd %s", cmd)
			}
		})
	}
}
