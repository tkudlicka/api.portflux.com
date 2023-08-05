package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tkudlicka/portflux-api/config"
	"github.com/tkudlicka/portflux-api/core/models"
	"github.com/tkudlicka/portflux-api/core/ports"
	"github.com/tkudlicka/portflux-api/test/mocks"
)

// TestCreateBroker_Ok checks that CreateBroker handler returns the expected response when a valid request is received
func TestCreateBroker_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedResponse := models.CreationResp{
		InsertedID: "new-id",
	}
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Create), mock.Anything, mock.AnythingOfType("models.CreateBrokerReq")).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker"
	body := models.CreateBrokerReq{
		Name:        "test",
		Description: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusCreated, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.CreationResp
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestCreateBroker_InvalidRequest checks that CreateBroker handler returns an error when the received request is not valid
func TestCreateBroker_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, nil)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker"
	invalidBody := []byte(`{"Name":invalid-type}`)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(invalidBody))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestCreateBroker_CreateError checks that CreateBroker handler returns an error when the Create function from the service fails
func TestCreateBroker_CreateError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Create), mock.Anything, mock.AnythingOfType("models.CreateBrokerReq")).Return(models.CreationResp{}, fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker"
	body := models.CreateBrokerReq{
		Name:        "test",
		Description: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestCreateManyBrokers_Ok checks that CreateManyBrokers handler returns the expected response when a valid request is received
func TestCreateManyBrokers_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedResponse := models.MultiCreationResp{
		InsertedIDs: []string{"new-id"},
	}
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateBrokerReq")).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker/many"
	body := []models.CreateBrokerReq{
		{
			Name:        "test",
			Description: "test",
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusCreated, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.MultiCreationResp
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestCreateManyBrokers_InvalidRequest checks that CreateManyBrokers handler returns an error when the received request is not valid
func TestCreateManyBrokers_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, nil)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker/many"
	invalidBody := []byte(`[{"Email":invalid-type}]`)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(invalidBody))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestCreateManyBrokers_CreateManyError checks that CreateManyBrokers handler returns an error when the CreateMany function from the service fails
func TestCreateManyBrokers_CreateManyError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateBrokerReq")).Return(models.MultiCreationResp{}, fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker/many"
	body := []models.CreateBrokerReq{
		{
			Name:        "test",
			Description: "test",
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetAllBrokers_Ok checks that GetAllBrokers handler returns the expected response when everything goes as expected
func TestGetAllBrokers_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedResponse := []models.BrokerResp{
		{
			Name:        "test",
			Description: "test",
		},
	}
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetAll), mock.Anything).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response []models.BrokerResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetAllBrokers_GetAllError checks that GetAllBrokers handler returns an error when the GetAll function from the service fails
func TestGetAllBrokers_GetAllError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetAll), mock.Anything).Return([]models.BrokerResp{}, fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/broker"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetBrokerBySlug_Ok checks that GetBrokerBySlug handler returns the expected response when everything goes as expected
func TestGetBrokerBySlug_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedResponse := models.BrokerResp{
		Name:        "test",
		Slug:        "test",
		Description: "test",
	}
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetBySlug), mock.Anything, expectedResponse.Slug).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/slug/%s", expectedResponse.Slug)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.BrokerResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetBrokerBySlug_GetBySlugError checks that GetBrokerBySlug handler returns an error when the GetBySlug function from the service fails
func TestGetBrokerBySlug_GetBySlugError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	testEmail := "test@test.com"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetBySlug), mock.Anything, testEmail).Return(models.BrokerResp{}, fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/slug/%s", testEmail)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetBrokerByID_Ok checks that GetBrokerByID handler returns the expected response when everything goes as expected
func TestGetBrokerByID_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedResponse := models.BrokerResp{
		BrokerID: "test-id",
	}
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetByID), mock.Anything, expectedResponse.BrokerID).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", expectedResponse.BrokerID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.BrokerResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetBrokerByID_GetBrokerByIDError checks that GetBrokerByID handler returns an error when the GetByID function from the service fails
func TestGetBrokerByID_GetByIDError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	testID := "test-id"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.GetByID), mock.Anything, testID).Return(models.BrokerResp{}, fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestUpdateBroker_Ok checks that UpdateBroker handler returns the expected response when a valid request is received
func TestUpdateBroker_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	testID := "test-id"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateBrokerReq")).Return(nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	testName := "test"
	body := models.UpdateBrokerReq{
		Name: testName,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

// TesUpdateBroker_InvalidRequest checks that UpdateBroker handler returns an error when the received request is not valid
func TestUpdateBroker_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, nil)

	rr := httptest.NewRecorder()
	testID := "test-id"
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	invalidBody := []byte(`{"Email":invalid-type}`)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(invalidBody))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestUpdateBroker_UpdateError checks that UpdateBroker handler returns an error when the Update function from the service fails
func TestUpdateBroker_UpdateError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	testID := "test-id"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateBrokerReq")).Return(fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	testName := "test"
	body := models.UpdateBrokerReq{
		Name: testName,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestDeleteBroker_Ok checks that DeleteBroker handler returns the expected response when everything goes as expected
func TestDeleteBroker_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	testID := "test-id"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Delete), mock.Anything, testID).Return(nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	headerName := "Authorization"
	jwtOkAdmin := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.ojgkTERSbf1y8rVLNNF-u70hp_GmJOcsmdB5PLcdews"
	req.Header.Add(headerName, jwtOkAdmin)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

// TestDeleteBroker_DeleteError checks that DeleteBroker handler returns an error when the Delete function from the service fails
func TestDeleteBroker_DeleteError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	brokerService := mocks.NewBrokerService(t)
	expectedError := "service-error"
	testID := "test-id"
	brokerService.On(testutils.FunctionName(t, ports.BrokerService.Delete), mock.Anything, testID).Return(fmt.Errorf(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	SetBrokerRoutes(context.Background(), cfg, r, brokerService)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/broker/%s", testID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	headerName := "Authorization"
	jwtOkAdmin := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.ojgkTERSbf1y8rVLNNF-u70hp_GmJOcsmdB5PLcdews"
	req.Header.Add(headerName, jwtOkAdmin)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}
