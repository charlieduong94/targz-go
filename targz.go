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

func writeToTar(tarWriter *tar.Writer, fileInfo os.FileInfo, sourcePath string, basePath string) (error) {
    header, infoErr := tar.FileInfoHeader(fileInfo, fileInfo.Name())
    fmt.Println(fileInfo.Name())
    header.Name = strings.Split(sourcePath, basePath)[1]
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
func makeTar(source string) (error) {
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
        return writeToTar(tarWriter, fileInfo, source, source)
    } else {
        // no err, walk the directory and tar all files
        return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
            if (info.IsDir()) {
                return nil
            }
            if (err != nil) {
                return err
            }
            return writeToTar(tarWriter, info, path, source)
        })
    }


}


func Pack(sourcePath string, targetPath string) (error) {
    tarFilePath := fmt.Sprintf("%s.tar", sourcePath)
    tarErr := makeTar(sourcePath)
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
        fmt.Println(header.Name)
        matches := regex.FindStringSubmatch(header.Name)
        fmt.Println(matches);
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
