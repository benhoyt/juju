// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package resources

import (
	"testing"

	"github.com/juju/tc"

	loggertesting "github.com/juju/juju/internal/logger/testing"
)

func TestFacadeSuite(t *testing.T) {
	tc.Run(t, &FacadeSuite{})
}

type FacadeSuite struct {
	BaseSuite
}

func (s *FacadeSuite) TestNewFacadeOkay(c *tc.C) {
	defer s.setupMocks(c).Finish()
	_, err := NewResourcesAPI(s.applicationService, s.resourceService, s.crossModelRelationService, s.factory, loggertesting.WrapCheckLog(c))
	c.Check(err, tc.ErrorIsNil)
}

func (s *FacadeSuite) TestNewFacadeMissingApplicationService(c *tc.C) {
	defer s.setupMocks(c).Finish()
	_, err := NewResourcesAPI(nil, s.resourceService, s.crossModelRelationService, s.factory, loggertesting.WrapCheckLog(c))
	c.Check(err, tc.ErrorMatches, ".*missing application service.*")
}

func (s *FacadeSuite) TestNewFacadeMissingResourceService(c *tc.C) {
	defer s.setupMocks(c).Finish()
	_, err := NewResourcesAPI(s.applicationService, nil, s.crossModelRelationService, s.factory, loggertesting.WrapCheckLog(c))
	c.Check(err, tc.ErrorMatches, ".*missing resource service.*")
}

func (s *FacadeSuite) TestNewFacadeMissingCrossModelRelationService(c *tc.C) {
	defer s.setupMocks(c).Finish()
	_, err := NewResourcesAPI(s.applicationService, s.resourceService, nil, s.factory, loggertesting.WrapCheckLog(c))
	c.Check(err, tc.ErrorMatches, ".*missing cross model relation service.*")
}

func (s *FacadeSuite) TestNewFacadeMissingFactory(c *tc.C) {
	defer s.setupMocks(c).Finish()
	_, err := NewResourcesAPI(s.applicationService, s.resourceService, s.crossModelRelationService, nil, loggertesting.WrapCheckLog(c))
	c.Check(err, tc.ErrorMatches, ".*missing factory for new repository.*")
}
