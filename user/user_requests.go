package user

import (
	"github.com/datatogether/errors"
	"github.com/datatogether/sqlutil"
	"strings"
)

// Requests holds all types of requests for users
type UserRequests struct {
	Store sqlutil.Transactable
}

// UsersRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type UsersListParams struct {
	// the user performing the request
	User *User `required:"true"`
	Type UserType
	// users requests embeds pagination info
	Limit  int
	Offset int
}

func (r UserRequests) List(p *UsersListParams, res *[]*User) error {
	users, err := ReadUsers(r.Store, p.Type, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*res = users
	return nil
}

type UsersGetParams struct {
	User    *User
	Subject *User
}

func (r UserRequests) Get(p *UsersGetParams, res *User) error {
	if err := p.Subject.Read(r.Store); err != nil {
		return err
	}

	*res = *p.Subject
	return nil
}

type UsersCreateParams struct {
	User     *User
	Password string
}

func (r UserRequests) Create(p *UsersCreateParams, res *User) error {
	p.User.password = p.Password
	if strings.TrimSpace(p.User.password) == "" {
		return errors.ErrPasswordRequired
	}
	if strings.TrimSpace(p.User.Email) == "" {
		return errors.ErrEmailRequired
	}

	if err := p.User.Save(r.Store); err != nil {
		return err
	}

	*res = *p.User
	return nil
}

type UsersSaveParams struct {
	User    *User
	Subject *User
}

func (r UserRequests) Save(p *UsersSaveParams, res *User) error {
	// TODO - restore w community membership lookup
	// if !p.User.isAdmin && p.User.Id != p.Subject.Id {
	// 	return errors.ErrAccessDenied
	// }

	if err := p.Subject.Save(r.Store); err != nil {
		return err
	}

	*res = *p.Subject
	return nil
}

type UsersSearchParams struct {
	User   *User
	Query  string
	Limit  int
	Offset int
}

func (r UserRequests) Search(p *UsersSearchParams, res *[]*User) error {
	users, err := UsersSearch(r.Store, p.Query, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*res = users
	return nil
}

type UsersCommunityMembersParams struct {
	User      *User
	Community *User
	Order     string
	Limit     int
	Offset    int
}

func (r UserRequests) CommunityMembers(p *UsersCommunityMembersParams, res *[]*User) error {
	// override order for now
	p.Order = "community_users.joined DESC"

	users, err := CommunityUsers(r.Store, p.Community, p.Order, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*res = users
	return nil
}

type UsersCommunitiesParams struct {
	User   *User
	Order  string
	Limit  int
	Offset int
}

func (r UserRequests) UserCommunities(p *UsersCommunitiesParams, res *[]*User) error {
	// override order for now
	p.Order = "users.created DESC"

	users, err := UserCommunities(r.Store, p.User, p.Order, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*res = users
	return nil
}

// type UsersCreateCommunityParams struct {
// 	User *User
// 	Name string
// }

// func (r UserRequests) CreateCommunity(p *UsersCreateCommunityParams, res *User) error {
// 	return nil
// }

// type UsersDeleteCommunityParams struct {
// }

// func (r UserRequests) DeleteCommunity(p *UsersDeleteCommunityParams, res *User) error {
// 	return nil
// }

// type UsersCreateCommunityInviteParams struct {
// 	User *User
// }

// func (r UserRequests) CreateCommunityInvite(p *UsersCreateCommunityParams, res *User) error {
// 	return nil
// }

// type UsersRemoveCommunityUserParams struct {
// }

// func (r UserRequests) UsersRemoveCommunityUser(p *UsersRemoveCommunityUserParams, res *bool) error {
// 	return nil
// }
