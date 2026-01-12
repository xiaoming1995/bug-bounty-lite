package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type organizationRepo struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) domain.OrganizationRepository {
	return &organizationRepo{db: db}
}

func (r *organizationRepo) Create(org *domain.Organization) error {
	return r.db.Create(org).Error
}

func (r *organizationRepo) FindByID(id uint) (*domain.Organization, error) {
	var org domain.Organization
	err := r.db.First(&org, id).Error
	return &org, err
}

func (r *organizationRepo) List() ([]domain.Organization, error) {
	var orgs []domain.Organization
	err := r.db.Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepo) Update(org *domain.Organization) error {
	return r.db.Save(org).Error
}

func (r *organizationRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Organization{}, id).Error
}
