package corp

import (
	"absurdlab.io/WeSuiteCred/internal/sqlitedb"
	"absurdlab.io/WeSuiteCred/internal/suite"
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	testify_suite "github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	testify_suite.Run(t, new(serviceTestSuite))
}

type serviceTestSuite struct {
	testify_suite.Suite
}

func (s *serviceTestSuite) mockSuiteTicket(db *bun.DB) error {
	ticket := suite.Ticket{ID: 1, Ticket: "mock_ticket"}
	_, err := db.NewInsert().Model(&ticket).Exec(context.TODO())
	return err
}

func (s *serviceTestSuite) mockAuthorization(db *bun.DB) error {
	authz := Authorization{
		ID:            1,
		CorpID:        "foo",
		CorpName:      "foo",
		PermanentCode: "bar",
		AuthInfo:      sqlitedb.WrapJSON(AuthInfo{}),
		Permissions:   sqlitedb.WrapJSON(Permissions{}),
	}

	_, err := db.NewInsert().Model(&authz).Exec(context.Background())

	return err
}

func (s *serviceTestSuite) TestOnNewAuthCode() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.MockGetSuiteAccessTokenEndpoint()
	MockGetCorpAccessTokenEndpoint()
	MockGetPermanentCodeEndpoint()
	MockGetPermissionsEndpoint()

	s.runWithDependencies(
		s.mockSuiteTicket,
		func(service *Service, db *bun.DB) {
			if err := service.OnNewAuthCode(context.TODO(), "mock_auth_code"); s.Assert().NoError(err) {
				var authz Authorization
				if count, err := db.NewSelect().Model(&authz).ScanAndCount(context.TODO()); s.Assert().NoError(err) {
					s.Assert().Equal(1, count)

					s.Assert().NotEmpty(authz.CorpID)
					s.Assert().NotEmpty(authz.CorpName)
					s.Assert().NotEmpty(authz.PermanentCode)

					s.Assert().Len(authz.Permissions.Unwrap().AppPermissions, 1)
				}
			}
		},
	)
}

func (s *serviceTestSuite) TestOnAuthorizationChanged() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.MockGetSuiteAccessTokenEndpoint()
	MockGetCorpAccessTokenEndpoint()
	MockGetAuthorizationInfoEndpoint()
	MockGetPermissionsEndpoint()

	s.runWithDependencies(
		s.mockSuiteTicket,
		s.mockAuthorization,
		func(service *Service, db *bun.DB) {
			if err := service.OnAuthorizationChanged(context.Background(), "foo"); s.Assert().NoError(err) {
				var authz Authorization
				if count, err := db.NewSelect().Model(&authz).ScanAndCount(context.TODO()); s.Assert().NoError(err) {
					s.Assert().Equal(1, count)

					s.Assert().NotEmpty(authz.CorpID)
					s.Assert().NotEmpty(authz.CorpName)
					s.Assert().NotEmpty(authz.PermanentCode)

					s.Assert().Len(authz.Permissions.Unwrap().AppPermissions, 1)
				}
			}
		},
	)
}

func (s *serviceTestSuite) TestOnAuthorizationRemoved() {
	s.runWithDependencies(
		s.mockAuthorization,
		func(service *Service, db *bun.DB) {
			if err := service.OnAuthorizationRemoved(context.TODO(), "foo"); s.Assert().NoError(err) {
				if n, err := db.NewSelect().Model(&Authorization{}).Count(context.Background()); s.Assert().NoError(err) {
					s.Assert().Equal(0, n)
				}
			}
		},
	)
}

func (s *serviceTestSuite) runWithDependencies(fn ...any) {
	logger := zerolog.New(zerolog.NewTestWriter(s.T()))

	err := fx.New(
		fx.NopLogger,
		fx.Supply(
			&logger,
			&suite.Properties{
				Id:                "wwddddccc7775555aaa",
				Secret:            "ldAE_H9anCRN21GKXVfdAAAAAAAAAAAAAAAAAA",
				AccessTokenLeeway: 30 * time.Second,
			},
		),
		fx.Provide(
			sqlitedb.NewMemory,
			NewService,
			suite.NewAccessTokenSupplier,
		),
		fx.Invoke(append([]any{
			func(db *bun.DB) { db.RegisterModel((*suite.Ticket)(nil), (*Authorization)(nil)) },
			sqlitedb.Migrate,
		}, fn...)...),
	).Start(context.TODO())

	s.Require().NoError(err)
}
