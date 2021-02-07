package imageresizeservice

import (
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3FileWriter struct{
	batchSize int
	readers map[string]io.Reader
	readerMutex sync.Mutex
	errChan chan error
}

func newS3FileUploader(batchSize int) (*s3FileWriter)  {
	s := &s3FileWriter{
		batchSize: batchSize,
		readers: make(map[string]io.Reader, 0),
		errChan: make(chan error),
	}

	go func(){
		for {
			s.readerMutex.Lock()
			if len(s.readers) >= batchSize {
				s.errChan <- s.batchUploadFiles("cm-personal-site-bucket", s.readers)
				s.readers = make(map[string]io.Reader, 0)
				close(s.errChan)
			}
			s.readerMutex.Unlock()
		}
	}()

	return s
}

func (s *s3FileWriter) getUploadErr() error {
	for err := range s.errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *s3FileWriter)  writeFile(path string, r io.Reader) error {
	s.readerMutex.Lock()
	s.readers[path] = r
	s.readerMutex.Unlock()
	return nil
}

func (s *s3FileWriter) batchUploadFiles(bucket string, files map[string]io.Reader) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	svc := s3manager.NewUploader(sess)
	
	objects := []s3manager.BatchUploadObject{}

	for path := range files {
		objects = append(objects, s3manager.BatchUploadObject{
			Object: &s3manager.UploadInput{
				Key: aws.String(path),
				Body: files[path],
				Bucket: aws.String(bucket),
			},
		},
	)}

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	err = svc.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	return nil
}