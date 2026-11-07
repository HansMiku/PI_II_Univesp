package activity

import "context"

type Service interface {
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

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPropositionsByPolitician(ctx context.Context, id string, limit, offset int, onlyVoted bool) ([]Proposition, error) {
	return s.repo.GetPropositionsByPolitician(ctx, id, limit, offset, onlyVoted)
}

func (s *service) GetVotesByPolitician(ctx context.Context, id string, limit, offset int, voteCode *string) ([]Vote, error) {
	return s.repo.GetVotesByPolitician(ctx, id, limit, offset, voteCode)
}

func (s *service) ListAllPropositions(ctx context.Context, limit, offset int) ([]Proposition, error) {
	return s.repo.ListAllPropositions(ctx, limit, offset)
}

func (s *service) ListAllVotes(ctx context.Context, limit, offset int) ([]Vote, error) {
	return s.repo.ListAllVotes(ctx, limit, offset)
}

func (s *service) GetVoteType(ctx context.Context, house, voteCode string) (*VoteType, error) {
	return s.repo.GetVoteType(ctx, house, voteCode)
}

func (s *service) GetVotacao(ctx context.Context, house, idVote string) (*Votacao, error) {
	return s.repo.GetVotacao(ctx, house, idVote)
}

func (s *service) GetPropositionsByVote(ctx context.Context, house, idVote string) ([]Proposition, error) {
	return s.repo.GetPropositionsByVote(ctx, house, idVote)
}

func (s *service) GetVotesByProposition(ctx context.Context, house, idProp string) ([]Vote, error) {
	return s.repo.GetVotesByProposition(ctx, house, idProp)
}

func (s *service) GetVotacoesByProposition(ctx context.Context, house, idProp string) ([]Votacao, error) {
	return s.repo.GetVotacoesByProposition(ctx, house, idProp)
}

func (s *service) GetVotesByVotacao(ctx context.Context, house, idVote string, limit, offset int, voteCode *string) ([]Vote, error) {
	return s.repo.GetVotesByVotacao(ctx, house, idVote, limit, offset, voteCode)
}
