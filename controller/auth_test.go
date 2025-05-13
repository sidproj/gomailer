package controller

import (
	"gomailer/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mock user model ---
type mockUserModel struct {
	users []models.UserSchema
	err   error
}

func (m *mockUserModel) Find(filter interface{}) ([]models.UserSchema, error) {
	return m.users, m.err
}

type LoginPOST struct{
	name              string
	form              url.Values
	mockUsers         []models.UserSchema
	mockFindErr       error
	expectedStatus    int
	expectedInBody    string
	expectRedirect    bool
	redirectLocation  string
}


func TestLoginControllerPOST(t *testing.T){

	testCases := []LoginPOST{
        {
			name:           "empty input",
			form:           url.Values{},
			expectedStatus: http.StatusOK,
			expectedInBody: "invalid request data",
		},
		{
			name:           "user not found",
			form:           url.Values{"Email": {"a@b.com"}, "Password": {"123"}},
			mockUsers:      []models.UserSchema{},
			expectedStatus: http.StatusOK,
			expectedInBody: "no user found",
		},
		{
			name:           "incorrect password",
			form:           url.Values{"Email": {"a@b.com"}, "Password": {"wrong"}},
			mockUsers:      []models.UserSchema{{Password: "correctpassword"}},
			expectedStatus: http.StatusOK,
			expectedInBody: "incorrect password",
		},
		{
			name:             "successful login",
			form:             url.Values{"Email": {"a@b.com"}, "Password": {"correctpassword"}},
			mockUsers:        []models.UserSchema{{
				ID: primitive.NewObjectID(), 
				Password: "correctpassword",
			}},
			expectedStatus:   http.StatusSeeOther,
			expectRedirect:   true,
			redirectLocation: "/",
		},
    }
	
	for _, tc := range testCases{
		t.Run(tc.name,func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost,"/login",strings.NewReader(tc.form.Encode()))
			req.Header.Set("Content-type","application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			LoginControllerPOST(w,req)

			res:= w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus{
				t.Errorf("expected status code %d, got code %d",tc.expectedStatus,res.StatusCode)
			}

			body := w.Body.String()

			if !strings.Contains(body,tc.expectedInBody){
				t.Errorf("expected response body %q, got %q",tc.expectedInBody,body)
			}
		})
	}
}