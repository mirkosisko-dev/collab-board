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

func createRandomTask(t *testing.T) Task {
	board := createRandomBoard(t)
	column := createRandomBoardColumn(t)
	assignee := createRandomUser(t)
	createdBy := createRandomUser(t)

	arg := CreateTaskParams{
		BoardID:     pgtype.Int4{Int32: board.ID, Valid: true},
		ColumnID:    pgtype.Int4{Int32: column.ID, Valid: true},
		AssigneeID:  pgtype.Int4{Int32: assignee.ID, Valid: true},
		Title:       pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
		Description: pgtype.Text{String: util.GenerateRandomString(20), Valid: true},
		Position:    pgtype.Int4{Int32: int32(util.GenerateRandomInt(1, 100)), Valid: true},
		CreatedBy:   pgtype.Int4{Int32: createdBy.ID, Valid: true},
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, arg.BoardID, task.BoardID)
	require.Equal(t, arg.ColumnID, task.ColumnID)
	require.Equal(t, arg.AssigneeID, task.AssigneeID)
	require.Equal(t, arg.Title, task.Title)
	require.Equal(t, arg.Description, task.Description)
	require.Equal(t, arg.Position, task.Position)
	require.Equal(t, arg.CreatedBy, task.CreatedBy)

	require.NotZero(t, task.ID)
	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t)
}

func TestGetTask(t *testing.T) {
	task1 := createRandomTask(t)
	task2, err := testQueries.GetTask(context.Background(), task1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, task2)

	require.Equal(t, task1.ID, task2.ID)
	require.Equal(t, task1.BoardID, task2.BoardID)
	require.Equal(t, task1.ColumnID, task2.ColumnID)
	require.Equal(t, task1.AssigneeID, task2.AssigneeID)
	require.Equal(t, task1.Title, task2.Title)
	require.Equal(t, task1.Description, task2.Description)
	require.Equal(t, task1.Position, task2.Position)
	require.Equal(t, task1.CreatedBy, task2.CreatedBy)
	require.WithinDuration(t, task1.CreatedAt.Time, task2.CreatedAt.Time, time.Second)
}

func TestUpdateTask(t *testing.T) {
	task1 := createRandomTask(t)
	newColumn := createRandomBoardColumn(t)
	newAssignee := createRandomUser(t)

	arg := UpdateTaskParams{
		ID:          task1.ID,
		ColumnID:    pgtype.Int4{Int32: newColumn.ID, Valid: true},
		AssigneeID:  pgtype.Int4{Int32: newAssignee.ID, Valid: true},
		Title:       pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
		Description: pgtype.Text{String: util.GenerateRandomString(20), Valid: true},
		Position:    pgtype.Int4{Int32: int32(util.GenerateRandomInt(1, 100)), Valid: true},
	}

	task2, err := testQueries.UpdateTask(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, task2)

	require.Equal(t, task1.ID, task2.ID)
	require.Equal(t, arg.ColumnID, task2.ColumnID)
	require.Equal(t, arg.AssigneeID, task2.AssigneeID)
	require.Equal(t, arg.Title, task2.Title)
	require.Equal(t, arg.Description, task2.Description)
	require.Equal(t, arg.Position, task2.Position)
}

func TestDeleteTask(t *testing.T) {
	task1 := createRandomTask(t)

	err := testQueries.DeleteTask(context.Background(), task1.ID)

	require.NoError(t, err)

	task2, err := testQueries.GetTask(context.Background(), task1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, task2)
}

func TestListTasks(t *testing.T) {
	board := createRandomBoard(t)
	column := createRandomBoardColumn(t)
	assignee := createRandomUser(t)
	createdBy := createRandomUser(t)

	for i := 0; i < 5; i++ {
		arg := CreateTaskParams{
			BoardID:     pgtype.Int4{Int32: board.ID, Valid: true},
			ColumnID:    pgtype.Int4{Int32: column.ID, Valid: true},
			AssigneeID:  pgtype.Int4{Int32: assignee.ID, Valid: true},
			Title:       pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
			Description: pgtype.Text{String: util.GenerateRandomString(20), Valid: true},
			Position:    pgtype.Int4{Int32: int32(i + 1), Valid: true},
			CreatedBy:   pgtype.Int4{Int32: createdBy.ID, Valid: true},
		}
		_, err := testQueries.CreateTask(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTasksParams{
		BoardID: pgtype.Int4{Int32: board.ID, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	tasks, err := testQueries.ListTasks(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tasks, 5)

	for _, task := range tasks {
		require.NotEmpty(t, task)
		require.Equal(t, board.ID, task.BoardID.Int32)
	}
}

func TestListTasksByColumn(t *testing.T) {
	board := createRandomBoard(t)
	column := createRandomBoardColumn(t)
	assignee := createRandomUser(t)
	createdBy := createRandomUser(t)

	for i := 0; i < 5; i++ {
		arg := CreateTaskParams{
			BoardID:     pgtype.Int4{Int32: board.ID, Valid: true},
			ColumnID:    pgtype.Int4{Int32: column.ID, Valid: true},
			AssigneeID:  pgtype.Int4{Int32: assignee.ID, Valid: true},
			Title:       pgtype.Text{String: util.GenerateRandomString(10), Valid: true},
			Description: pgtype.Text{String: util.GenerateRandomString(20), Valid: true},
			Position:    pgtype.Int4{Int32: int32(i + 1), Valid: true},
			CreatedBy:   pgtype.Int4{Int32: createdBy.ID, Valid: true},
		}
		_, err := testQueries.CreateTask(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTasksByColumnParams{
		ColumnID: pgtype.Int4{Int32: column.ID, Valid: true},
		Limit:    5,
		Offset:   0,
	}

	tasks, err := testQueries.ListTasksByColumn(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tasks, 5)

	for _, task := range tasks {
		require.NotEmpty(t, task)
		require.Equal(t, column.ID, task.ColumnID.Int32)
	}
}
