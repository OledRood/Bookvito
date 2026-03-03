package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) domain.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *domain.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) GetByID(id uuid.UUID) (*domain.Report, error) {
	var report domain.Report
	err := r.db.Preload("Book").Preload("User").First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetAll(limit, offset int) ([]*domain.Report, error) {
	var reports []*domain.Report
	q := r.db.Preload("Book").Preload("User").Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	return reports, q.Find(&reports).Error
}

func (r *reportRepository) GetPending(limit, offset int) ([]*domain.Report, error) {
	var reports []*domain.Report
	q := r.db.Preload("Book").Preload("User").
		Where("status = ?", domain.ReportPending).
		Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	return reports, q.Find(&reports).Error
}

func (r *reportRepository) UpdateStatus(id uuid.UUID, status domain.ReportStatus) error {
	return r.db.Model(&domain.Report{}).Where("id = ?", id).Update("status", status).Error
}

func (r *reportRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Report{}).Where("status = ?", domain.ReportPending).Count(&count).Error
	return count, err
}
