package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

func main() {
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true
	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("minioClient is now set up") // minioClient is now set up

	// make bucket
	bucketName := "tags-test"
	objectName := "upload.go"
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// upload with custom tags
	tagString := "OS=Ubuntu&Version=20.04&Build=10"
	tagsSet, err := tags.Parse(tagString, true)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := minioClient.FPutObject(ctx, bucketName, objectName,
		objectName, minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			UserTags:    tagsSet.ToMap(),
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	// retrieve tags
	prefix := ""
	for obj := range minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:       prefix,
		Recursive:    true,
		WithMetadata: true,
	}) {
		if obj.Err != nil {
			log.Fatalln(obj.Err)
		}
		log.Println("obj", obj.UserTags)
	}
}
