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

func createRandomMessage(t *testing.T) Message {
	board := createRandomBoard(t)
	organization := createRandomOrganization(t)
	user := createRandomUser(t)
	document := createRandomDocument(t)

	arg := CreateMessageParams{
		BoardID:        pgtype.Int4{Int32: board.ID, Valid: true},
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
		DocumentID:     pgtype.Int4{Int32: document.ID, Valid: true},
		Content:        pgtype.Text{String: utils.GenerateRandomString(20), Valid: true},
	}

	message, err := testQueries.CreateMessage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message)

	require.Equal(t, arg.BoardID, message.BoardID)
	require.Equal(t, arg.OrganizationID, message.OrganizationID)
	require.Equal(t, arg.UserID, message.UserID)
	require.Equal(t, arg.DocumentID, message.DocumentID)
	require.Equal(t, arg.Content, message.Content)

	require.NotZero(t, message.ID)
	require.NotZero(t, message.CreatedAt)

	return message
}

func TestCreateMessage(t *testing.T) {
	createRandomMessage(t)
}

func TestGetMessage(t *testing.T) {
	message1 := createRandomMessage(t)
	message2, err := testQueries.GetMessage(context.Background(), message1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, message2)

	require.Equal(t, message1.ID, message2.ID)
	require.Equal(t, message1.BoardID, message2.BoardID)
	require.Equal(t, message1.OrganizationID, message2.OrganizationID)
	require.Equal(t, message1.UserID, message2.UserID)
	require.Equal(t, message1.DocumentID, message2.DocumentID)
	require.Equal(t, message1.Content, message2.Content)
	require.WithinDuration(t, message1.CreatedAt.Time, message2.CreatedAt.Time, time.Second)
}

func TestUpdateMessage(t *testing.T) {
	message1 := createRandomMessage(t)

	arg := UpdateMessageParams{
		ID:      message1.ID,
		Content: pgtype.Text{String: utils.GenerateRandomString(20), Valid: true},
	}

	message2, err := testQueries.UpdateMessage(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, message2)

	require.Equal(t, message1.ID, message2.ID)
	require.Equal(t, arg.Content, message2.Content)
	require.Equal(t, message1.BoardID, message2.BoardID)
	require.Equal(t, message1.OrganizationID, message2.OrganizationID)
	require.Equal(t, message1.UserID, message2.UserID)
	require.Equal(t, message1.DocumentID, message2.DocumentID)
}

func TestDeleteMessage(t *testing.T) {
	message1 := createRandomMessage(t)

	err := testQueries.DeleteMessage(context.Background(), message1.ID)

	require.NoError(t, err)

	message2, err := testQueries.GetMessage(context.Background(), message1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, message2)
}

func TestListMessages(t *testing.T) {
	board := createRandomBoard(t)
	organization := createRandomOrganization(t)
	user := createRandomUser(t)
	document := createRandomDocument(t)

	for i := 0; i < 5; i++ {
		arg := CreateMessageParams{
			BoardID:        pgtype.Int4{Int32: board.ID, Valid: true},
			OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
			UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
			DocumentID:     pgtype.Int4{Int32: document.ID, Valid: true},
			Content:        pgtype.Text{String: utils.GenerateRandomString(20), Valid: true},
		}
		_, err := testQueries.CreateMessage(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListMessagesParams{
		BoardID: pgtype.Int4{Int32: board.ID, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	messages, err := testQueries.ListMessages(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, messages, 5)

	for _, message := range messages {
		require.NotEmpty(t, message)
		require.Equal(t, board.ID, message.BoardID.Int32)
	}
}

func TestListMessagesByOrganization(t *testing.T) {
	board := createRandomBoard(t)
	organization := createRandomOrganization(t)
	user := createRandomUser(t)
	document := createRandomDocument(t)

	for i := 0; i < 5; i++ {
		arg := CreateMessageParams{
			BoardID:        pgtype.Int4{Int32: board.ID, Valid: true},
			OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
			UserID:         pgtype.Int4{Int32: user.ID, Valid: true},
			DocumentID:     pgtype.Int4{Int32: document.ID, Valid: true},
			Content:        pgtype.Text{String: utils.GenerateRandomString(20), Valid: true},
		}
		_, err := testQueries.CreateMessage(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListMessagesByOrganizationParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Limit:          5,
		Offset:         0,
	}

	messages, err := testQueries.ListMessagesByOrganization(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, messages, 5)

	for _, message := range messages {
		require.NotEmpty(t, message)
		require.Equal(t, organization.ID, message.OrganizationID.Int32)
	}
}
