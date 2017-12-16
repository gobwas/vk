package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gobwas/vk"
	"github.com/gobwas/vk/internal/httputil"
)

func GetLargestSize(sizes []vk.PhotoSize) (max vk.PhotoSize) {
	// Range from the end of sizes cause there is a nice chance that 'w' type
	// is the last one.
	for i := len(sizes) - 1; i >= 0; i-- {
		size := sizes[i]
		if size.Type == vk.SizeW {
			// Largest possible size.
			return size
		}
		if max.Type.Less(size.Type) {
			max = size
		}
	}
	return max
}

func GetDefaultDest(suffix string) (path string) {
	user, err := user.Current()
	if err == nil {
		path = user.HomeDir
	}
	if path == "" {
		path = "/tmp"
	}
	return path + "/vk/" + suffix
}

func Photo(ctx context.Context, destDir string, photo vk.Photo, size vk.PhotoSize) error {
	photoID := strconv.Itoa(photo.ID)
	ext := path.Ext(size.Src)

	filepath := filepath.Clean(fmt.Sprintf("%s/%s%s", destDir, photoID, ext))

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", size.Src, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := httputil.CheckResponseStatus(resp); err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)

	return err
}
