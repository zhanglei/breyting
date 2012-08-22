package conf

import "container/list"
import dbg "fmt"

import "github.com/mewmew/breyting/download"

/// ### [ todo ] ###
///   - There are still some bugs the addition and removal of active page
///     watchers.
///   - Try to figure out a better way of implementing this.
/// ### [/ todo ] ###

var pageList *list.List

func ListAdd(page *download.Page) {
	if ListPresent(page) {
		return
	}
	pageList.PushBack(page)
}

func ListPresent(pageA *download.Page) bool {
	for e := pageList.Front(); e != nil; e = e.Next() {
		pageB := e.Value.(*download.Page)
		if pageA.Equal(pageB) {
			return true
		}
	}
	return false
}

func ListCleanup() {
	for e := pageList.Front(); e != nil; e = e.Next() {
		page := e.Value.(*download.Page)
		if !isPageActive(page) {
			dbg.Println("killing page:", page.RawUrl, page.RawSel)
			page.Kill()
			defer pageList.Remove(e)
		}
	}
}

func init() {
	pageList = list.New()
}
