package main

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3FileWriter struct{
	batchSize int
	readers map[string]io.Reader
}

func newS3FileWriter(batchSize int) (*s3FileWriter, chan error)  {
	fileWriter := &s3FileWriter{
		batchSize: batchSize,
		readers: make(map[string]io.Reader, 0),
	}

	errChan := make(chan error)

	go func(){
		for {
			if len(fileWriter.readers) >= batchSize {
				errChan <- fileWriter.batchUploadFiles("cm-personal-site-bucket", fileWriter.readers)
				fileWriter.readers = make(map[string]io.Reader, 0)
			}
		}
	}()

	return fileWriter, errChan

}

func (s *s3FileWriter)  writeFile(path string, r io.Reader) error {
	s.readers[path] = r
	return nil
}

func (s *s3FileWriter) batchUploadFiles(bucket string, files map[string]io.Reader) error {
	fmt.Println("In s3 uploader")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	fmt.Println(sess)

	region, err := s3manager.GetBucketRegion(context.Background(), sess, bucket, "us-west-2")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(region)

	svc := s3manager.NewUploader(sess)
	fmt.Println(svc)
	
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