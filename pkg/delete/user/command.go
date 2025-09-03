/*
Copyright 2024-2025 the Unikorn Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package user

import (
	"context"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"github.com/unikorn-cloud/core/pkg/constants"
	identityv1 "github.com/unikorn-cloud/identity/pkg/apis/unikorn/v1alpha1"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/factory"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/flags"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/util"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type deleteUserOptions struct {
	UnikornFlags *factory.UnikornFlags

	user  *flags.UserFlags
	email string
}

func (o *deleteUserOptions) AddFlags(cmd *cobra.Command, factory *factory.Factory) error {
	cmd.Flags().StringVar(&o.email, "email", "", "User's email address.")

	if err := cmd.MarkFlagRequired("email"); err != nil {
		return err
	}

	return nil
}

func (o *deleteUserOptions) deleteOrgUserFromGroups(ctx context.Context, cli client.Client, organizationID, orgUserID string) error {
	groups, err := util.GetOrgGroups(ctx, cli, organizationID)
	if err != nil {
		return err
	}

	for i := range groups.Items {
		if slices.Contains(groups.Items[i].Spec.UserIDs, orgUserID) {
			groups.Items[i].Spec.UserIDs = slices.DeleteFunc(groups.Items[i].Spec.UserIDs, func(id string) bool {
				return id == orgUserID
			})

			if err := cli.Update(ctx, &groups.Items[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *deleteUserOptions) deleteOrgUser(ctx context.Context, cli client.Client, orgUser *identityv1.OrganizationUser) error {
	err := cli.Delete(ctx, orgUser)
	if err != nil {
		return err
	}

	return nil
}

func (o *deleteUserOptions) deleteUser(ctx context.Context, cli client.Client, user *identityv1.User) error {
	err := cli.Delete(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func Command(factory *factory.Factory) *cobra.Command {
	unikornFlags := &factory.UnikornFlags

	o := deleteUserOptions{
		UnikornFlags: unikornFlags,
		user:         flags.NewUserFlags(unikornFlags),
	}

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Comprehensive user deletion - remove from all groups and organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			client, err := factory.Client()
			if err != nil {
				return err
			}

			user, err := util.GetUser(ctx, client, o.UnikornFlags.IdentityNamespace, o.email)
			if err != nil {
				return err
			}

			userOrgs, err := util.GetOrgUsers(ctx, client, user.Name)
			if err != nil {
				return err
			}
			for i := range userOrgs.Items {
				orgID := userOrgs.Items[i].Labels[constants.OrganizationLabel]
				orgUserID := userOrgs.Items[i].Name
				err := o.deleteOrgUserFromGroups(ctx, client, orgID, orgUserID)
				if err != nil {
					return err
				}
				err = o.deleteOrgUser(ctx, client, &userOrgs.Items[i])
				if err != nil {
					return err
				}
			}
			// Pew pew
			err = o.deleteUser(ctx, client, user)
			if err != nil {
				return err
			}

			return nil
		},
	}

	if err := o.AddFlags(cmd, factory); err != nil {
		panic(err)
	}

	return cmd
}
