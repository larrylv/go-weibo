package weibo

// User represents a Weibo user.
type User struct {
	ID               *int    `json:"id,omitempty"`
	ScreeName        *string `json:"screen_name,omitempty"`
	Name             *string `json:"name,omitempty"`
	Province         *string `json:"province,omitempty"`
	City             *string `json:"city,omitempty"`
	Location         *string `json:"location,omitempty"`
	Description      *string `json:"description,omitempty"`
	URL              *string `json:"url,omitempty"`
	ProfileImageUrl  *string `json:"profile_image_url,omitempty"`
	Domain           *string `json:"domain,omitempty"`
	Gender           *string `json:"gender,omitempty"`
	FollowersCount   *int    `json:"followers_count,omitempty"`
	FriendsCount     *int    `json:"friends_count,omitempty"`
	StatusesCount    *int    `json:"statuses_count,omitempty"`
	FavouritesCount  *int    `json:"favourites_count,omitempty"`
	CreatedAt        *string `json:"created_at,omitempty"`
	Following        *bool   `json:"following,omitempty"`
	AllowAllActMsg   *bool   `json:"allow_all_act_msg,omitempty"`
	GeoEnabled       *bool   `json:"geo_enabled,omitempty"`
	Verified         *bool   `json:"verified,omitempty"`
	AllowAllComment  *bool   `json:"allow_all_comment,omitempty"`
	AvatarLarge      *string `json:"avatar_large,omitempty"`
	VerifiedReason   *string `json:"verified_reason,omitempty"`
	FollowMe         *bool   `json:"follow_me,omitempty"`
	OnlineStatus     *int    `json:"online_status,omitempty"`
	BiFollowersCount *int    `json:"bi_followers_count,omitempty"`
}
