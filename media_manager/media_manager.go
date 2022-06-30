package media_manager

import (
	"os"
	"fmt"
	"sort"
	"errors"
	"strings"
	"net/http"
	"encoding/json"
	"html/template"
)

const DATA_PATH = "./static/media_manager_data/"

type dataConfig struct {
	Name string
	Type string
	Path string
}

var configs []dataConfig

var templates = template.Must(template.ParseFiles(
	"tmpl/media.html",
	"tmpl/videos.html",
	"tmpl/images.html",
	"tmpl/templates.html",
))

type VideoData struct {
	Path string
	SubtitlePath string
	Name string
}

type ImageData struct {
	Title string
	Paths []string
}

func getVideoPathAndNames(folderPath string) (result []VideoData) {
	fmt.Println("path: ", folderPath)
	paths := getFileList(folderPath + "/video")
	sort.Strings(paths)

	for _, path := range paths {
		name := getFileNameFromPath(path)
		// todo: don't hardcode
		subTitlePath := "/static/res/" + folderPath + "/sub/" + name + ".vtt"
		result = append(result, VideoData { path, subTitlePath, name })
	}
	return
}

func getImageNames(path string) []string {
	return getFileList(path)
}

func Init() {
	encoded_config, _ := os.ReadFile(DATA_PATH + "config.json")
	json.Unmarshal(encoded_config, &configs)
}

/**
 * images have their own page for each folder
 * videos have only one for all
*/

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	mediaName := r.URL.Path[len("/media/"):]

	if r.Method != "GET" {
		panic("media only deals with GET")
	}

	if len(mediaName) == 0 {
		templates.ExecuteTemplate(w, "media.html", configs)
	} else {
		mediaType, path, err := lookup(mediaName)
		if err != nil {
			http.NotFound(w, r)
		} else {
			switch mediaType {
				case "image":
					imageData := ImageData {path, getImageNames(path)}
					templates.ExecuteTemplate(w, "images.html", imageData)
				case "video":
					videos := getVideoPathAndNames(path)
					templates.ExecuteTemplate(w, "videos.html", videos)
				default:
					panic("invalid media type")
			}
		}
	}
}

func lookup(mediaName string) (mediaType, path string, err error) {
	for _, config := range configs {
		if (mediaName == config.Name) {
			return config.Type, config.Path, nil
		}
	}
	return "", "", errors.New("not found")
}

func getFileNameFromPath(path string) (name string) {
	elements := strings.Split(path, "/")
	filename := elements[len(elements) - 1]
	filename = strings.Split(filename, ".")[0]
	return filename
}

func getFileList(folder string) []string {
	path := "./static/res/" + folder
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []string{}

	for _, v := range files {
		result = append(result, path[len("."):] + "/" + v.Name())
	}
	return result
}
