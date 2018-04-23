package header

import "net/http"

var Header http.Header

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36"

func init() {
	Header = make(http.Header)
	Header["User-Agent"] = []string{UserAgent}
}
