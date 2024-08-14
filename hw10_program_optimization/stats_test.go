//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStatAdditional(t *testing.T) {
	testCases := []struct {
		name      string
		data      string
		domain    string
		expected  DomainStat
		expectErr bool
	}{
		{
			name:      "single user with matching domain",
			data:      `{"Id":1,"Name":"John Doe","Email":"john.doe@example.com"}`,
			domain:    "com",
			expected:  DomainStat{"example.com": 1},
			expectErr: false,
		},
		{
			name: "multiple users with matching domain",
			data: `{"Id":1,"Name":"John Doe","Email":"john.doe@example.com"}
{"Id":2,"Name":"Jane Doe","Email":"jane.doe@example.com"}`,
			domain:    "com",
			expected:  DomainStat{"example.com": 2},
			expectErr: false,
		},
		{
			name:      "no matching domain",
			data:      `{"Id":1,"Name":"John Doe","Email":"john.doe@example.org"}`,
			domain:    "com",
			expected:  DomainStat{},
			expectErr: false,
		},
		{
			name:      "invalid json",
			data:      `{"Id":1,"Name":"John Doe","Email":"john.doe@example.com"`,
			domain:    "com",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "empty data",
			data:      ``,
			domain:    "com",
			expected:  DomainStat{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewBufferString(tc.data)
			result, err := GetDomainStat(r, tc.domain)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}
