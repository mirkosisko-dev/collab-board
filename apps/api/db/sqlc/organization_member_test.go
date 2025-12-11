package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrganizationMember(t *testing.T) OrganizationMember {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	arg := CreateOrganizationMemberParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
		Role:           pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
	}

	member, err := testQueries.CreateOrganizationMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, member)

	require.Equal(t, arg.OrganizationID, member.OrganizationID)
	require.Equal(t, arg.UserID, member.UserID)
	require.Equal(t, arg.Role, member.Role)

	require.NotZero(t, member.ID)
	require.NotZero(t, member.CreatedAt)

	return member
}

func TestCreateOrganizationMember(t *testing.T) {
	createRandomOrganizationMember(t)
}

func TestGetOrganizationMember(t *testing.T) {
	member1 := createRandomOrganizationMember(t)
	member2, err := testQueries.GetOrganizationMember(context.Background(), member1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.ID, member2.ID)
	require.Equal(t, member1.OrganizationID, member2.OrganizationID)
	require.Equal(t, member1.UserID, member2.UserID)
	require.Equal(t, member1.Role, member2.Role)
	require.WithinDuration(t, member1.CreatedAt.Time, member2.CreatedAt.Time, time.Second)
}

func TestGetOrganizationMemberByOrgAndUser(t *testing.T) {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	arg := CreateOrganizationMemberParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
		Role:           pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
	}

	member1, err := testQueries.CreateOrganizationMember(context.Background(), arg)
	require.NoError(t, err)

	member2, err := testQueries.GetOrganizationMemberByOrgAndUser(context.Background(), GetOrganizationMemberByOrgAndUserParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, member2)
	require.Equal(t, member1.ID, member2.ID)
}

func TestUpdateOrganizationMember(t *testing.T) {
	member1 := createRandomOrganizationMember(t)

	arg := UpdateOrganizationMemberParams{
		ID:   member1.ID,
		Role: pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
	}

	member2, err := testQueries.UpdateOrganizationMember(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.ID, member2.ID)
	require.Equal(t, arg.Role, member2.Role)
	require.Equal(t, member1.OrganizationID, member2.OrganizationID)
	require.Equal(t, member1.UserID, member2.UserID)
}

func TestDeleteOrganizationMember(t *testing.T) {
	member1 := createRandomOrganizationMember(t)

	err := testQueries.DeleteOrganizationMember(context.Background(), member1.ID)

	require.NoError(t, err)

	member2, err := testQueries.GetOrganizationMember(context.Background(), member1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, member2)
}

func TestListOrganizationMembers(t *testing.T) {
	organization := createRandomOrganization(t)

	for i := 0; i < 5; i++ {
		user := createRandomUser(t)
		arg := CreateOrganizationMemberParams{
			OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
			UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
			Role:           pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
		}
		_, err := testQueries.CreateOrganizationMember(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListOrganizationMembersParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Limit:         5,
		Offset:        0,
	}

	members, err := testQueries.ListOrganizationMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, members, 5)

	for _, member := range members {
		require.NotEmpty(t, member)
		require.Equal(t, organization.ID, member.OrganizationID.Int32)
	}
}

