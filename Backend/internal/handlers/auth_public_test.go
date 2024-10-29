//go:build public

package handlers

import (
	"context"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util/testutil"
	"github.com/pashagolub/pgxmock/v4"
)

func TestAuthHandlerHandleGet(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleGet, mock, tests)
}

func TestAuthHandlerHandlePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandlePost, mock, tests)
}

func TestAuthHandlerHandlePut(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)
	
	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandlePut, mock, tests)
}

func TestAuthHandlerHandleDelete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleDelete, mock, tests)
}

func TestAuthHandlerHandleLogin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleLogin, mock, tests)
}
