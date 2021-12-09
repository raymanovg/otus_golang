package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	subDomain := "." + domain
	res, err := countDomains(r, subDomain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return res, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	sc := bufio.NewScanner(r)
	user := &User{}
	for sc.Scan() {
		*user = User{}
		if err := jsoniter.Unmarshal(sc.Bytes(), user); err != nil {
			return nil, err
		}
		matchedDomain, ok := matchDomain(user.Email, domain)
		if ok {
			result[matchedDomain]++
		}
	}
	return result, nil
}

func matchDomain(email string, subDomain string) (string, bool) {
	if strings.Contains(email, subDomain) {
		return strings.ToLower(strings.SplitN(email, "@", 2)[1]), true
	}
	return "", false
}
