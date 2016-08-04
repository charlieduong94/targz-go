package targz

import (
    "testing"
    "os"
    "fmt"
)

func TestPack(t *testing.T) {
    dir, _ := os.Getwd()
    target := fmt.Sprintf("%s/%s", dir, "test.tar.gz")
    source := dir
    err := Pack(source, target)
    if (err != nil) {
        fmt.Println(err)
        t.Fail()
    }
}


func TestUnpack(t *testing.T) {
    dir, _ := os.Getwd()
    source := fmt.Sprintf("%s/%s", dir, "test.tar.gz")
    target := fmt.Sprintf("%s/test", dir);
    mkdirErr := os.MkdirAll(target, 0777)
    if (mkdirErr != nil) {
        fmt.Println(mkdirErr)
        t.Fail()
    }
    err := Unpack(source, target)
    if (err != nil) {
        fmt.Println(err)
        t.Fail()
    }
}
