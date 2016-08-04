# targz-go

A simple module for packing and unpacking tarballs.

## Usage

#### Packing
```
packErr := targz.Pack(sourcePath, targetPath)
if (packErr != nil) {
    fmt.Println("Unable pack file/directory")
}
```

#### Unpacking
```
unpackErr := targz.Unpack(sourcePath, targetPath)
if (unpackErr != nil) {
    fmt.Println("Unable to unpack tar.gz content into the target path")
}

```
