package service

import (
	"bug-bounty-lite/internal/domain"
)

type organizationService struct {
	repo domain.OrganizationRepository
}

func NewOrganizationService(repo domain.OrganizationRepository) domain.OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) CreateOrganization(name string, description string) (*domain.Organization, error) {
	org := &domain.Organization{
		Name:        name,
		Description: description,
	}
	if err := s.repo.Create(org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *organizationService) GetOrganization(id uint) (*domain.Organization, error) {
	return s.repo.FindByID(id)
}

func (s *organizationService) ListOrganizations() ([]domain.Organization, error) {
	return s.repo.List()
}

func (s *organizationService) UpdateOrganization(id uint, name string, description string) (*domain.Organization, error) {
	org, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		org.Name = name
	}
	if description != "" {
		org.Description = description
	}
	if err := s.repo.Update(org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *organizationService) DeleteOrganization(id uint) error {
	return s.repo.Delete(id)
}
