package main

import "regexp"

const defaultCharset = "utf-8"

var metaRx = regexp.MustCompile("(?is)<meta charset=\"(.*?)\"")

func findCharset(text string) string {
	match := metaRx.FindStringSubmatch(text)
	if len(match) == 0 {
		return defaultCharset
	}
	return match[1]
}
