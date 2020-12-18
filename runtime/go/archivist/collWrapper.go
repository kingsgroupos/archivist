package archivist

import "time"

type CollectionWrapper struct {
	Collection
	compatibleVersions map[string]struct{}
	when               time.Time
}
