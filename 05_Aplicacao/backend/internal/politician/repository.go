package politician

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetPolitician(ctx context.Context, id string, house string) (*Politician, error)
	ListPoliticians(ctx context.Context, limit, offset int) ([]Politician, error)
	SearchPoliticiansByName(ctx context.Context, name string, limit, offset int) ([]Politician, error)
	ListHouses(ctx context.Context) ([]House, error)
	ListStates(ctx context.Context) ([]State, error)
	ListParties(ctx context.Context) ([]Party, error)
	ListLegislatures(ctx context.Context) ([]Legislature, error)
	GetObservedRegistration(ctx context.Context, house, idPolitician string) (*ObservedRegistration, error)
	GetPoliticianLegislatures(ctx context.Context, house, idPolitician string) ([]PoliticianLegislature, error)
	GetState(ctx context.Context, uf string) (*State, error)
	GetParty(ctx context.Context, sigla string) (*Party, error)
	GetLegislature(ctx context.Context, number int) (*Legislature, error)
	ListPoliticiansByState(ctx context.Context, uf string, house *string, limit, offset int) ([]Politician, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetPolitician(ctx context.Context, id string, house string) (*Politician, error) {
	var p Politician
	err := r.db.WithContext(ctx).Where("id_parlamentar = ? AND casa = ?", id, house).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("politician not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching politician: %w", err)
	}
	return &p, nil
}

func (r *repository) ListPoliticians(ctx context.Context, limit, offset int) ([]Politician, error) {
	var politicians []Politician
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&politicians).Error
	if err != nil {
		return nil, fmt.Errorf("error listing politicians: %w", err)
	}
	return politicians, nil
}

func (r *repository) SearchPoliticiansByName(ctx context.Context, name string, limit, offset int) ([]Politician, error) {
	var politicians []Politician
	searchTerm := "%" + name + "%"
	err := r.db.WithContext(ctx).
		Where("nome_referencia ILIKE ?", searchTerm).
		Limit(limit).Offset(offset).
		Find(&politicians).Error
	if err != nil {
		return nil, fmt.Errorf("error searching politicians by name: %w", err)
	}
	return politicians, nil
}

func (r *repository) ListHouses(ctx context.Context) ([]House, error) {
	var houses []House
	err := r.db.WithContext(ctx).Find(&houses).Error
	return houses, err
}

func (r *repository) ListStates(ctx context.Context) ([]State, error) {
	var states []State
	err := r.db.WithContext(ctx).Find(&states).Error
	return states, err
}

func (r *repository) ListParties(ctx context.Context) ([]Party, error) {
	var parties []Party
	err := r.db.WithContext(ctx).Find(&parties).Error
	return parties, err
}

func (r *repository) ListLegislatures(ctx context.Context) ([]Legislature, error) {
	var legislatures []Legislature
	err := r.db.WithContext(ctx).Find(&legislatures).Error
	return legislatures, err
}

func (r *repository) GetObservedRegistration(ctx context.Context, house, idPolitician string) (*ObservedRegistration, error) {
	var or ObservedRegistration
	err := r.db.WithContext(ctx).Where("casa = ? AND id_parlamentar = ?", house, idPolitician).First(&or).Error
	if err != nil {
		return nil, err
	}
	return &or, nil
}

func (r *repository) GetPoliticianLegislatures(ctx context.Context, house, idPolitician string) ([]PoliticianLegislature, error) {
	var pl []PoliticianLegislature
	err := r.db.WithContext(ctx).Where("casa = ? AND id_parlamentar = ?", house, idPolitician).Find(&pl).Error
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func (r *repository) GetState(ctx context.Context, uf string) (*State, error) {
	var s State
	err := r.db.WithContext(ctx).Where("uf = ?", uf).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetParty(ctx context.Context, sigla string) (*Party, error) {
	var p Party
	err := r.db.WithContext(ctx).Where("sigla = ?", sigla).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) GetLegislature(ctx context.Context, number int) (*Legislature, error) {
	var l Legislature
	err := r.db.WithContext(ctx).Where("numero_legislatura = ?", number).First(&l).Error
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *repository) ListPoliticiansByState(ctx context.Context, uf string, house *string, limit, offset int) ([]Politician, error) {
	registrations := r.db.WithContext(ctx).
		Model(&ObservedRegistration{}).
		Select("casa, id_parlamentar").
		Where("uf = ?", uf)

	if house != nil {
		registrations = registrations.Where("casa = ?", *house)
	}

	var politicians []Politician
	err := r.db.WithContext(ctx).
		Where("(casa, id_parlamentar) IN (?)", registrations).
		Limit(limit).Offset(offset).
		Find(&politicians).Error
	if err != nil {
		return nil, fmt.Errorf("error listing politicians by state: %w", err)
	}
	return politicians, nil
}
