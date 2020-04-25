package module

type HeaderInfo struct {
	ID        int
	IP        string
	Host      string
	Referer   string
	UserAgent string

	Method string
	Path   string
}
