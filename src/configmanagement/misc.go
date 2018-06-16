package configmanagement

import (
	"regexp"
	"fmt"
	"strings"
	)

var (
	filterValueAsRegexp = regexp.MustCompile("^\\/.+\\/$")
)

func filterValueToRegexp(value string) (filter *regexp.Regexp) {

	if filterValueAsRegexp.MatchString(value) {
		// remove trailing slash
		value = value[:len(value)-1]
		// remove leading slash
		value = value[1:]
	} else {
		// filter is not a regexp -> wildcards
		value = fmt.Sprintf("^%s$", regexp.QuoteMeta(value))
		value = strings.Replace(value, "\\?", ".", -1)
		value = strings.Replace(value, "\\*", ".+", -1)
	}

	filter = regexp.MustCompile(value)
	return
}
