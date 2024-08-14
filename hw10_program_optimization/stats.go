package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		var user User
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			domainPart := strings.SplitN(user.Email, "@", 2)[1]
			result[strings.ToLower(domainPart)]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return result, nil
}
