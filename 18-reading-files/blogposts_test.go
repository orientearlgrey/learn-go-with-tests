package blogposts_test

import (
	"blogposts"
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func assertPost(t *testing.T, got, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %+v but got %+v", want, got)
	}
}

func TestNewBlogPosts(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte("Title: Post 1")},
			"hello-world2.md": {Data: []byte("Title: Post 2")},
		}

		posts, err := blogposts.NewPostsFromFS(fs)
		if err != nil {
			t.Fatal(err)
		}

		assertPost(t, posts[0], blogposts.Post{Title: "Post 1"})
	})

	t.Run("load failure", func(t *testing.T) {
		_, err := blogposts.NewPostsFromFS(StubFailingFS{})
		if err == nil {
			t.Fatal("expected an error but got nil")
		}
	})
}