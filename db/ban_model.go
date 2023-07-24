package db

type BanType string

const (
	UserBan BanType = "user"
	IPBan   BanType = "ip"
)

type Ban struct {
	Type  BanType `json:"type"`
	Value string  `json:"value"`
}

var Bans = conn().Collection("bans")
