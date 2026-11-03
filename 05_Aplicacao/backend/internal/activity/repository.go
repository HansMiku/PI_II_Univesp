package activity

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)
type Repository interface {
	GetPropositionsByPolitician(ctx context.Context, id string, limit, offset int) ([]Proposition, error)
	GetVotesByPolitician(ctx context.Context, id string, limit, offset int, voteCode *string) ([]Vote, error)
	ListAllPropositions(ctx context.Context, limit, offset int) ([]Proposition, error)
	ListAllVotes(ctx context.Context, limit, offset int) ([]Vote, error)
	GetVoteType(ctx context.Context, house, voteCode string) (*VoteType, error)
	GetVotacao(ctx context.Context, house, idVote string) (*Votacao, error)
	GetPropositionsByVote(ctx context.Context, house, idVote string) ([]Proposition, error)
	GetVotesByProposition(ctx context.Context, house, idProp string) ([]Vote, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetPropositionsByPolitician(ctx context.Context, id string, limit, offset int) ([]Proposition, error) {
	var propositions []Proposition
	err := r.db.WithContext(ctx).
		Table("public.proposicao").
		Joins("JOIN public.autoria ON public.proposicao.id_proposicao = public.autoria.id_proposicao").
		Where("public.autoria.id_autor = ?", id).
		Limit(limit).Offset(offset).
		Find(&propositions).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching propositions: %w", err)
	}
	return propositions, nil
}

func (r *repository) GetVotesByPolitician(ctx context.Context, id string, limit, offset int, voteCode *string) ([]Vote, error) {
	var votes []Vote
	query := r.db.WithContext(ctx).Where("id_parlamentar = ?", id)

	if voteCode != nil {
		query = query.Where("codigo_voto = ?", *voteCode)
	}

	err := query.Limit(limit).Offset(offset).Find(&votes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching votes: %w", err)
	}
	return votes, nil
}

func (r *repository) ListAllPropositions(ctx context.Context, limit, offset int) ([]Proposition, error) {
	var propositions []Proposition
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&propositions).Error
	if err != nil {
		return nil, fmt.Errorf("error listing all propositions: %w", err)
	}
	return propositions, nil
}

func (r *repository) ListAllVotes(ctx context.Context, limit, offset int) ([]Vote, error) {
	var votes []Vote
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&votes).Error
	if err != nil {
		return nil, fmt.Errorf("error listing all votes: %w", err)
	}
	return votes, nil
}

func (r *repository) GetVoteType(ctx context.Context, house, voteCode string) (*VoteType, error) {
	var vt VoteType
	err := r.db.WithContext(ctx).Where("casa = ? AND codigo_voto = ?", house, voteCode).First(&vt).Error
	if err != nil {
		return nil, err
	}
	return &vt, nil
}

func (r *repository) GetVotacao(ctx context.Context, house, idVote string) (*Votacao, error) {
	var v Votacao
	err := r.db.WithContext(ctx).Where("casa = ? AND id_votacao = ?", house, idVote).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *repository) GetPropositionsByVote(ctx context.Context, house, idVote string) ([]Proposition, error) {
	var propositions []Proposition
	err := r.db.WithContext(ctx).
		Table("public.proposicao").
		Joins("JOIN public.votacao_proposicao ON public.proposicao.id_proposicao = public.votacao_proposicao.id_proposicao").
		Where("public.votacao_proposicao.casa = ? AND public.votacao_proposicao.id_votacao = ?", house, idVote).
		Find(&propositions).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching propositions by vote: %w", err)
	}
	return propositions, nil
}

func (r *repository) GetVotesByProposition(ctx context.Context, house, idProp string) ([]Vote, error) {
	var votes []Vote
	err := r.db.WithContext(ctx).
		Table("public.voto").
		Joins("JOIN public.votacao ON public.voto.id_votacao = public.votacao.id_votacao").
		Joins("JOIN public.votacao_proposicao ON public.votacao.id_votacao = public.votacao_proposicao.id_votacao").
		Where("public.votacao.casa = ? AND public.votacao_proposicao.id_proposicao = ?", house, idProp).
		Find(&votes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching votes by proposition: %w", err)
	}
	return votes, nil
}

func (r *repository) GetVotacoesByProposition(ctx context.Context, house, idProp string) ([]Votacao, error) {
	var votacoes []Votacao
	err := r.db.WithContext(ctx).
		Table("public.votacao").
		Joins("JOIN public.votacao_proposicao ON public.votacao.id_votacao = public.votacao_proposicao.id_votacao").
		Where("public.votacao_proposicao.casa = ? AND public.votacao_proposicao.id_proposicao = ?", house, idProp).
		Find(&votacoes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching votacoes by proposition: %w", err)
	}
	return votacoes, nil
}
