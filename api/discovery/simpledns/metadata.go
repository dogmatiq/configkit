package simpledns

type key string

const (
	// QueryHostKey is the target meta-data key for the hostname used in the DNS
	// query that discovered the target.
	QueryHostKey key = "simpledns.query-host"
)
