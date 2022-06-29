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
	Name string
}

type ImageData struct {
	Title string
	Paths []string
}

func getMp4PathAndNames(path string) (result []VideoData) {
	paths := getFileList(path + "/video")
	sort.Strings(paths)

	for _, path := range paths {
		result = append(result, VideoData { path, getFileNameFromPath(path) })
	}
	return
}

func getJpgNames(path string) []string {
	return getFileList(path)
}

func Init() {
	encoded_config, _ := os.ReadFile(DATA_PATH + "config.json")
	json.Unmarshal(encoded_config, &configs)
	// fmt.Println(getMp4Names(configs[3].Path))
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
		templates.ExecuteTemplate(w, "media.html", 0)
	} else {
		mediaType, path, err := lookup(mediaName)
		if err != nil {
			http.NotFound(w, r)
		} else {
			fmt.Println(mediaName)
			switch mediaType {
				case "jpg":
					imageData := ImageData {path, getJpgNames(path)}
					templates.ExecuteTemplate(w, "images.html", imageData)
				case "mp4":
					videos := getMp4PathAndNames(path)
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
