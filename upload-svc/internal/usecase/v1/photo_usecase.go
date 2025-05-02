package v1

// func (u *photoUsecase) UploadPhoto(ctx context.Context, file *multipart.FileHeader, request *model.CreatePhotoRequest) error {
// 	if request.Latitude != nil && request.Longitude == nil {
// 		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Longitude is required")
// 	} else if request.Latitude == nil && request.Longitude != nil {
// 		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Latitude is required")
// 	}

// 	srcFile, err := file.Open()
// 	if err != nil {
// 		return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "Cannot open uploaded file")
// 	}
// 	defer srcFile.Close()

// 	// Peek buffer from pool
// 	peekBuf := peekBufPool.Get().([]byte)
// 	defer peekBufPool.Put(peekBuf)

// 	n, err := io.ReadFull(srcFile, peekBuf)
// 	if err != nil && err != io.ErrUnexpectedEOF {
// 		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Cannot read image header", err)
// 	}

// 	// Decode metadata
// 	imgConfig, format, err := image.DecodeConfig(bytes.NewReader(peekBuf[:n]))
// 	if err != nil {
// 		return helper.NewUseCaseWithInternalError(errorcode.ErrInvalidArgument, "Not a valid image", err)
// 	}

// 	imageType := strings.ToUpper(format)
// 	if imageType == "JPEG" {
// 		imageType = "JPG"
// 	}

// 	// Hash while streaming to upload
// 	hasher := sha256.New()
// 	multiReader := io.MultiReader(bytes.NewReader(peekBuf[:n]), srcFile)
// 	teeReader := io.TeeReader(multiReader, hasher)

// 	upload, err := u.storageAdapter.UploadFileWithoutMultipart(ctx, file, io.NopCloser(teeReader), "photo")
// 	if err != nil {
// 		return helper.WrapInternalServerError(u.logs, "error when uploading file", err)
// 	}

// 	checksum := fmt.Sprintf("%x", hasher.Sum(nil))
// 	now := time.Now()

// 	// Buat entitas Photo baru
// 	newPhoto := &entity.Photo{
// 		Id:            ulid.Make().String(),
// 		UserId:        request.UserId,
// 		CreatorId:     request.CreatorId,
// 		Title:         upload.Filename,
// 		CollectionUrl: upload.URL,
// 		Price:         request.Price,
// 		PriceStr:      request.PriceStr,
// 		OriginalAt:    time.Now(),
// 		CreatedAt:     time.Now(),
// 		UpdatedAt:     time.Now(),
// 		Latitude:      nullable.ToSQLFloat64(request.Latitude),
// 		Longitude:     nullable.ToSQLFloat64(request.Longitude),
// 		Description:   nullable.ToSQLString(request.Description),
// 	}

// 	// Lanjutkan ke database insert, dll

// 	u.logs.CustomLog("Decoded image format:", format)
// 	u.logs.Log(fmt.Sprintf("Decoded image resolution: %d * %d", imgConfig.Width, imgConfig.Height))

// 	newPhotoDetail := &entity.PhotoDetail{
// 		Id:              ulid.Make().String(),
// 		PhotoId:         newPhoto.Id,
// 		FileName:        upload.Filename,
// 		FileKey:         upload.FileKey,
// 		Size:            upload.Size,
// 		Type:            imageType,
// 		Checksum:        checksum,
// 		Width:           imgConfig.Width,  // disesuaikan tipe data jika perlu
// 		Height:          imgConfig.Height, // disesuaikan tipe data jika perlu
// 		Url:             upload.URL,
// 		YourMomentsType: enum.YourMomentTypeCollection,
// 		CreatedAt:       time.Now(),
// 		UpdatedAt:       time.Now(),
// 	}

// 	if err := u.photoAdapter.CreatePhoto(ctx, newPhoto, newPhotoDetail); err != nil {
// 		return helper.WrapInternalServerError(u.logs, "failed to create photo :", err)
// 	}

// 	// Async compression & second upload
// 	go func() {
// 		tmpFilePath, err := u.compressAdapter.CompressImageToTempFile(file.Filename, io.NopCloser(io.MultiReader(bytes.NewReader(peekBuf[:n]), srcFile)))
// 		if err != nil {
// 			u.logs.CustomError("failed to compress image: %v", err)
// 			return
// 		}
// 		defer os.Remove(tmpFilePath)

// 		fileComp, err := os.Open(tmpFilePath)
// 		if err != nil {
// 			u.logs.CustomError("failed to open compressed file: %v", err)
// 			return
// 		}
// 		defer fileComp.Close()

// 		stat, _ := fileComp.Stat()
// 		header := &multipart.FileHeader{
// 			Filename: stat.Name(),
// 			Header:   textproto.MIMEHeader{},
// 			Size:     stat.Size(),
// 		}
// 		header.Header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, stat.Name()))

// 		compressedPhoto, err := u.storageAdapter.UploadFile(ctx, header, fileComp, "photo/compressed")
// 		if err != nil {
// 			u.logs.CustomError("failed to upload compressed file: %v", err)
// 			return
// 		}

// 		compressedPhotoDetail := &entity.PhotoDetail{
// 			Id:              ulid.Make().String(),
// 			PhotoId:         newPhoto.Id,
// 			FileName:        compressedPhoto.Filename,
// 			FileKey:         compressedPhoto.FileKey,
// 			Size:            compressedPhoto.Size,
// 			Url:             compressedPhoto.URL,
// 			YourMomentsType: enum.YourMomentTypeCompressed,
// 			CreatedAt:       now,
// 			UpdatedAt:       now,
// 		}

// 		if err := u.photoAdapter.UpdatePhotoDetail(ctx, compressedPhotoDetail); err != nil {
// 			u.logs.CustomError("failed to update compressed photo detail: %v", err)
// 			return
// 		}

// 		u.logs.CustomLog("compressed uploaded: %s", compressedPhoto.URL)
// 		u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, compressedPhoto.URL)
// 	}()

// 	return nil
// |}
