package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	gcs "cloud.google.com/go/storage"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

const OBJECT_KEY_TMPL = "%v/%v/image.%v"

func SaveProfilePic(cfg *config.Config, user *model.User) (bool, model.ServiceError) {
	data, err := base64.StdEncoding.DecodeString(user.CtPpic)
	if err != nil {
		return false, model.ImageDecodingError.WithCause(err)
	}

	ctx, client, bucket, serr := findBucket(cfg)
	if serr.IsError() {
		return false, serr
	}
	defer client.Close()

	objectKey := fmt.Sprintf(OBJECT_KEY_TMPL, user.CtUser, user.CtProf, user.CtImgt)
	opic := bucket.Object(objectKey)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	w := opic.NewWriter(ctx)
	if _, err = w.Write(data); err != nil {
		serr = model.CloudStorageError.WithCause(err)
		util.LogIt("Cloudtacts", fmt.Sprintf("Failed to save pic to cloud storage: %v", serr))
		return false, serr
	}
	user.CtPpic = fmt.Sprintf("%v%v", model.OBJK_TAG, objectKey)

	if err := w.Close(); err != nil {
		serr = model.CloudStorageError.WithCause(err)
		util.LogIt("Cloudtacts", fmt.Sprintf("Failed to close object writer: %v", serr))
		return false, serr
	}

	return true, model.NoError
}

func DeleteProfilePic(cfg *config.Config, user *model.User) (bool, model.ServiceError) {
	ctx, client, bucket, serr := findBucket(cfg)
	if !serr.IsError() {
		defer client.Close()

		objectKey := fmt.Sprintf(OBJECT_KEY_TMPL, user.CtUser, user.CtProf, user.CtImgt)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		err := bucket.Object(objectKey).Delete(ctx)
		if err != nil {
			serr = model.CloudStorageError.WithCause(err)
		}
	}

	return !serr.IsError(), serr
}

func ReadProfilePic(cfg *config.Config, imageKey string) ([]byte, model.ServiceError) {
	var ppic []byte

	ctx, client, bucket, serr := findBucket(cfg)
	if !serr.IsError() {
		defer client.Close()

		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		obj := bucket.Object(imageKey)
		if attrs, err := obj.Attrs(ctx); err != nil {
			serr = model.CloudStorageError.WithCause(err)
		} else {
			ppic = make([]byte, attrs.Size)
		}

		reader, err := obj.NewReader(ctx)
		if err != nil {
			serr = model.CloudStorageError.WithCause(err)
		} else {
			if _, err = reader.Read(ppic); err != nil {
				serr = model.CloudStorageError.WithCause(err)
			}
		}
	}

	return ppic, serr
}

func GetEncodedImage(cfg *config.Config, imageKey string) (string, model.ServiceError) {
	var serr model.ServiceError
	var encImg string

	var bbuff []byte
	if strings.HasPrefix(imageKey, model.OBJK_TAG) {
		bbuff, serr = ReadProfilePic(cfg, imageKey[2:])
	} else {
		bbuff, serr = ReadProfilePic(cfg, imageKey)
	}
	encImg = base64.StdEncoding.EncodeToString(bbuff)

	return encImg, serr
}

func findBucket(cfg *config.Config) (context.Context, *gcs.Client, *gcs.BucketHandle, model.ServiceError) {
	serr := model.NoError

	var ctx context.Context
	if ctx = cfg.Context(); ctx == nil {
		ctx = context.Background()
	}

	bucketName := cfg.ValueOf(model.KEY_STORAGE_BUCKET)
	var bucket *gcs.BucketHandle

	client, err := gcs.NewClient(ctx)
	if err != nil {
		serr = model.CloudStorageError.WithCause(err)
		util.LogIt("Cloudtacts", fmt.Sprintf("Failed to create cloud storage client: %v", serr))
	} else {
		bucket = client.Bucket(bucketName)
	}

	return ctx, client, bucket, serr
}
