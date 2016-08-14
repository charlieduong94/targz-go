package targz

import (
    "archive/tar"
    "compress/gzip"
    "os"
    "io"
    "fmt"
    "path/filepath"
    "strings"
    "regexp"
)

func listToMap(ignoreList []string) (map[string]bool) {
    ignoreMap := make(map[string]bool)
    for _, value := range ignoreList {
        ignoreMap[value] = true
    }
    return ignoreMap
}

func writeToTar(tarWriter *tar.Writer, fileInfo os.FileInfo, sourcePath string, basePath string, ignoreMap map[string]bool) (error) {
    header, infoErr := tar.FileInfoHeader(fileInfo, fileInfo.Name())
    name := strings.Split(sourcePath, basePath)[1]
    regex, _ := regexp.Compile("[^\\/]*")
    if (ignoreMap[regex.FindString(name)]) {
        return nil
    }
    header.Name = name
    if (infoErr != nil) {
        return infoErr
    }
    if writeErr := tarWriter.WriteHeader(header); writeErr != nil {
        return writeErr
    }

    fileReader, openErr := os.Open(sourcePath)
    if (openErr != nil) {
        return openErr
    }
    defer fileReader.Close()
    _, copyErr := io.Copy(tarWriter, fileReader)
    return copyErr
}
/**
 *  Creates a tar file based off of the tar file
 *
 *
 */
func makeTar(source string, ignore map[string]bool) (error) {
    // create a tarfile
    target := fmt.Sprintf("%s.tar", source)
    tarfile, createErr := os.Create(target)
    if (createErr != nil) {
        return createErr
    }
    defer tarfile.Close()
    // open up tarfile for writing
    tarWriter := tar.NewWriter(tarfile)
    defer tarWriter.Close()

    fileInfo, statErr := os.Stat(source)
    if (statErr != nil) {
        return statErr
    }
    source += "/"
    if (!fileInfo.IsDir()) {
        return writeToTar(tarWriter, fileInfo, source, source, ignore)
    } else {
        // no err, walk the directory and tar all files
        return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
            if (info.IsDir()) {
                return nil
            }
            if (err != nil) {
                return err
            }
            return writeToTar(tarWriter, info, path, source, ignore)
        })
    }


}

func packageAndCompress(sourcePath string, targetPath string, ignoreMap map[string]bool) (error) {
    tarFilePath := fmt.Sprintf("%s.tar", sourcePath)
    tarErr := makeTar(sourcePath, ignoreMap)
    if (tarErr != nil) {
        return tarErr
    }
    reader, readErr := os.Open(tarFilePath);
    if (readErr != nil) {
        return readErr
    }

    writer, createErr := os.Create(targetPath)
    if (createErr != nil) {
        return createErr
    }

    defer writer.Close()
    archiver := gzip.NewWriter(writer)
    defer archiver.Close()

    _, copyErr := io.Copy(archiver, reader)
    if (copyErr != nil) {
        return copyErr
    }
    deleteErr := os.Remove(tarFilePath)
    return deleteErr
}

func Pack(sourcePath string, targetPath string) (error) {
    var ignoreMap map[string]bool
    return packageAndCompress(sourcePath, targetPath, ignoreMap)
}

func PackIgnore(sourcePath string, targetPath string, ignoreList []string) (error) {
    ignoreMap := listToMap(ignoreList)
    return packageAndCompress(sourcePath, targetPath, ignoreMap)
}

// unpacks the content of sourcepath (gzipped tar), into the targetPath
func Unpack(sourcePath string, targetPath string) (error) {
    regex, _ := regexp.Compile("([\\s\\S]*)\\/[\\s\\S]*$");

    // open file
    reader, openErr := os.Open(sourcePath)
    if (openErr != nil) {
        return openErr
    }
    // create a gzip reader
    gReader, gzipErr := gzip.NewReader(reader)
    if (gzipErr != nil) {
        return gzipErr
    }

    tarReader := tar.NewReader(gReader)

    // while
    for {
        header, err := tarReader.Next()
        if (err == io.EOF) {
            break;
        }
        matches := regex.FindStringSubmatch(header.Name)
        var path string
        if (len(matches) > 0) {
            path = matches[1]
        } else {
            path = ""
        }

        if (len(path) > 0) {
            mkdirErr := os.MkdirAll(filepath.Join(targetPath, path), 0777)
            if (mkdirErr != nil) {
                return mkdirErr
            }
        }
        file, createErr := os.Create(filepath.Join(targetPath, header.Name))
        defer file.Close()
        if (createErr != nil) {
            return createErr
        }
        _, copyErr := io.Copy(file, tarReader)
        if (copyErr != nil) {
            return copyErr
        }
    }
    return nil
}
