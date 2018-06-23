package main

import (
  "os"
  "path"
  "io/ioutil"
  "crypto/md5"
  "io"
  "math"
  "fmt"
  "encoding/hex"
  "github.com/urfave/cli"
)

const (
  fileChunk = 8192 // 8KB
  name      = "md5-ls"
  version   = "0.0.1"
)

/**
Calculate md5 of a file
 */
func calcMd5(filePath string) (string) {
  file, err := os.Open(filePath)

  if err != nil {
    panic(err.Error())
  }

  defer file.Close()

  // calculate the file size
  info, _ := file.Stat()

  fileSize := info.Size()

  blocks := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

  hash := md5.New()

  for i := uint64(0); i < blocks; i++ {
    blockSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
    buf := make([]byte, blockSize)

    file.Read(buf)
    io.WriteString(hash, string(buf)) // append into the hash
  }

  md5string := hex.EncodeToString(hash.Sum(nil))

  return md5string
}

func main() {
  app := cli.NewApp()

  app.Name = name
  app.Usage = name + " [path]"
  app.Version = version
  app.Description = "Display the file/dir md5 value"

  app.Action = func(c *cli.Context) (err error) {

    var (
      cwd      string
      pathStat os.FileInfo
      relative = c.Args().Get(0)
      distPath string // 最终目录/文件
    )

    if path.IsAbs(relative) {
      distPath = relative
    } else {
      if cwd, err = os.Getwd(); err != nil {
        return
      } else {
        distPath = path.Join(cwd, distPath)
      }
    }

    pathStat, err = os.Stat(distPath)

    if err != nil {
      panic(err)
      return
    }

    // if is a dir
    if pathStat.IsDir() {
      if fileInfos, err := ioutil.ReadDir(distPath); err != nil {
        return err
      } else {
        for i := range fileInfos {
          info := fileInfos[i]

          if info.IsDir() {
            // skip if it is a dir
            continue
          } else {
            fileName := info.Name()
            absFilePath := path.Join(distPath, fileName)
            md5string := calcMd5(absFilePath)
            fmt.Println(md5string + "  " + fileName)
          }
        }
      }
    }

    return nil
  }

  app.Run(os.Args)
}
