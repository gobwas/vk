package vk

import "strconv"

type Scope uint64

const (
	ScopeNotify        Scope = 1 << 0
	ScopeFriends       Scope = 1 << 1
	ScopePhotos        Scope = 1 << 2
	ScopeAudio         Scope = 1 << 3
	ScopeVideo         Scope = 1 << 4
	ScopePages         Scope = 1 << 7
	ScopeStatus        Scope = 1 << 10
	ScopeNotes         Scope = 1 << 11
	ScopeMessages      Scope = 1 << 12
	ScopeWall          Scope = 1 << 13
	ScopeAds           Scope = 1 << 15
	ScopeOffline       Scope = 1 << 16
	ScopeDocs          Scope = 1 << 17
	ScopeGroups        Scope = 1 << 18
	ScopeNotifications Scope = 1 << 19
	ScopeStats         Scope = 1 << 20
	ScopeEmail         Scope = 1 << 22
	ScopeMarket        Scope = 1 << 27
)

func (s Scope) String() string {
	return strconv.FormatUint(uint64(s), 10)
}
