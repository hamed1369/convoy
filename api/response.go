package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"strings"
)

type ErrorResponse struct {
	Error string
}

type VolumesResponse struct {
	Volumes map[string]VolumeResponse
}

type VolumeResponse struct {
	UUID       string
	Name       string
	Base       string
	Size       int64
	MountPoint string
	Snapshots  map[string]SnapshotResponse
}

type SnapshotResponse struct {
	UUID       string
	VolumeUUID string
}

type BlockStoreResponse struct {
	UUID      string
	Kind      string
	BlockSize int64
}

type ImageResponse struct {
	UUID        string
	Name        string
	Size        int64
	Checksum    string
	RawChecksum string
}

func ResponseError(format string, a ...interface{}) {
	response := ErrorResponse{Error: fmt.Sprintf(format, a...)}
	j, err := json.MarshalIndent(&response, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Failed to generate response for error:", err))
	}
	fmt.Println(string(j[:]))
}

func ResponseLogAndError(v interface{}) {
	if e, ok := v.(*logrus.Entry); ok {
		e.Error(e.Message)
		oldFormatter := e.Logger.Formatter
		logrus.SetFormatter(&logrus.JSONFormatter{})
		s, err := e.String()
		logrus.SetFormatter(oldFormatter)
		if err != nil {
			ResponseError(err.Error())
			return
		}
		// Cosmetic since " would be escaped
		ResponseError(strings.Replace(s, "\"", "'", -1))
	} else if e, ok := v.(error); ok {
		logrus.Errorf(fmt.Sprint(e))
		ResponseError(fmt.Sprint(e))
	}
}

func ResponseOutput(v interface{}) ([]byte, error) {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return j, nil
}
