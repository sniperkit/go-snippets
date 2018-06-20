package command

import (
	"context"
	permsrv "github.com/chremoas/perms-srv/proto"
	"strings"
)

// Should perform.go setup its own client?
type Permissions struct {
	Client          permsrv.PermissionsService
	PermissionsList []string
}

func NewPermission(client permsrv.PermissionsService, permissionsList []string) *Permissions {
	// TODO: Check to make sure the permissions exist and if not, create them
	//perms, _ := client.ListPermissions(context.Background(), &permsrv.NilRequest{})

	return &Permissions{Client: client, PermissionsList: permissionsList}
}

func (p Permissions) CanPerform(ctx context.Context, sender string) (bool, error) {
	s := strings.Split(sender, ":")
	canPerform, err := p.Client.Perform(ctx,
		&permsrv.PermissionsRequest{
			User:            s[1],
			PermissionsList: p.PermissionsList,
		})

	if err != nil {
		return false, err
	}

	return canPerform.CanPerform, nil
}
