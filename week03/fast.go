package main

import (
	"bufio"
	"fmt"
	"hw3/userdata"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := regexp.MustCompile("@")
	seenBrowsers := make(map[string]struct{})
	foundUsers := strings.Builder{}

	scanner := bufio.NewScanner(file)

	var users []userdata.UserData = nil

	for scanner.Scan() {
		user := userdata.UserData{}
		// fmt.Printf("%v %v\n", err, line)
		err := easyjson.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {

		flags := map[string]bool{
			"Android": false,
			"MSIE":    false,
		}

		browsers := user.Browsers

		for _, browser := range browsers {
			var match bool
			for key := range flags {
				if ok := strings.Contains(browser, key); ok {
					flags[key] = true
					match = true
				}
			}

			if !match {
				continue
			}

			seenBrowsers[browser] = struct{}{}
		}

		var match bool = true
		for _, v := range flags {
			if !v {
				match = false
				break
			}
		}

		if !match {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user.Email, " [at] ")
		foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
