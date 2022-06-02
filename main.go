package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type CLI struct {
	CacheFile string
}

var cli *CLI

type File struct {
	Link string
	Path string
}

const (
	TEMP_DIR                = "/tmp/dxvk-cache-grabber"
	REDDIT_DXVK_CACHE_POST  = `https://www.reddit.com/r/linux_gaming/comments/t5xrho/.json`
	LINK_PATTERN            = `(https:\/\/cdn\.discordapp\.com\/attachments)[-A-Z0-9+&@#\/%?=~_|$!:,.;]*[A-Z0-9+&@#\/%=~_|$](r5apex\.dxvk-cache)`
	MERGE_TOOL              = "./dxvk-cache-tool"
)

func MergeDXVKCache(file File) {
	cmd := exec.Command(MERGE_TOOL, cli.CacheFile, file.Path, "-o", cli.CacheFile)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error while merging cache file:", err)
	}

	fmt.Println("DXVK cache merged", string(out))
}

func CreateTmpDirectory() {
	// create directory if not exists
	if _, err := os.Stat(TEMP_DIR); os.IsNotExist(err) {
		err := os.Mkdir(TEMP_DIR, 0777)

		if err != nil {
			fmt.Println("Error creating tmp directory", err)
			os.Exit(1)
		}
	} else {
		os.RemoveAll(TEMP_DIR)
		CreateTmpDirectory()
	}
}

func Request(link string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "dxvk-grabber")

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetPostData() ([]byte, error) {
	response, err := Request(REDDIT_DXVK_CACHE_POST)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func ExtractLinks(data string) []string {
	regex := regexp.MustCompile(LINK_PATTERN)
	matches := regex.FindAllString(data, -1)

	return matches
}

func DownloadFile(file File) error {
	// Get the data
	resp, err := Request(file.Link)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(file.Path)

	if err != nil {
		return err
	}

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func DownloadFiles(links []string) {
	for _, link := range links {
		fmt.Println("Downloading link:", link)
		name := fmt.Sprintf("%s/%s", TEMP_DIR, fmt.Sprintf("%d.%s", time.Now().Unix(), strings.Split(link, "/")[len(strings.Split(link, "/"))-1]))

		file := File{
			Link: link,
			Path: name,
		}

		err := DownloadFile(file)

		if err != nil {
			fmt.Println("error occurred while downloading link:", err)
		}

		MergeDXVKCache(file)
	}
}

func HandleFlags() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf(" %s [options]\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}
  var cacheFile string

  flag.StringVar(&cacheFile, "c", "", "Steam DXVK cache file")
  flag.Parse()

  if cacheFile == "" {
    fmt.Print("You must enter cache file location check help: '-h'")
    os.Exit(1)
  }

  if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
    fmt.Println("Cache file does not exist")
    os.Exit(1)
  }

  cli = &CLI{
    CacheFile: cacheFile,
  }
}

func main() {
  HandleFlags()
	CreateTmpDirectory()
	data, err := GetPostData()

	if err != nil {
		fmt.Println("Error getting DXVK cache", err)
		os.Exit(1)
	}

	jsonData := string(data)

	links := ExtractLinks(jsonData)

	DownloadFiles(links)
}
