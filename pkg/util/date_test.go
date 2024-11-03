package util

import "testing"

// Offset for UTC+1 (CET)
func TestUnixToLocal_CET(t *testing.T) {
	unixTimestamp := int64(1733054400)
	timeOffset := 3600

	expectedTime := "2024-12-01 13:00"

	result := UnixToLocal(unixTimestamp, timeOffset)

	if result != expectedTime {
		t.Errorf("UnixToLocal(%d, %d) = %s; want %s", unixTimestamp, timeOffset, result, expectedTime)
	}
}

func TestUnixToUTC(t *testing.T) {
	// Define test cases
	tests := []struct {
		unix        int64
		expectedUTC string
	}{
		{
			unix:        1733054400, // 2024-12-01 12:00:00 UTC
			expectedUTC: "2024-12-01T12:00:00Z",
		},
		{
			unix:        0, // Unix epoch
			expectedUTC: "1970-01-01T00:00:00Z",
		},
		{
			unix:        1609459200, // 2021-01-01 00:00:00 UTC
			expectedUTC: "2021-01-01T00:00:00Z",
		},
	}

	for _, test := range tests {
		result := UnixToUTC(test.unix)
		if result != test.expectedUTC {
			t.Errorf("UnixToUTC(%d) = %s; want %s", test.unix, result, test.expectedUTC)
		}
	}
}
