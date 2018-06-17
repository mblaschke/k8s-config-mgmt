package configmanagement

var (
	statsNamespaces stats
	statsClusterObjects stats
	statsNamespaceObjects stats
)

type stats struct {
	created int
	updated int
	recreated int
	deleted int
}
