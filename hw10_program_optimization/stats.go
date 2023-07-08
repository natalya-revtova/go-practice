package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat, 0)
	scanner := bufio.NewScanner(r)
	re := regexp.MustCompile(fmt.Sprintf("\\.%s$", domain))

	for scanner.Scan() {
		var user User
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if re.Match([]byte(user.Email)) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
