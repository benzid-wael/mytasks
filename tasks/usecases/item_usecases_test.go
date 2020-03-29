package usecases

import (
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ItemUseCaseTestSuite struct {
	suite.Suite
	Repo infrastructure.ItemRepository
}

func (suite *ItemUseCaseTestSuite) SetupTest() {
	dir := tasks.CreateTempDirectory("ItemUseCase")
	suite.Repo = infrastructure.NewItemRepository(dir)
}

func (suite *ItemUseCaseTestSuite) Test_ItemUseCase_CreateNote_CreatesNote() {
	// Given
	testee := NewItemUseCase(suite.Repo)
	// When
	actual, err := testee.CreateNote("My note", "@taskbook")
	// Then
	suite.Nil(err)
	suite.NotNil(actual)
}

func (suite *ItemUseCaseTestSuite) Test_ItemUseCase_CreateTask_CreatesTask() {
	// Given
	testee := NewItemUseCase(suite.Repo)
	// When
	actual, err := testee.CreateTask("Learn golang", "@golang")
	// Then
	suite.Nil(err)
	suite.NotNil(actual)
}

func (suite *ItemUseCaseTestSuite) Test_ItemUseCase_EditItem_UpdatesItemTitle() {
	// Given
	testee := NewItemUseCase(suite.Repo)
	task, _ := testee.CreateTask("Learn golang", "@golang")
	// When
	title := "Learn golang in 2 weeks"
	err := testee.EditItem(task.Id, title, "", nil, "@golang")
	// Then
	suite.Nil(err)
	item := suite.Repo.GetItem(task.Id)
	suite.Equal((item).GetTitle(), title)
}

//// In order for 'go test' to run this suite, we need to create
//// a normal test function and pass our suite to suite.Run
func TestItemUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ItemUseCaseTestSuite))
}
