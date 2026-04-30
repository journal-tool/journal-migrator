package replication

import (
	"fmt"
	"regexp"
)

var localhostRegex = regexp.MustCompile(
	fmt.Sprintf("(?i)(%s|%s)",
		regexp.QuoteMeta("localhost"),
		regexp.QuoteMeta("127.0.0.1"),
	),
)

type ReplicationHost struct {
	Host string
	Port string
}

func (host *ReplicationHost) IsSource() bool {
	return localhostRegex.Match([]byte(host.Host))
}
