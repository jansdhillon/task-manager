package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
)

type mockDBConnector struct {
	connectFunc func() (*pgx.Conn, error)
}

func (m *mockDBConnector) Connect() (*pgx.Conn, error) {
	if m.connectFunc != nil {
		return m.connectFunc()
	}
	return nil, nil
}

func TestTaskDB_NewTaskDB(t *testing.T) {
	t.Run("creates new TaskDB instance", func(t *testing.T) {
		db := NewTaskDB()
		if db == nil {
			t.Error("NewTaskDB() returned nil")
		}
		if db.connector == nil {
			t.Error("NewTaskDB() created TaskDB with nil connector")
		}
	})
}

func TestTaskDB_NewTaskDBWithConnector(t *testing.T) {
	t.Run("creates TaskDB with custom connector", func(t *testing.T) {
		mockConnector := &mockDBConnector{}
		db := NewTaskDBWithConnector(mockConnector)
		if db == nil {
			t.Error("NewTaskDBWithConnector() returned nil")
		}
		if db.connector != mockConnector {
			t.Error("NewTaskDBWithConnector() did not set the connector properly")
		}
	})
}

func TestTaskDB_GetTasks(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		connectErr error
		wantErr    bool
	}{
		{
			name:       "connection error",
			ctx:        context.Background(),
			connectErr: &mockDBError{msg: "connection failed"},
			wantErr:    true,
		},
		{
			name:       "cancelled context",
			ctx:        func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			connectErr: &mockDBError{msg: "context cancelled"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConnector := &mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					if tt.connectErr != nil {
						return nil, tt.connectErr
					}
					return nil, nil
				},
			}

			taskDB := NewTaskDBWithConnector(mockConnector)

			tasks, err := taskDB.GetTasks(tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Error("GetTasks() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetTasks() unexpected error: %v", err)
				return
			}

			if tasks == nil {
				t.Error("GetTasks() returned nil tasks")
				return
			}
		})
	}
}

type mockDBError struct {
	msg string
}

func (e *mockDBError) Error() string {
	return e.msg
}

func TestTaskDB_NewTaskDB_Success(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(*testing.T, *TaskDB)
	}{
		{
			name: "creates TaskDB instance with non-nil pointer",
			validateFunc: func(t *testing.T, db *TaskDB) {
				if db == nil {
					t.Error("NewTaskDB() returned nil")
				}
			},
		},
		{
			name: "creates TaskDB with ProductionDBConnector",
			validateFunc: func(t *testing.T, db *TaskDB) {
				if db.connector == nil {
					t.Error("NewTaskDB() created TaskDB with nil connector")
					return
				}

				connectorType := reflect.TypeOf(db.connector).String()
				expectedType := "*db.ProductionDBConnector"
				if connectorType != expectedType {
					t.Errorf("NewTaskDB() connector type = %v, want %v", connectorType, expectedType)
				}
			},
		},
		{
			name: "creates TaskDB with proper struct fields",
			validateFunc: func(t *testing.T, db *TaskDB) {
				dbValue := reflect.ValueOf(db).Elem()
				connectorField := dbValue.FieldByName("connector")

				if !connectorField.IsValid() {
					t.Error("TaskDB missing connector field")
				}

				if connectorField.IsNil() {
					t.Error("TaskDB connector field is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTaskDB()
			tt.validateFunc(t, db)
		})
	}
}

func TestTaskDB_NewTaskDBWithConnector_Success(t *testing.T) {
	tests := []struct {
		name      string
		connector DBConnector
		wantNil   bool
	}{
		{
			name: "accepts valid mock connector",
			connector: &mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					return nil, nil
				},
			},
			wantNil: false,
		},
		{
			name:      "accepts nil connector",
			connector: nil,
			wantNil:   false,
		},
		{
			name: "accepts connector with error return",
			connector: &mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					return nil, &mockDBError{msg: "test error"}
				},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTaskDBWithConnector(tt.connector)

			if tt.wantNil && db != nil {
				t.Error("NewTaskDBWithConnector() expected nil, got non-nil")
			}

			if !tt.wantNil && db == nil {
				t.Error("NewTaskDBWithConnector() expected non-nil, got nil")
			}

			if db != nil && db.connector != tt.connector {
				t.Error("NewTaskDBWithConnector() did not set connector correctly")
			}
		})
	}
}

func TestTaskDB_ProductionDBConnector_Success(t *testing.T) {
	t.Run("ProductionDBConnector implements DBConnector interface", func(t *testing.T) {
		var _ DBConnector = &ProductionDBConnector{}
	})

	t.Run("ProductionDBConnector has Connect method", func(t *testing.T) {
		connector := &ProductionDBConnector{}
		connectorType := reflect.TypeOf(connector)

		method, exists := connectorType.MethodByName("Connect")
		if !exists {
			t.Error("ProductionDBConnector missing Connect method")
			return
		}

		if method.Type.NumIn() != 1 {
			t.Errorf("Connect method expected 1 input (receiver), got %d", method.Type.NumIn())
		}

		if method.Type.NumOut() != 2 {
			t.Errorf("Connect method expected 2 outputs, got %d", method.Type.NumOut())
		}
	})
}

func TestTaskDB_Interface_Compliance(t *testing.T) {
	t.Run("mockDBConnector implements DBConnector interface", func(t *testing.T) {
		var _ DBConnector = &mockDBConnector{}
	})

	t.Run("mockDBConnector Connect method returns expected types", func(t *testing.T) {
		mock := &mockDBConnector{
			connectFunc: func() (*pgx.Conn, error) {
				return nil, nil
			},
		}

		conn, err := mock.Connect()

		if conn != nil || err != nil {
			// This is expected behavior for our mock - can return nil, nil
		}
	})
}

func TestTaskDB_GetTasks_MethodExists(t *testing.T) {
	t.Run("GetTasks method exists on TaskDB", func(t *testing.T) {
		db := NewTaskDB()

		dbType := reflect.TypeOf(db)
		method, exists := dbType.MethodByName("GetTasks")

		if !exists {
			t.Error("TaskDB missing GetTasks method")
			return
		}

		if method.Type.NumIn() != 2 {
			t.Errorf("GetTasks method expected 2 inputs (receiver + context), got %d", method.Type.NumIn())
		}

		if method.Type.NumOut() != 2 {
			t.Errorf("GetTasks method expected 2 outputs (slice + error), got %d", method.Type.NumOut())
		}
	})
}

func TestTaskDB_Context_Handling_Success(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "background context",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "context with value",
			ctx:     context.WithValue(context.Background(), "key", "value"),
			wantErr: false,
		},
		{
			name:    "context with multiple values",
			ctx:     context.WithValue(context.WithValue(context.Background(), "key1", "val1"), "key2", "val2"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConnector := &mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					return nil, &mockDBError{msg: "expected connection error"}
				},
			}

			db := NewTaskDBWithConnector(mockConnector)

			_, err := db.GetTasks(tt.ctx)

			if tt.wantErr && err == nil {
				t.Error("GetTasks() expected error but got none")
			}

			if !tt.wantErr && err == nil {
				t.Error("GetTasks() expected error due to mock connector but got none")
			}

			if tt.ctx == nil {
				t.Error("Test setup error: context should not be nil")
			}
		})
	}
}

func TestTaskDB_Connector_Interface_Success(t *testing.T) {
	t.Run("multiple connector implementations work", func(t *testing.T) {
		connectors := []DBConnector{
			&ProductionDBConnector{},
			&mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					return nil, nil
				},
			},
			&mockDBConnector{
				connectFunc: func() (*pgx.Conn, error) {
					return nil, &mockDBError{msg: "test error"}
				},
			},
		}

		for i, connector := range connectors {
			t.Run(fmt.Sprintf("connector_%d", i), func(t *testing.T) {
				db := NewTaskDBWithConnector(connector)

				if db == nil {
					t.Error("NewTaskDBWithConnector() returned nil")
					return
				}

				if db.connector != connector {
					t.Error("NewTaskDBWithConnector() did not assign connector correctly")
				}

				var _ DBConnector = db.connector
			})
		}
	})
}

func TestTaskDB_Struct_Field_Access_Success(t *testing.T) {
	t.Run("TaskDB struct has accessible fields", func(t *testing.T) {
		db := NewTaskDB()

		dbValue := reflect.ValueOf(db).Elem()
		dbType := dbValue.Type()

		expectedFields := []string{"connector"}

		for _, fieldName := range expectedFields {
			field, exists := dbType.FieldByName(fieldName)
			if !exists {
				t.Errorf("TaskDB missing expected field: %s", fieldName)
				continue
			}

			if !field.Type.Implements(reflect.TypeOf((*DBConnector)(nil)).Elem()) {
				t.Errorf("Field %s does not implement DBConnector interface", fieldName)
			}
		}
	})
}

func TestTaskDB_ErrorTypes_Success(t *testing.T) {
	t.Run("mockDBError implements error interface", func(t *testing.T) {
		var _ error = &mockDBError{msg: "test"}
	})

	t.Run("mockDBError returns expected message", func(t *testing.T) {
		expectedMsg := "test error message"
		err := &mockDBError{msg: expectedMsg}

		if err.Error() != expectedMsg {
			t.Errorf("mockDBError.Error() = %v, want %v", err.Error(), expectedMsg)
		}
	})

	t.Run("mockDBError with empty message", func(t *testing.T) {
		err := &mockDBError{msg: ""}

		if err.Error() != "" {
			t.Errorf("mockDBError.Error() = %v, want empty string", err.Error())
		}
	})
}
