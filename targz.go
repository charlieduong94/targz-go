package tarball

import (
    "archive/tar"
    "compress/gzip"
    "os"
    "io"
    "fmt"
    "path/filepath"
)

func writeToTar(tarWriter *tar.Writer, fileInfo os.FileInfo, targetPath string) (error) {
    header, infoErr := tar.FileInfoHeader(fileInfo, fileInfo.Name())
    if (infoErr != nil) {
        return infoErr
    }

    if writeErr := tarWriter.WriteHeader(header); writeErr != nil {
        return writeErr
    }

    fileReader, openErr := os.Open(targetPath)
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

    if (!fileInfo.IsDir()) {
        return writeToTar(tarWriter, fileInfo, source)
    }

    // no err, walk the directory and tar all files
    return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
        if (info.IsDir()) {
            return nil
        }

        if (err != nil) {
            return err
        }
        return writeToTar(tarWriter, info, path)
    })
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
    fmt.Println(deleteErr)
    return deleteErr
}

// unpacks the content of sourcepath (gzipped tar), into the targetPath
func Unpack(sourcePath string, targetPath string) (error) {
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
        info := header.FileInfo()
        if (info.IsDir()) {
            mkdirErr := os.MkdirAll(header.Name, info.Mode())
            if (mkdirErr != nil) {
                return mkdirErr
            }
        } else {
            file, createErr := os.Create(filepath.Join(targetPath, header.Name))
            if (createErr != nil) {
                return createErr
            }
            _, copyErr := io.Copy(file, tarReader)
            if (copyErr != nil) {
                return copyErr
            }
        }
    }
    return nil
}
