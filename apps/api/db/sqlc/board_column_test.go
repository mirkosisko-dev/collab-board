package sqlc

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/utils"
	"github.com/stretchr/testify/require"
)

func createRandomBoardColumn(t *testing.T) BoardColumn {
	board := createRandomBoard(t)

	arg := CreateBoardColumnParams{
		BoardID:  pgtype.Int4{Int32: board.ID, Valid: true},
		Name:     pgtype.Text{String: utils.GenerateRandomString(10), Valid: true},
		Position: pgtype.Int4{Int32: int32(utils.GenerateRandomInt(1, 100)), Valid: true},
	}

	column, err := testQueries.CreateBoardColumn(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, column)

	require.Equal(t, arg.BoardID, column.BoardID)
	require.Equal(t, arg.Name, column.Name)
	require.Equal(t, arg.Position, column.Position)

	require.NotZero(t, column.ID)

	return column
}

func TestCreateBoardColumn(t *testing.T) {
	createRandomBoardColumn(t)
}

func TestGetBoardColumn(t *testing.T) {
	column1 := createRandomBoardColumn(t)
	column2, err := testQueries.GetBoardColumn(context.Background(), column1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, column2)

	require.Equal(t, column1.ID, column2.ID)
	require.Equal(t, column1.BoardID, column2.BoardID)
	require.Equal(t, column1.Name, column2.Name)
	require.Equal(t, column1.Position, column2.Position)
}

func TestUpdateBoardColumn(t *testing.T) {
	column1 := createRandomBoardColumn(t)

	arg := UpdateBoardColumnParams{
		ID:       column1.ID,
		Name:     pgtype.Text{String: utils.GenerateRandomString(10), Valid: true},
		Position: pgtype.Int4{Int32: int32(utils.GenerateRandomInt(1, 100)), Valid: true},
	}

	column2, err := testQueries.UpdateBoardColumn(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, column2)

	require.Equal(t, column1.ID, column2.ID)
	require.Equal(t, arg.Name, column2.Name)
	require.Equal(t, arg.Position, column2.Position)
	require.Equal(t, column1.BoardID, column2.BoardID)
}

func TestDeleteBoardColumn(t *testing.T) {
	column1 := createRandomBoardColumn(t)

	err := testQueries.DeleteBoardColumn(context.Background(), column1.ID)

	require.NoError(t, err)

	column2, err := testQueries.GetBoardColumn(context.Background(), column1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, column2)
}

func TestListBoardColumns(t *testing.T) {
	board := createRandomBoard(t)

	for i := 0; i < 5; i++ {
		arg := CreateBoardColumnParams{
			BoardID:  pgtype.Int4{Int32: board.ID, Valid: true},
			Name:     pgtype.Text{String: utils.GenerateRandomString(10), Valid: true},
			Position: pgtype.Int4{Int32: int32(i + 1), Valid: true},
		}
		_, err := testQueries.CreateBoardColumn(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListBoardColumnsParams{
		BoardID: pgtype.Int4{Int32: board.ID, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	columns, err := testQueries.ListBoardColumns(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, columns, 5)

	for _, column := range columns {
		require.NotEmpty(t, column)
		require.Equal(t, pgtype.Int4{Int32: board.ID, Valid: true}, column.BoardID)
	}
}
