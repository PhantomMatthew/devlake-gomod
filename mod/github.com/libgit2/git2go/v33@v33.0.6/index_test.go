package git

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestCreateRepoAndStage(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("README")
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexReadTree(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	_, _ = seedTestRepo(t, repo)

	ref, err := repo.Head()
	checkFatal(t, err)

	obj, err := ref.Peel(ObjectTree)
	checkFatal(t, err)

	tree, err := obj.AsTree()
	checkFatal(t, err)

	idx, err := NewIndex()
	checkFatal(t, err)

	err = idx.ReadTree(tree)
	checkFatal(t, err)

	id, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if tree.Id().Cmp(id) != 0 {
		t.Fatalf("Read and written trees are not the same")
	}
}

func TestIndexWriteTreeTo(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	idx, err := NewIndex()
	checkFatal(t, err)

	odb, err := repo.Odb()
	checkFatal(t, err)

	content, err := ioutil.ReadFile(path.Join(repo.Workdir(), "README"))
	checkFatal(t, err)

	id, err := odb.Write(content, ObjectBlob)
	checkFatal(t, err)

	err = idx.Add(&IndexEntry{
		Mode: FilemodeBlob,
		Uid:  0,
		Gid:  0,
		Size: uint32(len(content)),
		Id:   id,
		Path: "README",
	})
	checkFatal(t, err)

	treeId, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexAddAndWriteTreeTo(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	odb, err := repo.Odb()
	checkFatal(t, err)

	blobID, err := odb.Write([]byte("foo\n"), ObjectBlob)
	checkFatal(t, err)

	idx, err := NewIndex()
	checkFatal(t, err)

	if idx.Path() != "" {
		t.Fatal("in-memory repo has a path")
	}

	entry := IndexEntry{
		Path: "README",
		Id:   blobID,
		Mode: FilemodeBlob,
	}

	err = idx.Add(&entry)
	checkFatal(t, err)

	treeId, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexRemoveDirectory(t *testing.T) {
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	odb, err := repo.Odb()
	checkFatal(t, err)

	blobID, err := odb.Write([]byte("fou\n"), ObjectBlob)
	checkFatal(t, err)

	idx, err := NewIndex()
	checkFatal(t, err)

	entryCount := idx.EntryCount()
	if entryCount != 0 {
		t.Fatal("Index should count 0 entry")
	}

	entry := IndexEntry{
		Path: "path/to/LISEZ_MOI",
		Id:   blobID,
		Mode: FilemodeBlob,
	}

	err = idx.Add(&entry)
	checkFatal(t, err)

	entryCount = idx.EntryCount()
	if entryCount != 1 {
		t.Fatal("Index should count 1 entry")
	}

	err = idx.RemoveDirectory("path", 0)

	entryCount = idx.EntryCount()
	if entryCount != 0 {
		t.Fatal("Index should count 0 entry")
	}
}

func TestIndexAddFromBuffer(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	idx, err := repo.Index()
	checkFatal(t, err)

	entry := IndexEntry{
		Path: "README",
		Mode: FilemodeBlob,
	}

	err = idx.AddFromBuffer(&entry, []byte("foo\n"))
	checkFatal(t, err)

	treeId, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexAddAllNoCallback(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	err := ioutil.WriteFile(repo.Workdir()+"/README", []byte("foo\n"), 0644)
	checkFatal(t, err)

	idx, err := repo.Index()
	checkFatal(t, err)

	err = idx.AddAll([]string{}, IndexAddDefault, nil)
	checkFatal(t, err)

	treeId, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexAddAllCallback(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	err := ioutil.WriteFile(repo.Workdir()+"/README", []byte("foo\n"), 0644)
	checkFatal(t, err)

	idx, err := repo.Index()
	checkFatal(t, err)

	cbPath := ""
	err = idx.AddAll([]string{}, IndexAddDefault, func(p, mP string) error {
		cbPath = p
		return nil
	})
	checkFatal(t, err)
	if cbPath != "README" {
		t.Fatalf("%v", cbPath)
	}

	treeId, err := idx.WriteTreeTo(repo)
	checkFatal(t, err)

	if treeId.String() != "b7119b11e8ef7a1a5a34d3ac87f5b075228ac81e" {
		t.Fatalf("%v", treeId.String())
	}
}

func TestIndexOpen(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	path := repo.Workdir() + "/heyindex"

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Fatal("new index file already exists")
	}

	idx, err := OpenIndex(path)
	checkFatal(t, err)

	if path != idx.Path() {
		t.Fatalf("mismatched index paths, expected %v, got %v", path, idx.Path())
	}

	err = idx.Write()
	checkFatal(t, err)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		t.Fatal("new index file did not get written")
	}
}

func checkFatal(t *testing.T, err error) {
	if err == nil {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Unable to get caller")
	}
	t.Fatalf("Fail at %v:%v; %v", file, line, err)
}
