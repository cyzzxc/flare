package FlareData

import (
	"testing"

	FlareModel "github.com/soulteary/flare/config/model"
)

func TestGetBookmarksDataAsJSON(t *testing.T) {
	categories, bookmarks := GetBookmarksForEditor()
	if len(categories) == 0 || len(bookmarks) == 0 {
		t.Fatal("GetBookmarksForEditor Failed")
	}
}

func TestGetAndUpdateBookmarksFromEditor(t *testing.T) {

	const categories = `1,链接分类1
2,链接分类2
3,链接分类3
4,链接分类4`
	const bookmarks = `1,示例链接,https://link.example.com,[Flare 应用],sticky-note-line,链接描述文本
2,示例链接,https://link.example.com,[Flare 应用],fire-line,链接描述文本
3,示例链接,https://link.example.com,[Flare 应用],mail-line,链接描述文本
4,示例链接,https://link.example.com,[Flare 应用],file-text-line,链接描述文本
5,示例链接,https://link.example.com,[Flare 应用],skull-line,
6,示例链接,https://link.example.com,[Flare 应用],plug-line,
7,示例链接,https://link.example.com,[Flare 应用],image-line,
8,示例链接,https://link.example.com,[Flare 应用],sun-cloudy-line,
9,示例链接,https://link.example.com,链接分类1,checkbox-circle-line,
10,示例链接,https://link.example.com,链接分类1,eraser-line,
11,示例链接,https://link.example.com,链接分类1,mastodon-line,
12,示例链接,https://link.example.com,链接分类1,a-b,
13,示例链接,https://link.example.com,链接分类1,flask-line,
14,示例链接,https://link.example.com,链接分类2,sofa-line,
15,示例链接,https://link.example.com,链接分类2,focus-line,
16,示例链接,https://link.example.com,链接分类2,chat-settings-line,
17,示例链接,https://link.example.com,链接分类2,link,
18,示例链接,https://link.example.com,链接分类2,building-line,
19,示例链接,https://link.example.com,链接分类3,restaurant-line,
20,示例链接,https://link.example.com,链接分类3,keyboard-line,
21,示例链接,https://link.example.com,链接分类3,book-line,
22,示例链接,https://link.example.com,链接分类3,quill-pen-line,
23,示例链接,https://link.example.com,链接分类3,palette-line,
24,示例链接,https://link.example.com,链接分类4,music-2-line,
25,示例链接,https://link.example.com,链接分类4,spy-line,
26,示例链接,https://link.example.com,链接分类4,bookmark-line,
27,示例链接,https://link.example.com,链接分类4,group-line,
28,示例链接,https://link.example.com,链接分类4,seedling-line,`

	updated := UpdateBookmarksFromEditor(categories, bookmarks)
	if !updated {
		t.Fatal("UpdateBookmarksFromEditor Failed")
	}

	bookmarkCategories, ok := getCategoriesFromCSV(categories)
	if ok != nil {
		t.Fatal("getCategoriesFromCSV Failed")
	}

	_, _, ok = getBookmarksFromCSV(bookmarks, bookmarkCategories)
	if ok != nil {
		t.Fatal("getBookmarksFromCSV Failed")
	}
}

func TestPropsRemoveAndRestore(t *testing.T) {
	var input []FlareModel.Bookmark
	input = append(input, FlareModel.Bookmark{Private: true})

	removed := restorePrivateProp(removePrivateProp(input))
	for i := 0; i < len(removed); i++ {
		if removed[i].Private != false {
			t.Fatal("Remove and restore private prop Failed")
		}
	}
}
