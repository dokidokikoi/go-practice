package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filePath, url string) error {
	out, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	out.Close()

	fmt.Print("\n")

	err = os.Rename(filePath+".tmp", filePath)
	if err != nil {
		return err
	}

	return nil
}
