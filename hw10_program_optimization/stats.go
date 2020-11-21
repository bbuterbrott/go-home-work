package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"unicode"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain), nil
}

// за быстродействие частенько приходится платить читаемостью кода
//nolint: funlen
func countDomains(r io.Reader, firstLvlDomain string) DomainStat {
	result := make(DomainStat)
	prefixString := "\"Email\":\""
	prefixStringRunes := []rune(prefixString)

	// подразумевается, что входной json нормализован, валиден и не содержит лишних пробелов
	// если нам важно быстродействие, то, я думаю, что это выполнимое условие
	br := bufio.NewReader(r)
	var topDomainRunes []rune
	var domainRunes []rune
	inEmail := false
	inDomain := false
	inTopDomain := false
	var index int
	for {
		r, _, err := br.ReadRune()

		//nolint: errorlint
		// т.к. error.Is делается медленно
		if err == io.EOF {
			break
		}

		if !inEmail {
			if index == len(prefixStringRunes)-1 {
				inEmail = true
				continue
			}

			if r == prefixStringRunes[index] {
				index++
				continue
			}

			index = 0
			continue
		}

		// конец email
		if r == '"' {
			topDomain := string(topDomainRunes)

			if topDomain == firstLvlDomain {
				domain := string(domainRunes)
				num := result[domain]
				num++
				result[domain] = num
			}

			topDomainRunes = nil
			domainRunes = nil
			inEmail = false
			inDomain = false
			inTopDomain = false
			index = 0

			continue
		}

		// начало домена
		if r == '@' {
			inDomain = true
			continue
		}

		if inDomain {
			domainRunes = append(domainRunes, unicode.ToLower(r))

			// начало домена первого уровня
			// не учитывается вариант с email с доменом третьего уровня, например: test@sub.domain.com
			if r == '.' {
				inTopDomain = true
				continue
			}

			if inTopDomain {
				topDomainRunes = append(topDomainRunes, unicode.ToLower(r))
			}
		}
	}

	return result
}
