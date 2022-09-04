package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/go-gosh/gestful/component/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type _testTaskSuite struct {
	suite.Suite
	db     *gorm.DB
	repo   repo.TaskRepo
	svc    *task
	engine *gin.Engine
}

func (t *_testTaskSuite) SetupTest() {
	var err error
	t.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	t.Require().NoError(err)
	t.db = t.db.Debug()
	t.Require().NoError(t.db.AutoMigrate(&model.Task{}))
	t.repo = repo.NewTaskRepo(t.db)
	t.svc = NewTask(t.repo).(*task)
	t.engine = gin.Default()
	service.RegisterGroupRoute[TaskViewResp, mapper.CRUDPageResult[TaskViewResp]](t.engine.Group("/api/v1"), "task", t.svc)
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
	req, _ := http.NewRequest("GET", "/api/v1/task?parent_id=0", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(1, body.Page)
	t.EqualValues(service.DefaultPageLimit, body.PageSize)
	t.EqualValues(0, body.Total)
	t.EqualValues(0, body.TotalPage)
	t.Len(body.Data, 0)
}

func (t *_testTaskSuite) Test_Paginate_DefaultTaskWhenJustFillData() {
	data := t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(1, body.Page)
	t.EqualValues(service.DefaultPageLimit, body.PageSize)
	t.EqualValues(10, body.Total)
	t.EqualValues(1, body.TotalPage)
	t.Len(body.Data, service.DefaultPageLimit)
	view, err := Map(data, t.svc.NewTaskViewResp)
	t.Require().NoError(err)
	for i, v := range body.Data {
		taskStr, err := json.Marshal(view[i])
		t.NoError(err)
		t.EqualValues(taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_SubTaskWhenHasData() {
	roots := make([]model.Task, 0, 10)
	for i := 0; i < 10; i++ {
		root := t.addData(i, 0)
		root.SubTask = make([]model.Task, 0, 10)
		for j := 0; j < 9; j++ {
			root.SubTask = append(root.SubTask, t.addData(i*10+j, root.ID))
		}
		sort.Slice(root.SubTask, func(i, j int) bool {
			return root.SubTask[i].ID > root.SubTask[j].ID
		})
		roots = append(roots, root)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?parent_id=31", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(1, body.Page)
	t.EqualValues(service.DefaultPageLimit, body.PageSize)
	t.EqualValues(9, body.Total)
	t.EqualValues(1, body.TotalPage)
	t.Len(body.Data, 9)
	for i, v := range body.Data {
		view, err := t.svc.NewTaskViewResp(&roots[3].SubTask[i])
		t.NoError(err)
		taskStr, err := json.Marshal(view)
		t.NoError(err)
		t.EqualValuesf(taskStr, v, "%s - %s", taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_RootTaskWhenHasData() {
	roots := make([]model.Task, 0, 10)
	for i := 0; i < 10; i++ {
		root := t.addData(i, 0)
		root.SubTask = make([]model.Task, 0, 10)
		for j := 0; j < 10; j++ {
			root.SubTask = append(root.SubTask, t.addData(i*10+j, root.ID))
		}
		sort.Slice(root.SubTask, func(i, j int) bool {
			return root.SubTask[i].ID > root.SubTask[j].ID
		})
		roots = append(roots, root)
	}
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].ID > roots[j].ID
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?parent_id=0", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(1, body.Page)
	t.EqualValues(service.DefaultPageLimit, body.PageSize)
	t.EqualValues(10, body.Total)
	t.EqualValues(1, body.TotalPage)
	t.Len(body.Data, service.DefaultPageLimit)
	for i, v := range body.Data {
		view, err := t.svc.NewTaskViewResp(&roots[i])
		t.NoError(err)
		taskStr, err := json.Marshal(view)
		t.NoError(err)
		t.EqualValuesf(taskStr, v, "%s - %s", taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_11RootTaskWhenMoreData() {
	data := t.addRootData(service.DefaultPageLimit + 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?page_size=11", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(1, body.Page)
	t.EqualValues(11, body.PageSize)
	t.EqualValues(12, body.Total)
	t.EqualValues(2, body.TotalPage)
	t.Len(body.Data, 11)
	view, err := Map(data, t.svc.NewTaskViewResp)
	t.Require().NoError(err)
	for i, v := range body.Data {
		taskStr, err := json.Marshal(view[i])
		t.NoError(err)
		t.EqualValues(taskStr, v)
	}
}

func (t *_testTaskSuite) Test_Paginate_NextPageRootTaskWhenMoreData() {
	data := t.addRootData(service.DefaultPageLimit + 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task?page_size=11&page=2", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	body := mapper.CRUDPageResult[json.RawMessage]{}
	t.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	t.EqualValues(2, body.Page)
	t.EqualValues(11, body.PageSize)
	t.EqualValues(12, body.Total)
	t.EqualValues(2, body.TotalPage)
	t.Len(body.Data, 1)
	view, err := t.svc.NewTaskViewResp(&data[11])
	t.NoError(err)
	taskStr, err := json.Marshal(view)
	t.NoError(err)
	t.EqualValues(taskStr, body.Data[0])
}

func (t *_testTaskSuite) Test_Retrieve_ShowTaskWhenHasData() {
	data := t.addRootData(service.DefaultPageLimit)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/task/5", nil)
	t.engine.ServeHTTP(w, req)
	t.Require().Equal(http.StatusOK, w.Code, w.Body.String())
	view, err := t.svc.NewTaskViewResp(&data[5])
	taskStr, err := json.Marshal(view)
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
	old := data[service.DefaultPageLimit-5]
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
		res = append(res, t.addData(i, 0))
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ID > res[j].ID
	})
	return res
}

func (t *_testTaskSuite) addData(no int, parent uint) model.Task {
	task := model.Task{
		ParentId: parent,
		Point:    100,
		IsCheck:  false,
		Star:     2,
		Category: fmt.Sprintf("test-category-%v", no),
		Title:    fmt.Sprintf("test-title-%v", no),
		Detail:   fmt.Sprintf("test-detail-%v", no),
		StartAt:  time.Now(),
	}
	t.NoError(t.db.Create(&task).Error)
	return task
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, &_testTaskSuite{})
}
