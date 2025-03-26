package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

type JSONEncoding struct{}

func (JSONEncoding) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONEncoding) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type JSONGzipEncoding struct{}

func (JSONGzipEncoding) Marshal(v any) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	err = json.NewEncoder(writer).Encode(v)
	if err != nil {
		writer.Close() // nolint: errcheck
		return nil, err
	}
	writer.Close() // nolint: errcheck
	return buf.Bytes(), nil
}

func (JSONGzipEncoding) Unmarshal(data []byte, v any) error {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer reader.Close() // nolint: errcheck
	return json.NewDecoder(reader).Decode(v)
}
