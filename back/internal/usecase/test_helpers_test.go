package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeUserRepo struct {
	usersByID    map[uuid.UUID]*domain.User
	usersByEmail map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		usersByID:    map[uuid.UUID]*domain.User{},
		usersByEmail: map[string]*domain.User{},
	}
}

func (r *fakeUserRepo) Create(user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	copyUser := *user
	r.usersByID[user.ID] = &copyUser
	r.usersByEmail[user.Email] = &copyUser
	return nil
}

func (r *fakeUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyUser := *user
	return &copyUser, nil
}

func (r *fakeUserRepo) GetByEmail(email string) (*domain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyUser := *user
	return &copyUser, nil
}

func (r *fakeUserRepo) Update(user *domain.User) error {
	if _, ok := r.usersByID[user.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	copyUser := *user
	r.usersByID[user.ID] = &copyUser
	r.usersByEmail[user.Email] = &copyUser
	return nil
}

func (r *fakeUserRepo) Delete(id uuid.UUID) error {
	user, ok := r.usersByID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	delete(r.usersByEmail, user.Email)
	delete(r.usersByID, id)
	return nil
}

func (r *fakeUserRepo) List(limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range r.usersByID {
		copyUser := *user
		users = append(users, &copyUser)
	}
	return users, nil
}

func (r *fakeUserRepo) GetByRefreshToken(refreshToken string) (*domain.User, error) {
	for _, user := range r.usersByID {
		if user.RefreshToken == refreshToken {
			copyUser := *user
			return &copyUser, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

type fakeMovementRepo struct {
	items []*domain.BookMovementHistory
}

func (r *fakeMovementRepo) Create(m *domain.BookMovementHistory) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	r.items = append(r.items, m)
	return nil
}
func (r *fakeMovementRepo) Update(m *domain.BookMovementHistory) error { return nil }
func (r *fakeMovementRepo) GetByID(id uuid.UUID) (*domain.BookMovementHistory, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeMovementRepo) GetByBookID(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	var out []*domain.BookMovementHistory
	for _, item := range r.items {
		if item.BookID == bookID {
			out = append(out, item)
		}
	}
	return out, nil
}
func (r *fakeMovementRepo) GetByExchangeID(exchangeID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	return nil, nil
}
func (r *fakeMovementRepo) GetByUserID(userID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	return r.items, nil
}
func (r *fakeMovementRepo) List(limit, offset int) ([]*domain.BookMovementHistory, error) {
	return r.items, nil
}

type fakeBookRepo struct {
	books map[uuid.UUID]*domain.Book
}

func newFakeBookRepo() *fakeBookRepo {
	return &fakeBookRepo{books: map[uuid.UUID]*domain.Book{}}
}

func (r *fakeBookRepo) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}
	copyBook := *book
	r.books[book.ID] = &copyBook
	return nil
}
func (r *fakeBookRepo) GetByID(id uuid.UUID) (*domain.Book, error) {
	book, ok := r.books[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyBook := *book
	return &copyBook, nil
}
func (r *fakeBookRepo) Update(book *domain.Book) error {
	if _, ok := r.books[book.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	copyBook := *book
	r.books[book.ID] = &copyBook
	return nil
}
func (r *fakeBookRepo) Delete(bookID uuid.UUID) error {
	delete(r.books, bookID)
	return nil
}
func (r *fakeBookRepo) List(limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	for _, book := range r.books {
		copyBook := *book
		books = append(books, &copyBook)
	}
	return books, nil
}
func (r *fakeBookRepo) ListFiltered(filter domain.BookListFilter) (*domain.BookListResponse, error) {
	books, _ := r.List(0, 0)
	return &domain.BookListResponse{Items: books, HasMore: false}, nil
}
func (r *fakeBookRepo) GetSummaryList(limit, offset int, excludeUserID *uuid.UUID) ([]*domain.BookSummary, error) {
	return nil, nil
}
func (r *fakeBookRepo) GetByStatus(status domain.BookStatus, limit, offset int) ([]*domain.Book, error) {
	return nil, nil
}
func (r *fakeBookRepo) GetByLocationID(locationID uuid.UUID) ([]*domain.Book, error) {
	return nil, nil
}
func (r *fakeBookRepo) GetByOwner(ownerID uuid.UUID) ([]*domain.Book, error) {
	var books []*domain.Book
	for _, book := range r.books {
		if book.OwnerID == ownerID {
			copyBook := *book
			books = append(books, &copyBook)
		}
	}
	return books, nil
}

type fakeExchangeRepo struct {
	items []*domain.Exchange
}

func (r *fakeExchangeRepo) Create(exchange *domain.Exchange) error {
	if exchange.ID == uuid.Nil {
		exchange.ID = uuid.New()
	}
	copyExchange := *exchange
	r.items = append(r.items, &copyExchange)
	return nil
}
func (r *fakeExchangeRepo) GetByID(id uuid.UUID) (*domain.Exchange, error) {
	for _, item := range r.items {
		if item.ID == id {
			copyExchange := *item
			return &copyExchange, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeExchangeRepo) GetByUserID(userID uuid.UUID) ([]*domain.Exchange, error) {
	var out []*domain.Exchange
	for _, item := range r.items {
		if item.UserID == userID {
			copyExchange := *item
			out = append(out, &copyExchange)
		}
	}
	return out, nil
}
func (r *fakeExchangeRepo) GetByBookID(bookID uuid.UUID) ([]*domain.Exchange, error) {
	var out []*domain.Exchange
	for _, item := range r.items {
		if item.BookID == bookID {
			copyExchange := *item
			out = append(out, &copyExchange)
		}
	}
	return out, nil
}
func (r *fakeExchangeRepo) Update(exchange *domain.Exchange) error {
	for i, item := range r.items {
		if item.ID == exchange.ID {
			copyExchange := *exchange
			r.items[i] = &copyExchange
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}
func (r *fakeExchangeRepo) Delete(id uuid.UUID) error { return nil }
func (r *fakeExchangeRepo) List(limit, offset int) ([]*domain.Exchange, error) {
	return r.items, nil
}
func (r *fakeExchangeRepo) GetExpired() ([]*domain.Exchange, error) {
	return nil, nil
}

type fakeLocationRepo struct {
	locations map[uuid.UUID]*domain.Location
}

func newFakeLocationRepo() *fakeLocationRepo {
	return &fakeLocationRepo{locations: map[uuid.UUID]*domain.Location{}}
}

func (r *fakeLocationRepo) Create(location *domain.Location) error {
	if location.ID == uuid.Nil {
		location.ID = uuid.New()
	}
	copyLocation := *location
	r.locations[location.ID] = &copyLocation
	return nil
}
func (r *fakeLocationRepo) GetByID(id uuid.UUID) (*domain.Location, error) {
	location, ok := r.locations[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyLocation := *location
	return &copyLocation, nil
}
func (r *fakeLocationRepo) GetByAddress(address string) (*domain.Location, error) { return nil, nil }
func (r *fakeLocationRepo) GetAll() ([]domain.Location, error) {
	var out []domain.Location
	for _, location := range r.locations {
		out = append(out, *location)
	}
	return out, nil
}
func (r *fakeLocationRepo) Update(location *domain.Location) error {
	copyLocation := *location
	r.locations[location.ID] = &copyLocation
	return nil
}
func (r *fakeLocationRepo) Delete(id uuid.UUID) error {
	delete(r.locations, id)
	return nil
}

type fakeReportRepo struct {
	reports []*domain.Report
}

func (r *fakeReportRepo) Create(report *domain.Report) error {
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}
	copyReport := *report
	r.reports = append(r.reports, &copyReport)
	return nil
}
func (r *fakeReportRepo) GetByID(id uuid.UUID) (*domain.Report, error)       { return nil, nil }
func (r *fakeReportRepo) GetAll(limit, offset int) ([]*domain.Report, error) { return r.reports, nil }
func (r *fakeReportRepo) GetPending(limit, offset int) ([]*domain.Report, error) {
	var out []*domain.Report
	for _, report := range r.reports {
		if report.Status == domain.ReportPending {
			out = append(out, report)
		}
	}
	return out, nil
}
func (r *fakeReportRepo) UpdateStatus(id uuid.UUID, status domain.ReportStatus) error {
	for _, report := range r.reports {
		if report.ID == id {
			report.Status = status
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}
func (r *fakeReportRepo) CountPending() (int64, error) {
	var total int64
	for _, report := range r.reports {
		if report.Status == domain.ReportPending {
			total++
		}
	}
	return total, nil
}

type fakeImageStorage struct {
	saved   map[string]string
	deleted []string
}

func newFakeImageStorage() *fakeImageStorage {
	return &fakeImageStorage{saved: map[string]string{}}
}

func (s *fakeImageStorage) Save(filename, contentType string, reader io.Reader) (string, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	url := s.PublicURL(filename)
	s.saved[url] = string(body)
	return url, nil
}

func (s *fakeImageStorage) Exists(publicURL string) (bool, error) {
	_, ok := s.saved[publicURL]
	return ok, nil
}

func (s *fakeImageStorage) Delete(publicURL string) error {
	s.deleted = append(s.deleted, publicURL)
	delete(s.saved, publicURL)
	return nil
}

func (s *fakeImageStorage) PublicURL(filename string) string {
	return "/images/" + strings.TrimLeft(filename, "/")
}

type fakeMetadataProvider struct {
	searchFn func(query string) (*domain.BookMeta, error)
}

func (p *fakeMetadataProvider) Search(query string) (*domain.BookMeta, error) {
	if p.searchFn != nil {
		return p.searchFn(query)
	}
	return &domain.BookMeta{Title: query}, nil
}

func mustUUID() uuid.UUID {
	return uuid.New()
}

func fixedTime(hours int) time.Time {
	return time.Now().Add(time.Duration(hours) * time.Hour)
}

func hasErrorCode(err error, code domain.ErrorCode) bool {
	actual, ok := domain.AppErrorCode(err)
	return ok && actual == code
}

func errContains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}

var errBoom = errors.New("boom")
