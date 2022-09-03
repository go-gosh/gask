package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gestful/component/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type _testTaskSuite struct {
	suite.Suite
	db     *gorm.DB
	repo   *repo.TaskRepo
	svc    *Task
	engine *gin.Engine
}

func (t *_testTaskSuite) SetupTest() {
	var err error
	t.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	t.Require().NoError(err)
	t.db = t.db.Debug()
	t.Require().NoError(t.db.AutoMigrate(&model.Task{}))
	t.repo = repo.NewTaskRepo(t.db)
	t.svc = NewTask(t.repo)
	t.engine = gin.Default()
	t.svc.RegisterGroupRoute(t.engine.Group("/api/v1"), "task")
}

func (t *_testTaskSuite) TearDownTest() {
	db, err := t.db.DB()
	t.Require().NoError(err)
	t.Require().NoError(db.Close())
}

func (t *_testTaskSuite) Test_Create() {
	input := map[string]interface{}{
		"parent_id": 0,
		"point":     100,
		"star":      4,
		"category":  "DEV",
		"title":     "test-title",
		"detail":    "",
	}
	w := httptest.NewRecorder()
	var b bytes.Buffer
	bs, err := json.Marshal(input)
	t.Require().NoError(err)
	b.Write(bs)
	req, _ := http.NewRequest("POST", "/api/v1/task", &b)
	t.engine.ServeHTTP(w, req)
	t.Equal(http.StatusOK, w.Code)
	t.Equal(`"success"`, w.Body.String())
	task := model.Task{}
	t.NoError(t.db.Model(&model.Task{}).Find(&task, 1).Error)
	t.EqualValues(input["parent_id"], task.ParentId)
	t.EqualValues(input["point"], task.Point)
	t.EqualValues(false, task.IsCheck)
	t.EqualValues(input["star"], task.Star)
	t.EqualValues(input["category"], task.Category)
	t.EqualValues(input["title"], task.Title)
	t.EqualValues(input["detail"], task.Detail)
	t.Nil(task.DeadLine)
}

func (t *_testTaskSuite) Test_Paginate_NoRootTask() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := struct {
		StartId int               `json:"start_id"`
		Limit   int               `json:"limit"`
		More    bool              `json:"more"`
		Data    []json.RawMessage `json:"data"`
	}{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(0, body.StartId)
	t.EqualValues(service.DefaultPageLimit, body.Limit)
	t.EqualValues(false, body.More)
	t.Len(body.Data, 0)
}

func (t *_testTaskSuite) Test_Paginate_DefaultRootTaskWhenJustFillData() {
	data := t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := struct {
		StartId int               `json:"start_id"`
		Limit   int               `json:"limit"`
		More    bool              `json:"more"`
		Data    []json.RawMessage `json:"data"`
	}{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(0, body.StartId)
	t.EqualValues(service.DefaultPageLimit, body.Limit)
	t.EqualValues(false, body.More)
	t.Len(body.Data, service.DefaultPageLimit)
	for i, v := range body.Data {
		taskStr, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_11RootTaskWhenMoreData() {
	data := t.addRootData(service.DefaultPageLimit + 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?limit=11", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := struct {
		StartId int               `json:"start_id"`
		Limit   int               `json:"limit"`
		More    bool              `json:"more"`
		Data    []json.RawMessage `json:"data"`
	}{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(0, body.StartId)
	t.EqualValues(11, body.Limit)
	t.EqualValues(true, body.More)
	t.Len(body.Data, 11)
	for i, v := range body.Data {
		taskStr, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_NextPageRootTaskWhenMoreData() {
	data := t.addRootData(service.DefaultPageLimit + 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?limit=11&start_id=11", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := struct {
		StartId int               `json:"start_id"`
		Limit   int               `json:"limit"`
		More    bool              `json:"more"`
		Data    []json.RawMessage `json:"data"`
	}{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(11, body.StartId)
	t.EqualValues(11, body.Limit)
	t.EqualValues(false, body.More)
	t.Len(body.Data, 1)
	taskStr, err := json.Marshal(data[11])
	t.NoError(err)
	t.EqualValues(taskStr, body.Data[0])
}

func (t *_testTaskSuite) Test_Retrieve_ShowTaskWhenHasData() {
	data := t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task/5", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	taskStr, err := json.Marshal(data[4])
	t.NoError(err)
	t.EqualValues(taskStr, w.Body.String())
}

func (t *_testTaskSuite) Test_Update_ModifiedWhenChangeData() {
	data := t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	var b bytes.Buffer
	b.WriteString(`{"title":"changed-title"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/task/5", &b)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	task := model.Task{}
	t.NoError(t.db.First(&task, 5).Error)
	t.EqualValues("changed-title", task.Title)
	ac, err := json.Marshal(task)
	t.NoError(err)
	old := data[4]
	old.Title = task.Title
	old.UpdatedAt = task.UpdatedAt
	ex, err := json.Marshal(old)
	t.EqualValues(ex, ac)
}

func (t *_testTaskSuite) Test_Update_NotFoundWhenNoData() {
	t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	var b bytes.Buffer
	b.WriteString(`{"title":"changed-title"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/task/12", &b)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusNotFound, w.Code, w.Body.String())
}

func (t *_testTaskSuite) Test_Retrieve_ShowTaskWhenNoData() {
	t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task/12", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusNotFound, w.Code)
}

func (t *_testTaskSuite) Test_Delete_NotFoundWhenNoData() {
	t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/task/12", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusNotFound, w.Code)
}

func (t *_testTaskSuite) Test_Delete_SuccessWhenHasData() {
	t.addRootData(service.DefaultPageLimit)
	task := model.Task{}
	t.NoError(t.db.First(&task, 1).Error)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/task/1", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code)
	t.ErrorIs(t.db.First(&task, 1).Error, gorm.ErrRecordNotFound)
}

func (t *_testTaskSuite) addRootData(num int) []model.Task {
	res := make([]model.Task, 0, num)
	for i := 0; i < num; i++ {
		task := model.Task{
			ParentId: 0,
			Point:    100,
			IsCheck:  false,
			Star:     2,
			Category: fmt.Sprintf("test-category-%v", i),
			Title:    fmt.Sprintf("test-title-%v", i),
			Detail:   fmt.Sprintf("test-detail-%v", i),
			StartAt:  time.Now(),
		}
		t.NoError(t.db.Create(&task).Error)
		res = append(res, task)
	}
	return res
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, &_testTaskSuite{})
}
