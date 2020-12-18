// Code generated by archivist. DO NOT EDIT.

package conf

import (
	"time"

	"github.com/pkg/errors"
	"gitlab-ee.funplus.io/watcher/watcher/archivist/runtime/go/archivist"
	"gitlab-ee.funplus.io/watcher/watcher/misc/wtime"
)

var (
	_ = time.After
	_ = errors.New
	_ = archivist.NewArchivist
	_ = wtime.ParseDuration
)

// easyjson:json
type ExampleArrayConf []*ExampleArrayConf_1771430469

// ExampleArrayConf_1771430469 represents /[]
type ExampleArrayConf_1771430469 struct {
	AvatarUrl         string `json:"avatar_url" bson:"avatar_url"`
	Blog              string `json:"blog" bson:"blog"`
	CreatedAt         string `json:"created_at" bson:"created_at"`
	EventsUrl         string `json:"events_url" bson:"events_url"`
	Followers         int64  `json:"followers" bson:"followers"`
	FollowersUrl      string `json:"followers_url" bson:"followers_url"`
	Following         int64  `json:"following" bson:"following"`
	FollowingUrl      string `json:"following_url" bson:"following_url"`
	GistsUrl          string `json:"gists_url" bson:"gists_url"`
	GravatarId        string `json:"gravatar_id" bson:"gravatar_id"`
	Hireable          bool   `json:"hireable" bson:"hireable"`
	HtmlUrl           string `json:"html_url" bson:"html_url"`
	Id                int64  `json:"id" bson:"_id"`
	Location          string `json:"location" bson:"location"`
	Login             string `json:"login" bson:"login"`
	Name              string `json:"name" bson:"name"`
	OrganizationsUrl  string `json:"organizations_url" bson:"organizations_url"`
	PublicGists       int64  `json:"public_gists" bson:"public_gists"`
	PublicRepos       int64  `json:"public_repos" bson:"public_repos"`
	ReceivedEventsUrl string `json:"received_events_url" bson:"received_events_url"`
	ReposUrl          string `json:"repos_url" bson:"repos_url"`
	SiteAdmin         bool   `json:"site_admin" bson:"site_admin"`
	StarredUrl        string `json:"starred_url" bson:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url" bson:"subscriptions_url"`
	Type              string `json:"type" bson:"type"`
	UpdatedAt         string `json:"updated_at" bson:"updated_at"`
	Url               string `json:"url" bson:"url"`
}

func (this *ExampleArrayConf) bindRefs(c *Collection) error {
	if this == nil {
		return nil
	}

	var ok bool
	_ = ok
	v := *this
	for _, v := range v {
		if err := v.bindRefs(c); err != nil {
			return err
		}
	}

	return nil
}

func (this *ExampleArrayConf_1771430469) bindRefs(c *Collection) error {
	if this == nil {
		return nil
	}

	var ok bool
	_ = ok

	return nil
}
