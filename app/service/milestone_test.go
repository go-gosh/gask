package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/query"
)

type _testServiceSuite struct {
	suite.Suite
	db  *gorm.DB
	svc *Milestone
}

func (s *_testServiceSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		QueryFields:    true,
	})
	s.Require().NoError(err)
	s.db = s.db.Debug()
	s.Require().NoError(s.db.AutoMigrate(&model.Milestone{}, &model.Checkpoint{}))
	s.svc = NewMilestone(query.Use(s.db))
}

func (s *_testServiceSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	s.Require().NoError(db.Close())
}

func (s *_testServiceSuite) TestCreateMilestone() {
	args := Create{
		Point:     100,
		Title:     "test 1",
		Content:   "test content 1",
		StartedAt: time.Now(),
		Deadline:  nil,
	}
	result, err := s.svc.CreateMilestone(args)
	s.NoError(err)
	s.NotEmpty(result.ID)
}

func (s *_testServiceSuite) TestViewAllMilestone() {
	now := time.Now()
	for i := 0; i < 20; i++ {
		args := Create{
			Point:     100,
			Title:     fmt.Sprintf("test %v", i),
			Content:   fmt.Sprintf("test content %v", i),
			StartedAt: now,
			Deadline:  nil,
		}
		m, _ := s.svc.CreateMilestone(args)
		s.svc.SplitMilestoneById(m.ID, CheckpointCreate{
			Point:     20,
			Content:   "checkpoint20",
			JoinedAt:  now,
			CheckedAt: &now,
		}, CheckpointCreate{
			Point:     30,
			Content:   "checkpoint30",
			JoinedAt:  now,
			CheckedAt: &now,
		})
	}
}

func TestService(t *testing.T) {
	suite.Run(t, &_testServiceSuite{})
}
