package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	gcs "cloud.google.com/go/storage"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

const OBJECT_KEY_TMPL = "%v/%v/image.%v"

func SaveProfilePic(cfg *config.Config, user *model.User) (bool, model.ServiceError) {
	var serr model.ServiceError
	var ctx context.Context
	if ctx = cfg.Context(); ctx == nil {
		ctx = context.Background()
	}
	bucketName := cfg.ValueOf(model.KEY_STORAGE_BUCKET)

	data, err := base64.StdEncoding.DecodeString(user.CtPpic)
	if err != nil {
		return false, model.ImageDecodingError.WithCause(err)
	}

	client, err := gcs.NewClient(ctx)
	if err != nil {
		serr = model.CloudStorageError.WithCause(err)
		util.LogIt("Cloudtacts", fmt.Sprintf("Failed to create cloud storage client: %v", serr))
		return false, serr
	}
	defer client.Close()

	objectKey := fmt.Sprintf(OBJECT_KEY_TMPL, user.CtUser, user.CtProf, user.CtImgt)
	bucket := client.Bucket(bucketName)
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
	var ctx context.Context
	if ctx = cfg.Context(); ctx == nil {
		ctx = context.Background()
	}

	bucketName := cfg.ValueOf(model.KEY_STORAGE_BUCKET)
	serr := model.NoError

	client, err := gcs.NewClient(ctx)
	if err != nil {
		serr = model.CloudStorageError.WithCause(err)
		util.LogIt("Cloudtacts", fmt.Sprintf("Failed to create cloud storage client: %v", serr))
		return false, serr
	}
	defer client.Close()

	objectKey := fmt.Sprintf(OBJECT_KEY_TMPL, user.CtUser, user.CtProf, user.CtImgt)
	bucket := client.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = bucket.Object(objectKey).Delete(ctx)
	if err != nil {
		serr = model.CloudStorageError.WithCause(err)
	}

	return true, serr
}
