package tcache

import (
	"github.com/golang/groupcache"
	"testing"
)

func TestGroupcache(t *testing.T) {
	me := "http://10.0.0.1"
	peers := groupcache.NewHTTPPool(me)

	// Whenever peers change:
	peers.Set("http://10.0.0.1", "http://10.0.0.2", "http://10.0.0.3")

	var thumbNails = groupcache.NewGroup("user", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			fileName := key
			dest.SetBytes(generateThumbnail(t, fileName))
			return nil
		}))

	{
		var data []byte
		err := thumbNails.Get(t.Context(), "first", groupcache.AllocatingByteSliceSink(&data))
		if err != nil {
			t.Fatalf("failed to get thumbnail for first: %v", err)
		}
		t.Logf("thumbnail for first: %s", data)
	}
	{
		var data []byte
		err := thumbNails.Get(t.Context(), "first", groupcache.AllocatingByteSliceSink(&data))
		if err != nil {
			t.Fatalf("failed to get thumbnail for first: %v", err)
		}
		t.Logf("thumbnail for first: %s", data)
	}
	{
		var data []byte
		err := thumbNails.Get(t.Context(), "second", groupcache.AllocatingByteSliceSink(&data))
		if err != nil {
			t.Fatalf("failed to get thumbnail for second: %v", err)
		}
		t.Logf("thumbnail for second: %s", data)
	}
}

func generateThumbnail(t *testing.T, name string) []byte {
	t.Logf("generating thumbnail for %s", name)
	return []byte("thumbnail for " + name + "")
}
