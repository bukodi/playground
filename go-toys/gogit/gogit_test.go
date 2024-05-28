package gogit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestInitRepo(t *testing.T) {

	directory := t.TempDir()

	r, err := git.PlainInit(directory, false)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// ... we need a file to commit so let's create a new file inside the
	// worktree of the project using the go standard library.
	filename := filepath.Join(directory, "example-git-file")
	err = os.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// Adds the new file to the staging area.
	_, err = w.Add("example-git-file")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
		Signer: nil,
	})

	if err != nil {
		t.Fatalf("%+v", err)
	}

	// Prints the current HEAD to verify that all worked well.
	obj, err := r.CommitObject(commit)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	t.Logf("Commit obj: %+v", obj)
	fileIter, err := obj.Files()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fileIter.ForEach(func(f *object.File) error {
		t.Logf("File name: %s", f.Name)
		return nil
	})
	tagRef, err := r.CreateTag("v0.0.1", commit, &git.CreateTagOptions{
		Tagger:  nil,
		Message: "",
		SignKey: nil,
	})
	t.Logf("Tag ref: %+v", tagRef)
}
