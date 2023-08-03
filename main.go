package main

import (
	"bufio"
	"fmt"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/panjf2000/ants/v2"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	verifier = emailverifier.
		NewVerifier().
		EnableSMTPCheck()
)

func main() {
	defer ants.Release()
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		input := i.(string)

		emailRegex := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`

		re := regexp.MustCompile(emailRegex)

		email := re.FindString(input)
		split := strings.Split(email, "@")
		domain := split[1]
		username := split[0]
		ret, _ := verifier.CheckSMTP(domain, username)

		if ret.Deliverable {
			fmt.Println(input)
		}
		wg.Done()
	})
	defer p.Release()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		input := scanner.Text()
		p.Invoke(input)

	}

	wg.Wait()

}
