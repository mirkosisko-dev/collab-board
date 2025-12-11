package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/util"
	"github.com/stretchr/testify/require"
)

func createRandomBoard(t *testing.T) Board {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	arg := CreateBoardParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Name:           pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
		CreatedBy:      pgtype.Int4{Int32: user.ID, Valid: true},
	}

	board, err := testQueries.CreateBoard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, board)

	require.Equal(t, arg.OrganizationID, board.OrganizationID)
	require.Equal(t, arg.Name, board.Name)
	require.Equal(t, arg.CreatedBy, board.CreatedBy)

	require.NotZero(t, board.ID)

	return board
}

func TestCreateBoard(t *testing.T) {
	createRandomBoard(t)
}

func TestGetBoard(t *testing.T) {
	board1 := createRandomBoard(t)
	board2, err := testQueries.GetBoard(context.Background(), board1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, board2)

	require.Equal(t, board1.ID, board2.ID)
	require.Equal(t, board1.OrganizationID, board2.OrganizationID)
	require.Equal(t, board1.Name, board2.Name)
	require.Equal(t, board1.CreatedBy, board2.CreatedBy)
}

func TestUpdateBoard(t *testing.T) {
	board1 := createRandomBoard(t)

	arg := UpdateBoardParams{
		ID:   board1.ID,
		Name: pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
	}

	board2, err := testQueries.UpdateBoard(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, board2)

	require.Equal(t, board1.ID, board2.ID)
	require.Equal(t, arg.Name, board2.Name)
	require.Equal(t, board1.OrganizationID, board2.OrganizationID)
	require.Equal(t, board1.CreatedBy, board2.CreatedBy)
}

func TestDeleteBoard(t *testing.T) {
	board1 := createRandomBoard(t)

	err := testQueries.DeleteBoard(context.Background(), board1.ID)

	require.NoError(t, err)

	board2, err := testQueries.GetBoard(context.Background(), board1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, board2)
}

func TestListBoards(t *testing.T) {
	organization := createRandomOrganization(t)
	user := createRandomUser(t)

	for i := 0; i < 5; i++ {
		arg := CreateBoardParams{
			OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
			Name:           pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
			CreatedBy:      pgtype.Int4{Int32: user.ID, Valid: true},
		}
		_, err := testQueries.CreateBoard(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListBoardsParams{
		OrganizationID: pgtype.Int4{Int32: organization.ID, Valid: true},
		Limit:          5,
		Offset:         0,
	}

	boards, err := testQueries.ListBoards(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, boards, 5)

	for _, board := range boards {
		require.NotEmpty(t, board)
		require.Equal(t, organization.ID, board.OrganizationID.Int32)
	}
}

