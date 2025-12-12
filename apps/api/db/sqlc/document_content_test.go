package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/utils"
	"github.com/stretchr/testify/require"
)

func createRandomDocumentContent(t *testing.T) DocumentContent {
	document := createRandomDocument(t)

	ydocState := []byte(utils.GenerateRandomString(20))

	arg := CreateDocumentContentParams{
		DocumentID: pgtype.Int4{Int32: document.ID, Valid: true},
		YdocState:  ydocState,
	}

	content, err := testQueries.CreateDocumentContent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, content)

	require.Equal(t, arg.DocumentID, content.DocumentID)
	require.Equal(t, arg.YdocState, content.YdocState)

	require.NotZero(t, content.UpdatedAt)

	return content
}

func TestCreateDocumentContent(t *testing.T) {
	createRandomDocumentContent(t)
}

func TestGetDocumentContent(t *testing.T) {
	content1 := createRandomDocumentContent(t)
	content2, err := testQueries.GetDocumentContent(context.Background(), content1.DocumentID)

	require.NoError(t, err)
	require.NotEmpty(t, content2)

	require.Equal(t, content1.DocumentID, content2.DocumentID)
	require.Equal(t, content1.YdocState, content2.YdocState)
	require.WithinDuration(t, content1.UpdatedAt.Time, content2.UpdatedAt.Time, time.Second)
}

func TestUpdateDocumentContent(t *testing.T) {
	content1 := createRandomDocumentContent(t)

	newYdocState := []byte(utils.GenerateRandomString(25))

	arg := UpdateDocumentContentParams{
		DocumentID: content1.DocumentID,
		YdocState:  newYdocState,
	}

	content2, err := testQueries.UpdateDocumentContent(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, content2)

	require.Equal(t, content1.DocumentID, content2.DocumentID)
	require.Equal(t, arg.YdocState, content2.YdocState)
}

func TestDeleteDocumentContent(t *testing.T) {
	content1 := createRandomDocumentContent(t)

	err := testQueries.DeleteDocumentContent(context.Background(), content1.DocumentID)

	require.NoError(t, err)

	content2, err := testQueries.GetDocumentContent(context.Background(), content1.DocumentID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, content2)
}
