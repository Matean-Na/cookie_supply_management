package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadXML(url, filename string) error {
	//create media if not exist
	if _, err := os.Stat("media"); os.IsNotExist(err) {
		if err = os.Mkdir("media", 0755); err != nil {
			return err
		}
	}

	//open file for recording
	file, err := os.Create(filepath.Join("media", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	//get file xml from url
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//recording body response in file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
