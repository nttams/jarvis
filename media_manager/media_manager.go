package media_manager

import (
	"encoding/json"
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const DATA_PATH = "/data/media_manager_data/"

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
	Title string
	Paths []VideoPath
}

type VideoPath struct {
	Path         string
	SubtitlePath string
	Name         string
}

type ImageData struct {
	Title string
	Paths []string
}

func getVideoPaths(folderPath string) (result []VideoPath) {
	paths := getFileList(folderPath + "/video")
	sort.Strings(paths)

	for _, path := range paths {
		name := getFileNameFromPath(path)
		// todo: don't hardcode
		subTitlePath := "/static/res/" + folderPath + "/sub/" + name + ".vtt"
		result = append(result, VideoPath{path, subTitlePath, name})
	}
	return
}

func getImagePaths(path string) []string {
	return getFileList(path)
}

func Init() {
	encodedConfig, _ := os.ReadFile(DATA_PATH + "config.json")
	json.Unmarshal(encodedConfig, &configs)
	rand.Seed(time.Now().UnixNano())
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
				paths := getImagePaths(path)

				rand.Shuffle(len(paths), func(i, j int) { paths[i], paths[j] = paths[j], paths[i] })

				imageData := ImageData{path, paths}
				templates.ExecuteTemplate(w, "images.html", imageData)
			case "video":
				videoData := VideoData{path, getVideoPaths(path)}
				templates.ExecuteTemplate(w, "videos.html", videoData)
			default:
				panic("invalid media type")
			}
		}
	}
}

func lookup(mediaName string) (mediaType, path string, err error) {
	for _, config := range configs {
		if mediaName == config.Name {
			return config.Type, config.Path, nil
		}
	}
	return "", "", errors.New("not found")
}

func getFileNameFromPath(path string) (name string) {
	elements := strings.Split(path, "/")
	filename := elements[len(elements)-1]
	filename = strings.Split(filename, ".")[0]
	return filename
}

// todo: improve this
func getFileList(folder string) []string {
	path := "./static/res/" + folder
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []string{}

	for _, v := range files {
		result = append(result, path[len("."):]+"/"+v.Name())
	}
	return result
}
