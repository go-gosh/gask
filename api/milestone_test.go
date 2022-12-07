package api

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/service"
)

type _testMilestoneApiSuite struct {
	suite.Suite
	router    *gin.Engine
	timestamp time.Time
	initData  []model.Milestone
}

func (s *_testMilestoneApiSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	s.Require().NoError(err)
	s.Require().NoError(db.AutoMigrate(&model.Milestone{}, &model.Checkpoint{}))
	milestone := service.NewMilestone(db)
	checkpoint := service.NewCheckpoint(db)
	milestoneTag := service.NewMilestoneTag(db)
	s.router = New(milestone, checkpoint, milestoneTag)
	s.timestamp = time.Date(2022, 10, 21, 10, 28, 0, 0, time.Local)
	s.initData = []model.Milestone{
		{
			Model: model.Model{
				CreatedAt: s.timestamp,
				UpdatedAt: s.timestamp,
			},
			Point:     100,
			Title:     "init test data 1",
			StartedAt: s.timestamp,
		},
	}
	s.Require().NoError(db.Create(&s.initData).Error)
}

func (s *_testMilestoneApiSuite) TestMilestonePaginate() {
	page := 1
	limit := 10
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/milestone?page=%d&pageSize=%d&orderBy=id", page, limit), nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(page, jsoniter.Get(w.Body.Bytes(), "data", "page").GetInterface())
	s.EqualValues(limit, jsoniter.Get(w.Body.Bytes(), "data", "pageSize").GetInterface())
	s.EqualValues(len(s.initData), jsoniter.Get(w.Body.Bytes(), "data", "total").GetInterface())
	s.EqualValues(math.Min(10, float64(len(s.initData))), jsoniter.Get(w.Body.Bytes(), "data", "data").Size())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "id").GetInterface())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "point").GetInterface())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "progress").GetInterface())
	s.EqualValues("init test data 1", jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "title").GetInterface())
	s.EqualValues(s.timestamp.Format(time.RFC3339), jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "startedAt").GetInterface())
	s.Nil(jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "deadline").GetInterface())
	s.EqualValues(false, jsoniter.Get(w.Body.Bytes(), "data", "data", 0, "isDeleted").GetInterface())
}

func (s *_testMilestoneApiSuite) TestMilestoneCreate() {
	var b bytes.Buffer
	b.WriteString(`{"title":"test","point":101,"startedAt":"2022-10-21T10:28:00+08:00"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/milestone", &b)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(2, jsoniter.Get(w.Body.Bytes(), "data", "id").ToInt())
	s.EqualValues(101, jsoniter.Get(w.Body.Bytes(), "data", "point").ToInt())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())
	s.EqualValues("test", jsoniter.Get(w.Body.Bytes(), "data", "title").ToString())
	s.EqualValues("2022-10-21T10:28:00+08:00", jsoniter.Get(w.Body.Bytes(), "data", "startedAt").ToString())
	s.EqualValues(false, jsoniter.Get(w.Body.Bytes(), "data", "isDeleted").ToBool())
}

func (s *_testMilestoneApiSuite) TestMilestoneDelete() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(404, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
}

func (s *_testMilestoneApiSuite) TestMilestoneRetrieve() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").ToInt())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "point").GetInterface())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())
	s.EqualValues("init test data 1", jsoniter.Get(w.Body.Bytes(), "data", "title").ToString())
	s.EqualValues(s.timestamp.Format(time.RFC3339), jsoniter.Get(w.Body.Bytes(), "data", "startedAt").ToString())
	s.EqualValues(false, jsoniter.Get(w.Body.Bytes(), "data", "isDeleted").ToBool())
}

func (s *_testMilestoneApiSuite) TestMilestoneUpdate() {
	var b bytes.Buffer
	b.WriteString(`{"title":"update your title","deadline":"2022-11-11T08:00:00+08:00"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/milestone/1", &b)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code, w.Body.String())
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").ToInt())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "point").GetInterface())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())
	s.EqualValues("update your title", jsoniter.Get(w.Body.Bytes(), "data", "title").ToString())
	s.EqualValues(s.timestamp.Format(time.RFC3339), jsoniter.Get(w.Body.Bytes(), "data", "startedAt").ToString())
	s.EqualValues("2022-11-11T08:00:00+08:00", jsoniter.Get(w.Body.Bytes(), "data", "deadline").ToString())
	s.EqualValues(false, jsoniter.Get(w.Body.Bytes(), "data", "isDeleted").ToBool())
}

func TestMilestoneApi(t *testing.T) {
	suite.Run(t, &_testMilestoneApiSuite{})
}
