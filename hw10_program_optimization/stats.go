package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type (
	DomainStat map[string]int
	emails     [100_000]EmailOnly
	EmailOnly  struct {
		Email string `json:"email"`
	}
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	e, err := getEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get emails error: %w", err)
	}
	return countDomains(e, domain)
}

func getEmails(r io.Reader) (result emails, err error) {
	scanner := bufio.NewScanner(r)
	var e EmailOnly
	i := 0

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return result, err
		}
		result[i] = e
		i++
	}

	return result, scanner.Err()
}

func countDomains(e emails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domRE := regexp.MustCompile(`\.` + regexp.QuoteMeta(domain) + `$`)

	for _, email := range e {
		if domRE.MatchString(email.Email) {
			result[strings.ToLower(strings.SplitN(email.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
