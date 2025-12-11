package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mirkosisko-dev/api/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrganization(t *testing.T) Organization {
	orgName := util.GenerateRandomString(5)

	organization, err := testQueries.CreateOrganization(context.Background(), orgName)
	require.NoError(t, err)
	require.NotEmpty(t, organization)

	require.Equal(t, orgName, organization.Name)

	require.NotZero(t, organization.ID)
	require.NotZero(t, organization.CreatedAt)

	return organization
}

func TestCreateOrganization(t *testing.T) {
	createRandomOrganization(t)
}

func TestGetOrganization(t *testing.T) {
	organization1 := createRandomOrganization(t)
	organization2, err := testQueries.GetOrganization(context.Background(), organization1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, organization2)

	require.Equal(t, organization1.ID, organization2.ID)
	require.Equal(t, organization1.Name, organization2.Name)
	require.WithinDuration(t, organization1.CreatedAt.Time, organization2.CreatedAt.Time, time.Second)
}

func TestUpdateOrganization(t *testing.T) {
	organization1 := createRandomOrganization(t)

	arg := UpdateOrganizationParams{
		ID:   organization1.ID,
		Name: util.GenerateRandomString(5),
	}

	organization2, err := testQueries.UpdateOrganization(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, organization2)

	require.Equal(t, organization1.ID, organization2.ID)
	require.Equal(t, arg.Name, organization2.Name)
	require.WithinDuration(t, organization1.CreatedAt.Time, organization2.CreatedAt.Time, time.Second)
}

func TestDeleteOrganization(t *testing.T) {
	organization1 := createRandomOrganization(t)

	err := testQueries.DeleteOrganization(context.Background(), organization1.ID)

	require.NoError(t, err)

	organization2, err := testQueries.GetOrganization(context.Background(), organization1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, organization2)
}

func TestListOrganizations(t *testing.T) {
	for i := 0; i <= 10; i++ {
		createRandomOrganization(t)
	}

	arg := ListOrganizationsParams{
		Limit:  5,
		Offset: 5,
	}

	organizations, err := testQueries.ListOrganizations(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, organizations, 5)

	for _, organization := range organizations {
		require.NotEmpty(t, organization)
	}
}
