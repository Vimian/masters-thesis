package exe

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/minio/minio-go/v7"
)

type Exe struct{
	algorithmName string
}

func (e Exe) Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	var parts []string = strings.Split(objectInfo.Key, "/")
	var fileName string = parts[len(parts)-1]
	var compressedName string = fileName+".pmm"

	if err = os.WriteFile(fileName, data, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, "", err
	}
	cmd := exec.Command("wine", "app/" + e.algorithmName + ".exe", "e", "-f" + compressedName, fileName)
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

func (e Exe) Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	compressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	var parts []string = strings.Split(objectInfo.Key, "/")
	var compressedName string = parts[len(parts)-1]
	var fileName string = compressedName[:len(compressedName)-4]
	println(compressedName, fileName)

	if err = os.WriteFile(compressedName, compressedData, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, "", err
	}

	cmd := exec.Command("wine", "app/" + e.algorithmName + ".exe", "d", compressedName)
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