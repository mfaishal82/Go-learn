package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo}
}

func (s *ReportService) GetSalesReport() (*models.Report, error) {
	return s.repo.GetReport()
}
