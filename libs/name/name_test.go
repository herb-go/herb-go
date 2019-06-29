package name

import "testing"

func TestName(t *testing.T) {
	_, err := New(false, "你好")
	if err == nil {
		t.Fatal(err)
	}

	n, err := New(false, "test")
	if err != nil {
		t.Fatal(err)
	}

	if n.Camel != "test" {
		t.Fatal(n.Camel)
	}
	if n.Lower != "test" {
		t.Fatal(n.Lower)
	}
	if n.LowerPath("123.txt") != "test/123.txt" {
		t.Fatal(n.LowerPath("123.txt"))
	}
	if n.LowerWithParentPath != "test" {
		t.Fatal(n.LowerWithParentPath)
	}
	if n.LowerWithParentDotSeparated != "test" {
		t.Fatal(n.LowerWithParentDotSeparated)
	}
	if n.Parents != "" {
		t.Fatal(n.Parents)
	}
	if len(n.ParentsList) != 0 {
		t.Fatal(n.ParentsList)
	}
	if n.Pascal != "Test" {
		t.Fatal(n.Pascal)
	}
	if n.PascalWithParents != "Test" {
		t.Fatal(n.PascalWithParents)
	}
	if n.Raw != "test" {
		t.Fatal(n.Raw)
	}
	if n.Title != "Test" {
		t.Fatal(n.Title)
	}

	n, err = New(false, "test Test2-test3_id")
	if err != nil {
		t.Fatal(err)
	}

	if n.Camel != "testTest2Test3ID" {
		t.Fatal(n.Camel)
	}
	if n.Lower != "testtest2test3id" {
		t.Fatal(n.Lower)
	}
	if n.LowerPath("123.txt") != "testtest2test3id/123.txt" {
		t.Fatal(n.LowerPath("123.txt"))
	}
	if n.LowerWithParentPath != "testtest2test3id" {
		t.Fatal(n.LowerWithParentPath)
	}
	if n.LowerWithParentDotSeparated != "testtest2test3id" {
		t.Fatal(n.LowerWithParentDotSeparated)
	}
	if n.Parents != "" {
		t.Fatal(n.Parents)
	}
	if len(n.ParentsList) != 0 {
		t.Fatal(n.ParentsList)
	}
	if n.Pascal != "TestTest2Test3ID" {
		t.Fatal(n.Pascal)
	}
	if n.PascalWithParents != "TestTest2Test3ID" {
		t.Fatal(n.PascalWithParents)
	}
	if n.Raw != "test Test2-test3_id" {
		t.Fatal(n.Raw)
	}
	if n.Title != "Test Test2-test3_id" {
		t.Fatal(n.Title)
	}

	n, err = New(true, "test Test2-test3_id")
	if err != nil {
		t.Fatal(err)
	}

	if n.Camel != "testTest2Test3ID" {
		t.Fatal(n.Camel)
	}
	if n.Lower != "testtest2test3id" {
		t.Fatal(n.Lower)
	}
	if n.LowerPath("123.txt") != "testtest2test3id/123.txt" {
		t.Fatal(n.LowerPath("123.txt"))
	}
	if n.LowerWithParentPath != "testtest2test3id" {
		t.Fatal(n.LowerWithParentPath)
	}
	if n.LowerWithParentDotSeparated != "testtest2test3id" {
		t.Fatal(n.LowerWithParentDotSeparated)
	}
	if n.Parents != "" {
		t.Fatal(n.Parents)
	}
	if len(n.ParentsList) != 0 {
		t.Fatal(n.ParentsList)
	}
	if n.Pascal != "TestTest2Test3ID" {
		t.Fatal(n.Pascal)
	}
	if n.PascalWithParents != "TestTest2Test3ID" {
		t.Fatal(n.PascalWithParents)
	}
	if n.Raw != "test Test2-test3_id" {
		t.Fatal(n.Raw)
	}
	if n.Title != "Test Test2-test3_id" {
		t.Fatal(n.Title)
	}
	n, err = New(true, "my folder-1_Id/folder2/test Test2-test3_id")
	if err != nil {
		t.Fatal(err)
	}
	if n.Camel != "testTest2Test3ID" {
		t.Fatal(n.Camel)
	}
	if n.Lower != "testtest2test3id" {
		t.Fatal(n.Lower)
	}
	if n.LowerPath("123.txt") != "my folder-1_Id/folder2/testtest2test3id/123.txt" {
		t.Fatal(n.LowerPath("123.txt"))
	}
	if n.LowerWithParentPath != "my folder-1_Id/folder2/testtest2test3id" {
		t.Fatal(n.LowerWithParentPath)
	}
	if n.LowerWithParentDotSeparated != "myfolder1id.folder2.testtest2test3id" {
		t.Fatal(n.LowerWithParentDotSeparated)
	}
	if n.Parents != "my folder-1_Id/folder2" {
		t.Fatal(n.Parents)
	}
	if len(n.ParentsList) != 2 {
		t.Fatal(n.ParentsList)
	}
	if n.Pascal != "TestTest2Test3ID" {
		t.Fatal(n.Pascal)
	}
	if n.PascalWithParents != "MyFolder1IDFolder2TestTest2Test3ID" {
		t.Fatal(n.PascalWithParents)
	}
	if n.Raw != "test Test2-test3_id" {
		t.Fatal(n.Raw)
	}
	if n.Title != "Test Test2-test3_id" {
		t.Fatal(n.Title)
	}

	n, err = New(false, "my folder-1_Id/folder2/test Test2-test3_id")
	if err == nil {
		t.Fatal(err)
	}
}
