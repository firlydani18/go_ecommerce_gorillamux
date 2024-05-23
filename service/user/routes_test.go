package user

import (
	"bytes"
	"encoding/json"
	"go-ecommerce/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := model.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "",
			Password:  "asd",
		}
		jsonMarshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonMarshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.registerHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserById(id int) (*model.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user model.User) error {
	return nil
}
