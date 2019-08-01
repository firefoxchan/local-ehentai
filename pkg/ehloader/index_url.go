package ehloader

import (
	"bufio"
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//  File content format:
//    [{:schema}://][{:host}][/g]/{:gid}/{:token}[/]
//    {:gid}
//  File content example:
//    https://e-hentai.org/g/1111111/abcdef123
//    https://exhentai.org/g/1111111/abcdef123
//    /g/1111111/abcdef123
//    /1111111/abcdef123
//    1111111
func indexURLList(urlPath string) (map[int]struct{}, error) {
	matchedGs := make(map[int]struct{})
	f, e := os.OpenFile(urlPath, os.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		u, e := url.Parse(line)
		if e != nil {
			logger.Printf("Unable to parse url: %s, %s", e, line)
			continue
		}
		paths := strings.Split(strings.Trim(u.Path, "/"), "/")
		var gid64 int64
		switch len(paths) {
		case 1:
			// /{:gid}
			// {:gid}
			gid64, e = strconv.ParseInt(paths[0], 10, 64)
		case 2:
			// /{:gid}/{:token}
			gid64, e = strconv.ParseInt(paths[0], 10, 64)
		case 3:
			// /g/{:gid}/{:token}
			gid64, e = strconv.ParseInt(paths[1], 10, 64)
		default:
			gid64, e = 0, errors.New("invalid url path")
		}
		if e != nil {
			logger.Printf("Unable to parse gid / token in url: %s", line)
			continue
		}
		matchedGs[int(gid64)] = struct{}{}
	}
	if e := scanner.Err(); e != nil {
		return nil, e
	}
	return matchedGs, nil
}
