package activity

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetPropositionsByPolitician(ctx context.Context, id string, limit, offset int, onlyVoted bool) ([]Proposition, error)
	GetVotesByPolitician(ctx context.Context, id string, limit, offset int, voteCode *string) ([]Vote, error)
	ListAllPropositions(ctx context.Context, limit, offset int) ([]Proposition, error)
	ListAllVotes(ctx context.Context, limit, offset int) ([]Vote, error)
	GetVoteType(ctx context.Context, house, voteCode string) (*VoteType, error)
	GetVotacao(ctx context.Context, house, idVote string) (*Votacao, error)
	GetPropositionsByVote(ctx context.Context, house, idVote string) ([]Proposition, error)
	GetVotesByProposition(ctx context.Context, house, idProp string) ([]Vote, error)
	GetVotacoesByProposition(ctx context.Context, house, idProp string) ([]Votacao, error)
	GetVotesByVotacao(ctx context.Context, house, idVote string, limit, offset int, voteCode *string) ([]Vote, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

const propositionVoteExistsSQL = `EXISTS (
	SELECT 1 FROM public.votacao_proposicao vp
	WHERE vp.casa = public.proposicao.casa
	  AND vp.id_proposicao = public.proposicao.id_proposicao
)`

func (r *repository) GetPropositionsByPolitician(ctx context.Context, id string, limit, offset int, onlyVoted bool) ([]Proposition, error) {
	var propositions []Proposition
	query := r.db.WithContext(ctx).
		Table("public.proposicao").
		Select("public.proposicao.*, "+propositionVoteExistsSQL+" AS tem_votacao").
		Joins("JOIN public.autoria ON public.proposicao.id_proposicao = public.autoria.id_proposicao AND public.proposicao.casa = public.autoria.casa").
		Where("public.autoria.id_autor = ? AND public.autoria.principal = ?", id, true)

	if onlyVoted {
		query = query.Where(propositionVoteExistsSQL)
	}

	err := query.
		Order("public.proposicao.data_apresentacao DESC NULLS LAST, public.proposicao.id_proposicao DESC").
		Limit(limit).
		Offset(offset).
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
	err := r.db.WithContext(ctx).
		Table("public.proposicao").
		Select("public.proposicao.*, " + propositionVoteExistsSQL + " AS tem_votacao").
		Limit(limit).
		Offset(offset).
		Find(&propositions).Error
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
		Select("public.proposicao.*, "+propositionVoteExistsSQL+" AS tem_votacao").
		Joins("JOIN public.votacao_proposicao ON public.proposicao.id_proposicao = public.votacao_proposicao.id_proposicao AND public.proposicao.casa = public.votacao_proposicao.casa").
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
		Select("public.voto.*").
		Joins("JOIN public.votacao ON public.voto.id_votacao = public.votacao.id_votacao AND public.voto.casa = public.votacao.casa").
		Joins("JOIN public.votacao_proposicao ON public.votacao.id_votacao = public.votacao_proposicao.id_votacao AND public.votacao.casa = public.votacao_proposicao.casa").
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
		Select("public.votacao.*").
		Joins("JOIN public.votacao_proposicao ON public.votacao.id_votacao = public.votacao_proposicao.id_votacao AND public.votacao.casa = public.votacao_proposicao.casa").
		Where("public.votacao_proposicao.casa = ? AND public.votacao_proposicao.id_proposicao = ?", house, idProp).
		Order("public.votacao.data_votacao DESC").
		Find(&votacoes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching votacoes by proposition: %w", err)
	}
	return votacoes, nil
}

func (r *repository) GetVotesByVotacao(ctx context.Context, house, idVote string, limit, offset int, voteCode *string) ([]Vote, error) {
	var votes []Vote
	query := r.db.WithContext(ctx).
		Where("casa = ? AND id_votacao = ?", house, idVote)

	if voteCode != nil {
		query = query.Where("codigo_voto = ?", *voteCode)
	}

	err := query.Order("id_parlamentar").Limit(limit).Offset(offset).Find(&votes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching votes by votacao: %w", err)
	}
	return votes, nil
}
