package services

import (
	"context"
	"regexp"

	"github.com/Qwepo/rusprofile/gen/rusprof"
	"github.com/Qwepo/rusprofile/internal/parser"

	"google.golang.org/grpc/status"
)

func (s *Services) GetCompany(ctx context.Context, r *rusprof.CreateReqest) (*rusprof.CreateResponse, error) {
	inn := r.GetInn()
	if len(inn) != 10 && len(inn) != 12 {
		s.log.Warn().Str("Servic", "services").Msg("INN is not valid")
		return nil, status.Error(400, "Invalid INN. INN length must be exactly 10 or 12 digits")
	}
	matched, _ := regexp.MatchString(`\D`, inn)
	if matched {
		s.log.Warn().Str("Services", "services").Msg("INN is not valid")
		return nil, status.Error(400, "Invalid INN. The INN must contain only digits")

	}
	data := parser.GoParser(inn)
	if data.INN == "" && data.KPP == "" && data.Leader == "" && data.CompanyName == "" {
		s.log.Info().Str("Services", "services").Msg("Company not found")
		return nil, status.Error(404, "Company not found")
	}
	s.log.Info().Str("Services", "services").Msgf("Company is found, inn: %s", inn)

	return &rusprof.CreateResponse{Inn: inn, Kpp: data.KPP, Name: data.CompanyName, Leader: data.Leader}, nil
}
