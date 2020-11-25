package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"runtime"
	"strings"
	"sync"
	"unicode"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain), nil
}

func countDomains(r io.Reader, firstLvlDomain string) DomainStat {
	queueCh := make(chan string)
	resultCh := make(chan string)
	wg := &sync.WaitGroup{}

	prefixString := "\"Email\":\""
	prefixStringRunes := []rune(prefixString)

	for i := 0; i < runtime.NumCPU(); i++ {
		go processLine(firstLvlDomain, prefixStringRunes, wg, queueCh, resultCh)
	}

	result := make(DomainStat)

	doneCh := make(chan struct{})
	go aggregateResults(result, resultCh, doneCh)

	readFile(r, wg, queueCh)

	wg.Wait()
	close(resultCh)

	<-doneCh

	return result
}

func readFile(r io.Reader, wg *sync.WaitGroup, queueCh chan<- string) {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		//nolint: errorlint
		// т.к. error.Is делается медленно
		if err != nil && err != io.EOF {
			break
		}

		wg.Add(1)

		queueCh <- line

		//nolint: errorlint
		// т.к. error.Is делается медленно
		if err == io.EOF {
			break
		}
	}

	defer close(queueCh)
}

func aggregateResults(result DomainStat, resultCh <-chan string, doneCh chan<- struct{}) {
	for {
		domain, ok := <-resultCh
		if !ok {
			break
		}

		num := result[domain]
		num++
		result[domain] = num
	}

	defer close(doneCh)
}

// подразумевается, что входной json нормализован, валиден и не содержит лишних пробелов
// если нам важно быстродействие, то, я думаю, что это выполнимое условие
// (пояснение к nolint) за быстродействие частенько приходится платить читаемостью кода.
//nolint: funlen, gocognit
func processLine(firstLvlDomain string, prefixStringRunes []rune, wg *sync.WaitGroup, queueCh <-chan string, resultCh chan<- string) {
	for {
		line, ok := <-queueCh
		if !ok {
			break
		}

		var topDomain strings.Builder
		var domain strings.Builder
		inEmail := false
		inDomain := false
		inTopDomain := false
		var index int

		for _, r := range line {
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
				if topDomain.String() == firstLvlDomain {
					resultCh <- domain.String()
				}

				wg.Done()

				break
			}

			// начало домена
			if r == '@' {
				inDomain = true
				continue
			}

			if inDomain {
				domain.WriteRune(unicode.ToLower(r))

				// начало домена первого уровня
				// не учитывается вариант с email с доменом третьего уровня, например: test@sub.domain.com
				if r == '.' {
					inTopDomain = true
					continue
				}

				if inTopDomain {
					topDomain.WriteRune(unicode.ToLower(r))
				}
			}
		}
	}
}
