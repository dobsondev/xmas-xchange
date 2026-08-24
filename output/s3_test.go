package output

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeS3PutObjectAPI struct {
	gotBucket string
	gotKey    string
	gotBody   []byte
	err       error
}

func (f *fakeS3PutObjectAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.gotBucket = *params.Bucket
	f.gotKey = *params.Key
	body, err := io.ReadAll(params.Body)
	if err != nil {
		return nil, err
	}
	f.gotBody = body

	return &s3.PutObjectOutput{}, nil
}

func TestUploadToS3_Success(t *testing.T) {
	fake := &fakeS3PutObjectAPI{}
	data := []byte("hello = \"world\"\n")

	location, err := uploadToS3(context.Background(), fake, "my-bucket", "my-key.toml", data)
	if err != nil {
		t.Fatalf("uploadToS3 returned error: %v", err)
	}

	if location != "s3://my-bucket/my-key.toml" {
		t.Errorf("Expected location %q, got %q", "s3://my-bucket/my-key.toml", location)
	}
	if fake.gotBucket != "my-bucket" {
		t.Errorf("Expected bucket %q, got %q", "my-bucket", fake.gotBucket)
	}
	if fake.gotKey != "my-key.toml" {
		t.Errorf("Expected key %q, got %q", "my-key.toml", fake.gotKey)
	}
	if !bytes.Equal(fake.gotBody, data) {
		t.Errorf("Expected body %q, got %q", data, fake.gotBody)
	}
}

func TestUploadToS3_Error(t *testing.T) {
	fake := &fakeS3PutObjectAPI{err: errors.New("boom")}

	_, err := uploadToS3(context.Background(), fake, "my-bucket", "my-key.toml", []byte("data"))
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

type fakeS3GetObjectAPI struct {
	gotBucket string
	gotKey    string
	body      []byte
	err       error
}

func (f *fakeS3GetObjectAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.gotBucket = *params.Bucket
	f.gotKey = *params.Key

	return &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader(f.body))}, nil
}

func TestDownloadFromS3_Success(t *testing.T) {
	data := []byte("hello = \"world\"\n")
	fake := &fakeS3GetObjectAPI{body: data}

	got, err := downloadFromS3(context.Background(), fake, "my-bucket", "my-key.toml")
	if err != nil {
		t.Fatalf("downloadFromS3 returned error: %v", err)
	}

	if !bytes.Equal(got, data) {
		t.Errorf("Expected body %q, got %q", data, got)
	}
	if fake.gotBucket != "my-bucket" {
		t.Errorf("Expected bucket %q, got %q", "my-bucket", fake.gotBucket)
	}
	if fake.gotKey != "my-key.toml" {
		t.Errorf("Expected key %q, got %q", "my-key.toml", fake.gotKey)
	}
}

func TestDownloadFromS3_Error(t *testing.T) {
	fake := &fakeS3GetObjectAPI{err: errors.New("boom")}

	_, err := downloadFromS3(context.Background(), fake, "my-bucket", "my-key.toml")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestParseS3URI(t *testing.T) {
	testCases := []struct {
		name       string
		location   string
		wantBucket string
		wantKey    string
		wantOK     bool
	}{
		{"valid uri", "s3://my-bucket/my-key.toml", "my-bucket", "my-key.toml", true},
		{"nested key", "s3://my-bucket/path/to/key.toml", "my-bucket", "path/to/key.toml", true},
		{"local path", "exchange-toml/my-key.toml", "", "", false},
		{"missing key", "s3://my-bucket/", "", "", false},
		{"missing key no slash", "s3://my-bucket", "", "", false},
		{"missing bucket", "s3:///my-key.toml", "", "", false},
		{"empty string", "", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bucket, key, ok := ParseS3URI(tc.location)
			if ok != tc.wantOK || bucket != tc.wantBucket || key != tc.wantKey {
				t.Errorf("ParseS3URI(%q) = (%q, %q, %v), want (%q, %q, %v)", tc.location, bucket, key, ok, tc.wantBucket, tc.wantKey, tc.wantOK)
			}
		})
	}
}
