package db

import "testing"

func TestTable_GetFilePath(t *testing.T) {
	path := database.GetFilePath("超话")
	t.Log(path)
}
