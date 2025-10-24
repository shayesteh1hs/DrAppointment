package medical

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainMedical "github.com/shayesteh1hs/DrAppointment/internal/domain/medical"
	medicalFilter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

// Mock repository
type MockDoctorRepository struct {
	mock.Mock
}

func (m *MockDoctorRepository) GetAllPaginated(ctx context.Context, filters medicalFilter.DoctorQueryParam, paginator *pagination.LimitOffsetPaginator[domainMedical.Doctor]) ([]domainMedical.Doctor, error) {
	args := m.Called(ctx, filters, paginator)
	var out []domainMedical.Doctor
	if v := args.Get(0); v != nil {
		if cast, ok := v.([]domainMedical.Doctor); ok {
			out = cast
		}
	}
	return out, args.Error(1)
}

func (m *MockDoctorRepository) Count(ctx context.Context, filters medicalFilter.DoctorQueryParam) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockDoctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainMedical.Doctor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainMedical.Doctor), args.Error(1)
}

func TestHandler_GetAllPaginated(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Test data
	doctorID1 := uuid.New()
	doctorID2 := uuid.New()
	specialtyID := uuid.New()
	now := time.Now()

	doctors := []domainMedical.Doctor{
		{
			ID:          doctorID1,
			Name:        "Dr. Smith",
			SpecialtyID: specialtyID,
			PhoneNumber: "1234567890",
			AvatarURL:   "avatar1.jpg",
			Description: "Description 1",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          doctorID2,
			Name:        "Dr. Johnson",
			SpecialtyID: specialtyID,
			PhoneNumber: "0987654321",
			AvatarURL:   "avatar2.jpg",
			Description: "Description 2",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Test cases
	tests := []struct {
		name               string
		queryParams        string
		mockSetup          func(*MockDoctorRepository)
		expectedStatusCode int
		expectedItemCount  int
	}{
		{
			name:        "Success - Get all doctors",
			queryParams: "?page=1&limit=10",
			mockSetup: func(repo *MockDoctorRepository) {
				repo.On("Count", mock.Anything, mock.Anything).Return(2, nil)
				repo.On("GetAllPaginated", mock.Anything, mock.Anything, mock.Anything).Return(doctors, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedItemCount:  2,
		},
		{
			name:        "Success - Filter by name",
			queryParams: "?page=1&limit=10&name=Smith",
			mockSetup: func(repo *MockDoctorRepository) {
				repo.On("Count", mock.Anything, mock.Anything).Return(1, nil)
				repo.On("GetAllPaginated", mock.Anything, mock.Anything, mock.Anything).Return(doctors[:1], nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedItemCount:  1,
		},
		{
			name:        "Error - Count fails",
			queryParams: "?page=1&limit=10",
			mockSetup: func(repo *MockDoctorRepository) {
				repo.On("Count", mock.Anything, mock.Anything).Return(0, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedItemCount:  0,
		},
		{
			name:        "Error - GetAllPaginated fails",
			queryParams: "?page=1&limit=10",
			mockSetup: func(repo *MockDoctorRepository) {
				repo.On("Count", mock.Anything, mock.Anything).Return(2, nil)
				repo.On("GetAllPaginated", mock.Anything, mock.Anything, mock.Anything).Return([]domainMedical.Doctor{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedItemCount:  0,
		},
		{
			name:               "Error - Invalid pagination params",
			queryParams:        "?page=0&limit=10",
			mockSetup:          func(repo *MockDoctorRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedItemCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockDoctorRepository)
			tt.mockSetup(mockRepo)
			handler := NewHandler(mockRepo)
			router := gin.New()

			router.GET("/doctors", handler.GetAllPaginated)
			req, err := http.NewRequest(http.MethodGet, "/doctors"+tt.queryParams, nil)
			require.NoError(t, err)

			// Set the RequestURI which is needed for pagination base URL
			req.RequestURI = "/doctors" + tt.queryParams

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code, w.Body.String())

			if tt.expectedStatusCode == http.StatusOK {
				var response pagination.Result[domainMedical.Doctor]
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response.Items, tt.expectedItemCount)
				assert.Equal(t, tt.expectedItemCount, response.TotalCount)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
