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

func createRandomDocument(t *testing.T) Document {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	arg := CreateDocumentParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Title:          pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
		CreatedBy:      pgtype.Int4{Int32: user.ID, Valid: true},
	}

	document, err := testQueries.CreateDocument(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, document)

	require.Equal(t, arg.OrganizationID, document.OrganizationID)
	require.Equal(t, arg.Title, document.Title)
	require.Equal(t, arg.CreatedBy, document.CreatedBy)

	require.NotZero(t, document.ID)
	require.NotZero(t, document.CreatedAt)

	return document
}

func TestCreateDocument(t *testing.T) {
	createRandomDocument(t)
}

func TestGetDocument(t *testing.T) {
	document1 := createRandomDocument(t)
	document2, err := testQueries.GetDocument(context.Background(), document1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, document2)

	require.Equal(t, document1.ID, document2.ID)
	require.Equal(t, document1.OrganizationID, document2.OrganizationID)
	require.Equal(t, document1.Title, document2.Title)
	require.Equal(t, document1.CreatedBy, document2.CreatedBy)
	require.WithinDuration(t, document1.CreatedAt.Time, document2.CreatedAt.Time, time.Second)
}

func TestUpdateDocument(t *testing.T) {
	document1 := createRandomDocument(t)

	arg := UpdateDocumentParams{
		ID:    document1.ID,
		Title: pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
	}

	document2, err := testQueries.UpdateDocument(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, document2)

	require.Equal(t, document1.ID, document2.ID)
	require.Equal(t, arg.Title, document2.Title)
	require.Equal(t, document1.OrganizationID, document2.OrganizationID)
	require.Equal(t, document1.CreatedBy, document2.CreatedBy)
}

func TestDeleteDocument(t *testing.T) {
	document1 := createRandomDocument(t)

	err := testQueries.DeleteDocument(context.Background(), document1.ID)

	require.NoError(t, err)

	document2, err := testQueries.GetDocument(context.Background(), document1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, document2)
}

func TestListDocuments(t *testing.T) {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	for i := 0; i < 5; i++ {
		arg := CreateDocumentParams{
			OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
			Title:          pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
			CreatedBy:      pgtype.Int4{Int32: user.ID, Valid: true},
		}
		_, err := testQueries.CreateDocument(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListDocumentsParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Limit:          5,
		Offset:         0,
	}

	documents, err := testQueries.ListDocuments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, documents, 5)

	for _, document := range documents {
		require.NotEmpty(t, document)
		require.Equal(t, organization.ID, document.OrganizationID.Int32)
	}
}
