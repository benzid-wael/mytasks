package infrastructure

import (
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"testing"
)

func createTempDirectory(prefix string) string {
	dir, err := ioutil.TempDir("/tmp", prefix+"-*")
	if err != nil {
		log.Fatal("Cannot create temp directory: ", err)
		panic(err)
	}
	return dir
}

func assertItemInDir(t *testing.T, key string, item map[string]interface{}, dir string) {
	actual := loadItems(dir)
	expected, _ := tasks.ToMap(item)
	assert.Contains(t, actual, key)
	assert.Equal(t, *expected, actual[key])
}

func Test_GetKey(t *testing.T) {
	tests := []struct {
		value    int
		expected string
	}{
		{value: 1, expected: "1"},
		{value: 12, expected: "12"},
	}
	for _, test := range tests {
		t.Run("GetKey for: "+test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, getKey(test.value))
		})
	}
}

func Test_NewItemRepository_InitializeEmptyDirectories(t *testing.T) {
	// Given
	dir := createTempDirectory("GetItems")
	// When
	NewItemRepository(dir)
	// Then
	assert.DirExists(t, path.Join(dir, "mytasks/repository"))
	assert.DirExists(t, path.Join(dir, "mytasks/archive"))
}

func TestFilesystemItemRepository_GetItems_ReturnsEmptyArray_ForNonInitializedDirectory(t *testing.T) {
	// Given
	dir := createTempDirectory("GetItems")
	testee := NewItemRepository(dir)
	// When
	actual := testee.GetItems()
	// Then
	assert.Empty(t, actual)
}

func TestFilesystemItemRepository_GetItems_ReturnsArrayOfItems(t *testing.T) {
	// Given
	dir := createTempDirectory("GetItems")
	testee := NewItemRepository(dir)
	item := entities.GenerateDummyRawItem()
	key := strconv.Itoa(entities.GetId(item))
	testee.store(key, item, testee.StorageDir)
	// When
	actual := testee.GetItems()
	// Then
	assert.NotEmpty(t, actual)
	assert.Len(t, actual, 1)
}

func Test_storeItems_StoreItemsIntoJsonFile(t *testing.T) {
	// Given
	items := map[string]map[string]interface{}{"1": {"id": 1}}
	dir := createTempDirectory("storeItems")
	NewItemRepository(dir)
	// When
	err := storeItems(items, dir)
	// Then
	assert.FileExists(t, path.Join(dir, "items.json"))
	assert.Nil(t, err)
}

func Test_store_StoreNoteIntoJsonFile(t *testing.T) {
	// Given
	item := entities.GenerateDummyRawItem()
	key := strconv.Itoa(entities.GetId(item))
	dir := createTempDirectory("store")
	testee := NewItemRepository(dir)
	// When
	err := testee.store(key, item, dir)
	// Then
	assert.Nil(t, err)
	assertItemInDir(t, key, item, dir)
}

func Test_FilesystemItemRepository_archiveItem_StoresItemIntoArchiveDirectory(t *testing.T) {
	// Given
	item := entities.GenerateDummyRawItem()
	key := strconv.Itoa(entities.GetId(item))
	dir := createTempDirectory("archiveItem")
	testee := NewItemRepository(dir)
	// When
	err := testee.archiveItem(key, item)
	// Then
	assert.Nil(t, err)
	assertItemInDir(t, key, item, testee.ArchiveDir)
}

func TestFilesystemItemRepository_getItem_ThrowsDoesNotExistError_WhenItemDoesNotExist(t *testing.T) {
	// Given
	dir := createTempDirectory("getItem")
	testee := NewItemRepository(dir)
	// When
	_, err := testee.getItem(entities.GenerateDummyRawItems(1), 12)
	// Then
	assert.Error(t, err)
}

func Test_FilesystemItemRepository_UpdateItem_ThrowsDoesNotExistError_WhenItemDoesNotExist(t *testing.T) {
	// Given
	dir := createTempDirectory("UpdateItem")
	testee := NewItemRepository(dir)
	// When
	err := testee.UpdateItem(2, nil, nil, nil)
	// Then
	assert.Error(t, err)
}

func Test_FilesystemItemRepository_UpdateItem_UpdatesItem_WhenItemExist(t *testing.T) {
	// Given
	dir := createTempDirectory("UpdateItem")
	testee := NewItemRepository(dir)
	item := entities.GenerateDummyRawItem()
	testee.store("1", item, testee.StorageDir)
	title := "My note"
	item["title"] = title
	expected := item
	// When
	err := testee.UpdateItem(1, &title, nil)
	// Then
	assert.Nil(t, err)
	assertItemInDir(t, "1", expected, testee.StorageDir)
}
