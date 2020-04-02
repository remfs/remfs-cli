package main

import (
        "fmt"
        "os"
        "strings"
        "net/http"
        "io"
        "io/ioutil"
        "encoding/json"
        "path"
)

type RemFS struct {
        Type string `json:"type"`
        Children map[string]*RemFS`json:"children"`
}


func main() {
        if len(os.Args) < 3 {
                os.Exit(1)
        }

        baseAddress := os.Args[1]
        outDir := os.Args[2]

        downloadDirectory(baseAddress, outDir)
}

func downloadDirectory(baseAddress string, parentPath string) {

        var remfsAddress string
        if strings.HasSuffix(baseAddress, "/remfs.json") {
                // noop
        } else if strings.HasSuffix(baseAddress, "/") {
                remfsAddress = baseAddress + "remfs.json"
        } else {
                baseAddress = baseAddress + "/"
                remfsAddress = baseAddress + "remfs.json"
        }

        req, err := http.Get(remfsAddress)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        remfsJson, err := ioutil.ReadAll(req.Body)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        remfs := &RemFS{}

        err = json.Unmarshal(remfsJson, remfs)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        err = os.MkdirAll(parentPath, 0755)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        for name, child := range remfs.Children {
                if child.Type == "file" {
                        downloadFile(baseAddress, parentPath, name)
                } else {
                        nextParentPath := path.Join(parentPath, name)
                        nextBaseAddress := baseAddress + "/" + name
                        downloadDirectory(nextBaseAddress, nextParentPath)
                }
        }

}

func downloadFile(parentUrl string, parentPath string, filename string) {

        url := parentUrl + filename

        req, err := http.Get(url)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        outPath := path.Join(parentPath, filename)
        f, err := os.Create(outPath)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        io.Copy(f, req.Body)
}
