package politician

import (
	"context"
	"fmt"
)

type Service interface {
	GetPoliticianFullData(ctx context.Context, id string, house string) (*Politician, error)
	ListAllPoliticians(ctx context.Context, limit, offset int) ([]Politician, error)
	SearchPoliticiansByName(ctx context.Context, name string, limit, offset int) ([]Politician, error)
	ListAllHouses(ctx context.Context) ([]House, error)
	ListAllStates(ctx context.Context) ([]State, error)
	ListAllParties(ctx context.Context) ([]Party, error)
	ListAllLegislatures(ctx context.Context) ([]Legislature, error)
	GetObservedRegistration(ctx context.Context, house, idPolitician string) (*ObservedRegistration, error)
	GetPoliticianLegislatures(ctx context.Context, house, idPolitician string) ([]PoliticianLegislature, error)
	GetState(ctx context.Context, uf string) (*State, error)
	GetParty(ctx context.Context, sigla string) (*Party, error)
	GetLegislature(ctx context.Context, number int) (*Legislature, error)
	ListPoliticiansByState(ctx context.Context, uf string, house *string, limit, offset int) ([]Politician, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPoliticianFullData(ctx context.Context, id string, house string) (*Politician, error) {
	return s.repo.GetPolitician(ctx, id, house)
}

func (s *service) ListAllPoliticians(ctx context.Context, limit, offset int) ([]Politician, error) {
	return s.repo.ListPoliticians(ctx, limit, offset)
}

func (s *service) SearchPoliticiansByName(ctx context.Context, name string, limit, offset int) ([]Politician, error) {
	return s.repo.SearchPoliticiansByName(ctx, name, limit, offset)
}

func (s *service) ListAllHouses(ctx context.Context) ([]House, error) {
	return s.repo.ListHouses(ctx)
}

func (s *service) ListAllStates(ctx context.Context) ([]State, error) {
	return s.repo.ListStates(ctx)
}

func (s *service) ListAllParties(ctx context.Context) ([]Party, error) {
	return s.repo.ListParties(ctx)
}

func (s *service) ListAllLegislatures(ctx context.Context) ([]Legislature, error) {
	return s.repo.ListLegislatures(ctx)
}

func (s *service) GetObservedRegistration(ctx context.Context, house, idPolitician string) (*ObservedRegistration, error) {
	return s.repo.GetObservedRegistration(ctx, house, idPolitician)
}

func (s *service) GetPoliticianLegislatures(ctx context.Context, house, idPolitician string) ([]PoliticianLegislature, error) {
	return s.repo.GetPoliticianLegislatures(ctx, house, idPolitician)
}

func (s *service) GetState(ctx context.Context, uf string) (*State, error) {
	return s.repo.GetState(ctx, uf)
}

func (s *service) GetParty(ctx context.Context, sigla string) (*Party, error) {
	return s.repo.GetParty(ctx, sigla)
}

func (s *service) GetLegislature(ctx context.Context, number int) (*Legislature, error) {
	return s.repo.GetLegislature(ctx, number)
}

func (s *service) ListPoliticiansByState(ctx context.Context, uf string, house *string, limit, offset int) ([]Politician, error) {
	if len(uf) != 2 {
		return nil, fmt.Errorf("uf inválida: %q", uf)
	}
	return s.repo.ListPoliticiansByState(ctx, uf, house, limit, offset)
}
