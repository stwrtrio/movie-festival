// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repositories/movie_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/stwrtrio/movie-festival/internal/models"
)

// MockMovieRepository is a mock of MovieRepository interface.
type MockMovieRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMovieRepositoryMockRecorder
}

// MockMovieRepositoryMockRecorder is the mock recorder for MockMovieRepository.
type MockMovieRepositoryMockRecorder struct {
	mock *MockMovieRepository
}

// NewMockMovieRepository creates a new mock instance.
func NewMockMovieRepository(ctrl *gomock.Controller) *MockMovieRepository {
	mock := &MockMovieRepository{ctrl: ctrl}
	mock.recorder = &MockMovieRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovieRepository) EXPECT() *MockMovieRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, movie)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockMovieRepositoryMockRecorder) Create(ctx, movie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMovieRepository)(nil).Create), ctx, movie)
}

// CreateVote mocks base method.
func (m *MockMovieRepository) CreateVote(ctx context.Context, userID, movieID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVote", ctx, userID, movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateVote indicates an expected call of CreateVote.
func (mr *MockMovieRepositoryMockRecorder) CreateVote(ctx, userID, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVote", reflect.TypeOf((*MockMovieRepository)(nil).CreateVote), ctx, userID, movieID)
}

// DeleteVote mocks base method.
func (m *MockMovieRepository) DeleteVote(ctx context.Context, voteID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVote", ctx, voteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVote indicates an expected call of DeleteVote.
func (mr *MockMovieRepositoryMockRecorder) DeleteVote(ctx, voteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVote", reflect.TypeOf((*MockMovieRepository)(nil).DeleteVote), ctx, voteID)
}

// FindArtistByMovieID mocks base method.
func (m *MockMovieRepository) FindArtistByMovieID(ctx context.Context, movieID string) (models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindArtistByMovieID", ctx, movieID)
	ret0, _ := ret[0].(models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindArtistByMovieID indicates an expected call of FindArtistByMovieID.
func (mr *MockMovieRepositoryMockRecorder) FindArtistByMovieID(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindArtistByMovieID", reflect.TypeOf((*MockMovieRepository)(nil).FindArtistByMovieID), ctx, movieID)
}

// FindGenreByMovieID mocks base method.
func (m *MockMovieRepository) FindGenreByMovieID(ctx context.Context, movieID string) (models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindGenreByMovieID", ctx, movieID)
	ret0, _ := ret[0].(models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindGenreByMovieID indicates an expected call of FindGenreByMovieID.
func (mr *MockMovieRepositoryMockRecorder) FindGenreByMovieID(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindGenreByMovieID", reflect.TypeOf((*MockMovieRepository)(nil).FindGenreByMovieID), ctx, movieID)
}

// FindMovieByID mocks base method.
func (m *MockMovieRepository) FindMovieByID(ctx context.Context, movieID string) (models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMovieByID", ctx, movieID)
	ret0, _ := ret[0].(models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMovieByID indicates an expected call of FindMovieByID.
func (mr *MockMovieRepositoryMockRecorder) FindMovieByID(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMovieByID", reflect.TypeOf((*MockMovieRepository)(nil).FindMovieByID), ctx, movieID)
}

// GetAllMovies mocks base method.
func (m *MockMovieRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMovies", ctx, limit, offset)
	ret0, _ := ret[0].([]models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMovies indicates an expected call of GetAllMovies.
func (mr *MockMovieRepositoryMockRecorder) GetAllMovies(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMovies", reflect.TypeOf((*MockMovieRepository)(nil).GetAllMovies), ctx, limit, offset)
}

// GetMostViewedGenre mocks base method.
func (m *MockMovieRepository) GetMostViewedGenre(ctx context.Context, page, pageSize int, sortOrder string) ([]models.GenreView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostViewedGenre", ctx, page, pageSize, sortOrder)
	ret0, _ := ret[0].([]models.GenreView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostViewedGenre indicates an expected call of GetMostViewedGenre.
func (mr *MockMovieRepositoryMockRecorder) GetMostViewedGenre(ctx, page, pageSize, sortOrder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostViewedGenre", reflect.TypeOf((*MockMovieRepository)(nil).GetMostViewedGenre), ctx, page, pageSize, sortOrder)
}

// GetMostViewedMovie mocks base method.
func (m *MockMovieRepository) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostViewedMovie", ctx)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostViewedMovie indicates an expected call of GetMostViewedMovie.
func (mr *MockMovieRepositoryMockRecorder) GetMostViewedMovie(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostViewedMovie", reflect.TypeOf((*MockMovieRepository)(nil).GetMostViewedMovie), ctx)
}

// GetMostVotedMovie mocks base method.
func (m *MockMovieRepository) GetMostVotedMovie(ctx context.Context) (*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostVotedMovie", ctx)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostVotedMovie indicates an expected call of GetMostVotedMovie.
func (mr *MockMovieRepositoryMockRecorder) GetMostVotedMovie(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostVotedMovie", reflect.TypeOf((*MockMovieRepository)(nil).GetMostVotedMovie), ctx)
}

// GetMoviesByIDs mocks base method.
func (m *MockMovieRepository) GetMoviesByIDs(ctx context.Context, movieIDs []string) ([]models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesByIDs", ctx, movieIDs)
	ret0, _ := ret[0].([]models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesByIDs indicates an expected call of GetMoviesByIDs.
func (mr *MockMovieRepositoryMockRecorder) GetMoviesByIDs(ctx, movieIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesByIDs", reflect.TypeOf((*MockMovieRepository)(nil).GetMoviesByIDs), ctx, movieIDs)
}

// GetUserVotedMovieIDs mocks base method.
func (m *MockMovieRepository) GetUserVotedMovieIDs(ctx context.Context, userID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserVotedMovieIDs", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserVotedMovieIDs indicates an expected call of GetUserVotedMovieIDs.
func (mr *MockMovieRepositoryMockRecorder) GetUserVotedMovieIDs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserVotedMovieIDs", reflect.TypeOf((*MockMovieRepository)(nil).GetUserVotedMovieIDs), ctx, userID)
}

// GetVoteByUserAndMovie mocks base method.
func (m *MockMovieRepository) GetVoteByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Vote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVoteByUserAndMovie", ctx, userID, movieID)
	ret0, _ := ret[0].(*models.Vote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVoteByUserAndMovie indicates an expected call of GetVoteByUserAndMovie.
func (mr *MockMovieRepositoryMockRecorder) GetVoteByUserAndMovie(ctx, userID, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVoteByUserAndMovie", reflect.TypeOf((*MockMovieRepository)(nil).GetVoteByUserAndMovie), ctx, userID, movieID)
}

// SearchMovies mocks base method.
func (m *MockMovieRepository) SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchMovies", ctx, query, limit, offset)
	ret0, _ := ret[0].([]models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchMovies indicates an expected call of SearchMovies.
func (mr *MockMovieRepositoryMockRecorder) SearchMovies(ctx, query, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchMovies", reflect.TypeOf((*MockMovieRepository)(nil).SearchMovies), ctx, query, limit, offset)
}

// TrackMovieView mocks base method.
func (m *MockMovieRepository) TrackMovieView(ctx context.Context, movieID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TrackMovieView", ctx, movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// TrackMovieView indicates an expected call of TrackMovieView.
func (mr *MockMovieRepositoryMockRecorder) TrackMovieView(ctx, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackMovieView", reflect.TypeOf((*MockMovieRepository)(nil).TrackMovieView), ctx, movieID)
}

// Update mocks base method.
func (m *MockMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, movie)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockMovieRepositoryMockRecorder) Update(ctx, movie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMovieRepository)(nil).Update), ctx, movie)
}
