package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jansdhillon/task-manager/server/.gen/postgres/public/model"
	"github.com/jansdhillon/task-manager/server/internal/task"
)

func TestAddTask(t *testing.T) {
	nonNilDesc := "w"
	testCases := []struct {
		desc        string
		taskTitle   string
		taskDesc    *string
		expectError bool
	}{
		{
			desc:      "add new task success",
			taskTitle: "mytask",
			taskDesc:  nil,
		},
		{
			desc:      "add new task with desc",
			taskTitle: "mytask",
			taskDesc:  &nonNilDesc,
		},
		{
			desc:        "add new task fail",
			taskTitle:   "mytask",
			expectError: true,
		},
	}
	for _, tC := range testCases {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New(): %s", err)
		}

		taskDB := NewTaskDB(mockDB)

		t.Cleanup(func() {
			_ = taskDB.Close()
			_ = mockDB.Close()
		})

		expectedSQL := `INSERT INTO public\.task \(title, description\)
        VALUES \(\$1, \$2\)
        RETURNING task\.id AS "task\.id",
                  task\.created_at AS "task\.created_at",
                  task\.description AS "task\.description",
                  task\.title AS "task\.title",
                  task\.last_updated_at AS "task\.last_updated_at",
                  task\.deleted AS "task\.deleted",
                  task\.status AS "task\.status"`

		if tC.expectError {
			mock.ExpectQuery(expectedSQL).
				WithArgs(tC.taskTitle, tC.taskDesc).
				WillReturnError(fmt.Errorf("database error"))
		} else {
			expectedID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
			expectedTime := time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC)

			mock.ExpectQuery(expectedSQL).
				WithArgs(tC.taskTitle, tC.taskDesc).
				WillReturnRows(sqlmock.NewRows([]string{
					"task.id", "task.created_at", "task.description", "task.title",
					"task.last_updated_at", "task.deleted", "task.status",
				}).AddRow(
					expectedID, expectedTime, tC.taskDesc, tC.taskTitle,
					expectedTime, false, model.Status_NotStarted,
				))
		}

		addedTask, err := taskDB.AddTask(t.Context(), tC.taskTitle, tC.taskDesc)
		if err != nil {
			if !tC.expectError {
				t.Fatalf("failed to add task to db: %s", err)
			}
		} else if tC.expectError {
			t.Fatal("expected error but got none")
		}

		if !tC.expectError {
			if addedTask == nil {
				t.Fatal("expected task but got nil")
			}
			if addedTask.Title != tC.taskTitle {
				t.Errorf("expected title %q, got %q", tC.taskTitle, addedTask.Title)
			}
			if (tC.taskDesc == nil && addedTask.Description != nil) ||
				(tC.taskDesc != nil && addedTask.Description == nil) ||
				(tC.taskDesc != nil && addedTask.Description != nil && *tC.taskDesc != *addedTask.Description) {
				t.Errorf("expected description %v, got %v", tC.taskDesc, addedTask.Description)
			}
			t.Logf("task: %s", addedTask)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetTask(t *testing.T) {
	testCases := []struct {
		desc        string
		taskID      uuid.UUID
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock, taskID uuid.UUID)
	}{
		{
			desc:   "get task success",
			taskID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			setupMock: func(mock sqlmock.Sqlmock, taskID uuid.UUID) {
				expectedTime := time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC)
				expectedSQL := `SELECT task\.id AS "task\.id",
                  task\.created_at AS "task\.created_at",
                  task\.description AS "task\.description",
                  task\.title AS "task\.title",
                  task\.last_updated_at AS "task\.last_updated_at",
                  task\.deleted AS "task\.deleted",
                  task\.status AS "task\.status"
FROM public\.task
WHERE task\.id = \$1`

				mock.ExpectQuery(expectedSQL).
					WithArgs(taskID).
					WillReturnRows(sqlmock.NewRows([]string{
						"task.id", "task.created_at", "task.description", "task.title",
						"task.last_updated_at", "task.deleted", "task.status",
					}).AddRow(
						taskID, expectedTime, nil, "Test Task",
						expectedTime, false, model.Status_NotStarted,
					))
			},
		},
		{
			desc:        "get task database error",
			taskID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock, taskID uuid.UUID) {
				expectedSQL := `SELECT task\.id AS "task\.id",
                  task\.created_at AS "task\.created_at",
                  task\.description AS "task\.description",
                  task\.title AS "task\.title",
                  task\.last_updated_at AS "task\.last_updated_at",
                  task\.deleted AS "task\.deleted",
                  task\.status AS "task\.status"
FROM public\.task
WHERE task\.id = \$1`

				mock.ExpectQuery(expectedSQL).
					WithArgs(taskID).
					WillReturnError(fmt.Errorf("database error"))
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("sqlmock.New(): %s", err)
			}

			taskDB := NewTaskDB(mockDB)
			t.Cleanup(func() {
				_ = taskDB.Close()
				_ = mockDB.Close()
			})

			tC.setupMock(mock, tC.taskID)

			result, err := taskDB.GetTask(context.Background(), tC.taskID)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result == nil {
					t.Fatal("expected task but got nil")
				}
				if result.ID != tC.taskID {
					t.Errorf("expected ID %v, got %v", tC.taskID, result.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	newDesc := "Updated description"

	testCases := []struct {
		desc        string
		taskID      uuid.UUID
		title       string
		description *string
		status      task.Status
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			desc:        "update task success",
			taskID:      testID,
			title:       "Updated Task",
			description: &newDesc,
			status:      task.InProgress,
			setupMock: func(mock sqlmock.Sqlmock) {
				expectedTime := time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC)

				getSQL := `SELECT task\.id AS "task\.id",
                  task\.created_at AS "task\.created_at",
                  task\.description AS "task\.description",
                  task\.title AS "task\.title",
                  task\.last_updated_at AS "task\.last_updated_at",
                  task\.deleted AS "task\.deleted",
                  task\.status AS "task\.status"
FROM public\.task
WHERE task\.id = \$1`

				mock.ExpectQuery(getSQL).
					WithArgs(testID).
					WillReturnRows(sqlmock.NewRows([]string{
						"task.id", "task.created_at", "task.description", "task.title",
						"task.last_updated_at", "task.deleted", "task.status",
					}).AddRow(
						testID, expectedTime, "Old desc", "Old Task",
						expectedTime, false, model.Status_NotStarted,
					))

				updateSQL := `UPDATE public\.task SET \(created_at, description, title, last_updated_at, deleted, status\) = \(\$1, \$2, \$3, \$4, \$5, \$6\)
WHERE task\.id = \$7
RETURNING task\.id AS "task\.id",
          task\.created_at AS "task\.created_at",
          task\.description AS "task\.description",
          task\.title AS "task\.title",
          task\.last_updated_at AS "task\.last_updated_at",
          task\.deleted AS "task\.deleted",
          task\.status AS "task\.status"`

				mock.ExpectQuery(updateSQL).
					WithArgs(expectedTime, &newDesc, "Updated Task", sqlmock.AnyArg(), false, model.Status_InProgress, testID).
					WillReturnRows(sqlmock.NewRows([]string{
						"task.id", "task.created_at", "task.description", "task.title",
						"task.last_updated_at", "task.deleted", "task.status",
					}).AddRow(
						testID, expectedTime, &newDesc, "Updated Task",
						expectedTime, false, model.Status_InProgress,
					))
			},
		},
		{
			desc:        "update task not found",
			taskID:      testID,
			title:       "Updated Task",
			status:      task.InProgress,
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				getSQL := `SELECT task\.id AS "task\.id",
                  task\.created_at AS "task\.created_at",
                  task\.description AS "task\.description",
                  task\.title AS "task\.title",
                  task\.last_updated_at AS "task\.last_updated_at",
                  task\.deleted AS "task\.deleted",
                  task\.status AS "task\.status"
FROM public\.task
WHERE task\.id = \$1`

				mock.ExpectQuery(getSQL).
					WithArgs(testID).
					WillReturnError(fmt.Errorf("task not found"))
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("sqlmock.New(): %s", err)
			}

			taskDB := NewTaskDB(mockDB)
			t.Cleanup(func() {
				_ = taskDB.Close()
				_ = mockDB.Close()
			})

			tC.setupMock(mock)

			result, err := taskDB.UpdateTask(context.Background(), tC.taskID, tC.title, tC.description, tC.status)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result == nil {
					t.Fatal("expected task but got nil")
				}
				if result.Title != tC.title {
					t.Errorf("expected title %q, got %q", tC.title, result.Title)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc        string
		taskID      uuid.UUID
		expectError bool
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			desc:   "delete task success",
			taskID: testID,
			setupMock: func(mock sqlmock.Sqlmock) {
				deleteSQL := `DELETE FROM public\.task WHERE task\.id = \$1`
				mock.ExpectExec(deleteSQL).
					WithArgs(testID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			desc:        "delete task database error",
			taskID:      testID,
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				deleteSQL := `DELETE FROM public\.task WHERE task\.id = \$1`
				mock.ExpectExec(deleteSQL).
					WithArgs(testID).
					WillReturnError(fmt.Errorf("database error"))
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("sqlmock.New(): %s", err)
			}

			taskDB := NewTaskDB(mockDB)
			t.Cleanup(func() {
				_ = taskDB.Close()
				_ = mockDB.Close()
			})

			tC.setupMock(mock)

			err = taskDB.DeleteTask(context.Background(), tC.taskID)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskFromDBTask(t *testing.T) {
	testCases := []struct {
		desc        string
		dbTask      *model.Task
		expectError bool
		expectedErr string
	}{
		{
			desc:        "nil db task",
			dbTask:      nil,
			expectError: true,
			expectedErr: "db task is nil",
		},
		{
			desc: "valid db task",
			dbTask: &model.Task{
				ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Title:         "Test Task",
				Description:   nil,
				CreatedAt:     time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC),
				LastUpdatedAt: time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC),
				Deleted:       false,
				Status:        model.Status_NotStarted,
			},
		},
		{
			desc: "db task with description",
			dbTask: &model.Task{
				ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Title:         "Test Task",
				Description:   stringPtr("Test description"),
				CreatedAt:     time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC),
				LastUpdatedAt: time.Date(2025, 10, 26, 12, 0, 0, 0, time.UTC),
				Deleted:       false,
				Status:        model.Status_InProgress,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := TaskFromDBTask(context.Background(), tC.dbTask)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tC.expectedErr != "" && err.Error() != tC.expectedErr {
					t.Errorf("expected error %q, got %q", tC.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result == nil {
					t.Fatal("expected task but got nil")
				}

				if result.ID != tC.dbTask.ID {
					t.Errorf("expected ID %v, got %v", tC.dbTask.ID, result.ID)
				}
				if result.Title != tC.dbTask.Title {
					t.Errorf("expected title %q, got %q", tC.dbTask.Title, result.Title)
				}
				if !equalStringPtr(result.Description, tC.dbTask.Description) {
					t.Errorf("expected description %v, got %v", tC.dbTask.Description, result.Description)
				}
			}
		})
	}
}

func TestTaskStatusFromModelStatus(t *testing.T) {
	testCases := []struct {
		desc         string
		modelStatus  model.Status
		expectedTask task.Status
		expectError  bool
	}{
		{
			desc:         "NotStarted status",
			modelStatus:  model.Status_NotStarted,
			expectedTask: task.NotStarted,
		},
		{
			desc:         "InProgress status",
			modelStatus:  model.Status_InProgress,
			expectedTask: task.InProgress,
		},
		{
			desc:         "Completed status",
			modelStatus:  model.Status_Completed,
			expectedTask: task.Completed,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := taskStatusFromModelStatus(tC.modelStatus)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != tC.expectedTask {
					t.Errorf("expected status %v, got %v", tC.expectedTask, result)
				}
			}
		})
	}
}

func TestModelStatusFromTaskStatus(t *testing.T) {
	testCases := []struct {
		desc          string
		taskStatus    task.Status
		expectedModel model.Status
	}{
		{
			desc:          "NotStarted status",
			taskStatus:    task.NotStarted,
			expectedModel: model.Status_NotStarted,
		},
		{
			desc:          "InProgress status",
			taskStatus:    task.InProgress,
			expectedModel: model.Status_InProgress,
		},
		{
			desc:          "Completed status",
			taskStatus:    task.Completed,
			expectedModel: model.Status_Completed,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := modelStatusFromTaskStatus(tC.taskStatus)
			if result != tC.expectedModel {
				t.Errorf("expected status %v, got %v", tC.expectedModel, result)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
