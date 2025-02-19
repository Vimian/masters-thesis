package sevenzip

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/minio/minio-go/v7"
)

type Sevenzip struct{
	algorithmName string
}

func (s Sevenzip) Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	var parts []string = strings.Split(objectInfo.Key, "/")
	var fileName string = parts[len(parts)-1]
	var compressedName string = fileName+".7z"

	if err = os.WriteFile(fileName, data, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, "", err
	}

	cmd := exec.Command("7z", "a", "-t7z", compressedName, fileName, "-mx=9", "-m0="+s.algorithmName)
	if err = cmd.Run(); err != nil {
		return nil, "", err
	}

	compressedData, err := os.ReadFile(compressedName)
	if err != nil {
		return nil, "", err
	}

	cmd = exec.Command("rm", fileName, compressedName)
	if err = cmd.Run(); err != nil {
		return nil, "", err
	}

	return bytes.NewReader(compressedData), compressedName, nil
}

func (s Sevenzip) Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	compressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	var parts []string = strings.Split(objectInfo.Key, "/")
	var compressedName string = parts[len(parts)-1]
	var fileName string = compressedName[:len(compressedName)-3]

	if err = os.WriteFile(compressedName, compressedData, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, "", err
	}

	cmd := exec.Command("7z", "e", compressedName)
	if err = cmd.Run(); err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, "", err
	}
	
	cmd = exec.Command("rm", fileName, compressedName)
	if err = cmd.Run(); err != nil {
		return nil, "", err
	}

	return bytes.NewReader(data), fileName, nil
}