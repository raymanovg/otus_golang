package hw10programoptimization

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	done := make(chan struct{})

	usersCh, errCh := getUsers(r)
	domainCh := countDomains(usersCh, domain, done)

	for d := range domainCh {
		select {
		case err := <-errCh:
			close(done)
			return nil, fmt.Errorf("get users error: %w", err)
		default:
			result[d]++
		}
	}

	return result, nil
}

type users chan User

func getUsers(r io.Reader) (users, chan error) {
	usersCh := make(users)
	errCh := make(chan error)

	go func() {
		defer close(usersCh)
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			var user User
			if err := jsoniter.Unmarshal(sc.Bytes(), &user); err != nil {
				errCh <- err
				break
			}
			usersCh <- user
		}
	}()

	return usersCh, errCh
}

func countDomains(u users, domain string, done chan struct{}) chan string {
	domainCh := make(chan string)

	subst := "." + domain
	go func() {
		defer close(domainCh)
		for user := range u {
			select {
			case <-done:
				return
			default:
			}
			if strings.Contains(user.Email, subst) {
				domainCh <- strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			}
		}
	}()

	return domainCh
}
