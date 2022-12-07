package api

import (
	"bytes"
	"fmt"
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

type _testCheckpointApiSuite struct {
	suite.Suite
	router    *gin.Engine
	timestamp time.Time
	initData  []model.Milestone
}

func (s *_testCheckpointApiSuite) SetupTest() {
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
			Title:     "init test data 2",
			StartedAt: s.timestamp,
		},
	}
	s.Require().NoError(db.Create(&s.initData).Error)
}

func (s *_testCheckpointApiSuite) TestCheckpointPaginate() {
	page := 1
	limit := 10
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/checkpoint?page=%d&pageSize=%d&orderBy=id", page, limit), nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.T().Log(w.Body.String())
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(page, jsoniter.Get(w.Body.Bytes(), "data", "page").GetInterface())
	s.EqualValues(limit, jsoniter.Get(w.Body.Bytes(), "data", "pageSize").GetInterface())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "total").GetInterface())
}

func (s *_testCheckpointApiSuite) TestCheckpointCreate_IsChecked() {
	var b bytes.Buffer
	b.WriteString(`{"content":"test","point":100,"joinedAt":"2022-10-21T10:28:00+08:00","checkedAt":"2022-10-21T10:28:01+08:00","milestoneId":1}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/checkpoint", &b)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.T().Log(w.Body.String())
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").GetInterface())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "point").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "milestoneId").GetInterface())
	s.EqualValues("test", jsoniter.Get(w.Body.Bytes(), "data", "content").GetInterface())
	s.EqualValues("2022-10-21T10:28:00+08:00", jsoniter.Get(w.Body.Bytes(), "data", "joinedAt").GetInterface())
	s.EqualValues("2022-10-21T10:28:01+08:00", jsoniter.Get(w.Body.Bytes(), "data", "checkedAt").GetInterface())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").GetInterface())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/checkpoint/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").GetInterface())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "point").GetInterface())
	s.EqualValues("test", jsoniter.Get(w.Body.Bytes(), "data", "content").GetInterface())
	s.EqualValues("2022-10-21T10:28:00+08:00", jsoniter.Get(w.Body.Bytes(), "data", "joinedAt").GetInterface())
	s.EqualValues("2022-10-21T10:28:01+08:00", jsoniter.Get(w.Body.Bytes(), "data", "checkedAt").GetInterface())
}

func (s *_testCheckpointApiSuite) TestMilestoneDelete() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/checkpoint/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/checkpoint/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(404, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
}

func (s *_testCheckpointApiSuite) TestMilestoneRetrieve_NotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/checkpoint/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(404, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
}

func (s *_testCheckpointApiSuite) TestMilestoneUpdate() {
	var b bytes.Buffer
	b.WriteString(`{"content":"test","point":100,"joinedAt":"2022-10-21T10:28:00+08:00","milestoneId":1}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/checkpoint", &b)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.T().Log(w.Body.String())
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	b.Reset()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").GetInterface())
	s.EqualValues(0, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())

	b.WriteString(`{"content":"update your content","checkedAt":"2022-10-21T10:28:01+08:00","isChecked":true}`)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/checkpoint/1", &b)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code, w.Body.String())
	s.T().Log(w.Body.String())
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/milestone/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").GetInterface())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").GetInterface())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "progress").GetInterface())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/checkpoint/1", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.EqualValues(200, jsoniter.Get(w.Body.Bytes(), "code").ToInt())
	s.EqualValues(1, jsoniter.Get(w.Body.Bytes(), "data", "id").ToInt())
	s.EqualValues(100, jsoniter.Get(w.Body.Bytes(), "data", "point").GetInterface())
	s.EqualValues("update your content", jsoniter.Get(w.Body.Bytes(), "data", "content").ToString())
	s.EqualValues("2022-10-21T10:28:00+08:00", jsoniter.Get(w.Body.Bytes(), "data", "joinedAt").ToString())
	s.EqualValues("2022-10-21T10:28:01+08:00", jsoniter.Get(w.Body.Bytes(), "data", "checkedAt").ToString())
}

func TestCheckpointApi(t *testing.T) {
	suite.Run(t, &_testCheckpointApiSuite{})
}
