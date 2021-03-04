package imageresizeservice

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"path/filepath"
	"strings"
	"io"
	"io/ioutil"

	uuid "github.com/satori/go.uuid"

	// depends on libvips-dev
	"github.com/daddye/vips"
	//"github.com/go-chi/chi"
)

var defaultBreakpoints = breakpointMap{
	"xs": {200, 75},
	"sm": {600, 75},
	"md": {960, 80},
	"lg": {1280, 90},
	"xl": {1920, 95},
}

type imageResizeService struct {
	breakpoints breakpointMap
	namespace   uuid.UUID
	writer      fileWriter
}

type breakpointMap map[string]breakpoint

type breakpoint struct {
	size    int
	quality int
}

type fileWriter interface {
	writeFile(string, io.Reader) error
}

func newImageResizeService(breakpoints breakpointMap, namespace string, writer fileWriter) *imageResizeService {
	return &imageResizeService{
		breakpoints: breakpoints,
		namespace:   uuid.NewV5(uuid.Nil, namespace),
		writer:      writer,
	}
}


// CreateImageHTTP saves an image at a unique address
func CreateImageHTTP(rootDir, namespace string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		s := newS3FileUploader("cm-personal-site-v2-bucket", len(defaultBreakpoints))
		//writer := writer{}
		ir := newImageResizeService(defaultBreakpoints, namespace, s)

		path, err := ir.saveImageAllSizesUUID(file, rootDir)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.getUploadErr()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("No errors uploading to ", path)

		jsonResp := struct {
			Path string `json:"path"`
		}{
			Path: path,
		}

		fmt.Println(jsonResp)

		//w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResp)
	}
}

// CreateStaticImageHTTP saves an image at the given name
func CreateStaticImageHTTP(rootDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		s := newS3FileUploader("cm-personal-site-v2-bucket", len(defaultBreakpoints))
		//writer := writer{}
		ir := newImageResizeService(defaultBreakpoints, "", s)

		path, err := ir.saveImageAllSizesStatic(file, fileHeader.Filename)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.getUploadErr()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("No errors uploading to ", path)

		jsonResp := struct {
			Path string `json:"path"`
		}{
			Path: path,
		}

		fmt.Println(jsonResp)

		//w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResp)
	}
}

func (ir *imageResizeService) newUUID(r io.Reader) uuid.UUID {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil
	}

	encoded := base32.StdEncoding.EncodeToString(content)

	return uuid.NewV5(ir.namespace, encoded)
}

func pathFromUUID(u uuid.UUID) string {
	id := u.String()
	return filepath.Join(id[0:2], id[2:4], id[4:])
}

func (ir *imageResizeService) uuidPath(r io.Reader, rootDir string) string {
	uid := pathFromUUID(ir.newUUID(r))
	return filepath.Join(rootDir, uid)
}

func (ir *imageResizeService) saveImageAllSizesUUID(image io.Reader, rootDir string) (string, error) {
	var imgBuf bytes.Buffer
	uuidReader := io.TeeReader(image, &imgBuf)

	path := ir.uuidPath(uuidReader, rootDir)

	return ir.saveImageAllSizes(imgBuf, path)
}

func (ir *imageResizeService) saveImageAllSizesStatic(image io.Reader, path string) (string, error) {
	var imgBuf bytes.Buffer
	imgBuf.ReadFrom(image)

	return ir.saveImageAllSizes(imgBuf, path)
}

func (ir *imageResizeService) saveImageAllSizes(imgBuf bytes.Buffer, path string) (string, error) {
	readers, writers := createPipesForBreakpoints(ir.breakpoints)
	saveSuccess := make(chan error)

	ir.resizeAndSaveFromReaders(readers, path, saveSuccess)
	go ir.copyWritersToReaders(writers, imgBuf)

	err := ir.blockUntilAllSizesSaved(saveSuccess)
	if err != nil {
		return "", err
	}

	return path, nil
}

// saveAllSizes launches separate go routines to resize and save all sizes at once
func (ir *imageResizeService) resizeAndSaveFromReaders(readers map[string]*io.PipeReader, path string, done chan error) {
	for size, bkpt := range ir.breakpoints {
		go func(size string, bkpt breakpoint) {
			done <- ir.saveImageAtSize(readers[size], path, size, bkpt)
		}(size, bkpt)
	}
}

func (ir *imageResizeService) copyWritersToReaders(writers map[string]*io.PipeWriter, imgBuf bytes.Buffer) {
	for size := range ir.breakpoints {
		defer writers[size].Close()
	}

	writerSlice := getPipeWriterMapAsWriterSlice(writers)
	mw := io.MultiWriter(writerSlice...)

	io.Copy(mw, &imgBuf)
}

func (ir *imageResizeService) blockUntilAllSizesSaved(done chan error) error {
	for range ir.breakpoints {
		err := <-done
		if err != nil {
			return err
		}
	}

	return nil
}

func createPipesForBreakpoints(b breakpointMap) (map[string]*io.PipeReader, map[string]*io.PipeWriter) {
	ioReaders := make(map[string]*io.PipeReader, len(b))
	ioWriters := make(map[string]*io.PipeWriter, len(b))

	for i := range b {
		ioReaders[i], ioWriters[i] = io.Pipe()
	}

	return ioReaders, ioWriters
}

func getPipeWriterMapAsWriterSlice(writers map[string]*io.PipeWriter) []io.Writer {
	writerSlice := make([]io.Writer, 0)

	for _, writer := range writers {
		writerSlice = append(writerSlice, writer)
	}

	return writerSlice
}

func (ir *imageResizeService) saveImageAtSize(image io.Reader, path, size string, b breakpoint) error {
	suffix := "-" + size
	path = insertSuffixBeforeExt(path, suffix)

	image = resizeImage(image, b)

	return ir.writer.writeFile(path, image)
}

func insertSuffixBeforeExt(file, suffix string) string {
	ext := filepath.Ext(file)
	file = strings.TrimSuffix(file, ext)

	return file + suffix + ext
}

func resizeImage(image io.Reader, b breakpoint) io.Reader {
	options := vips.Options{
		Width:        b.size,
		Height:       b.size,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      b.quality,
	}

	inBuf, _ := ioutil.ReadAll(image)
	outBuf, _ := vips.Resize(inBuf, options)

	return bytes.NewReader(outBuf)
}

type dirWriter struct{}

func (w dirWriter) writeFile(path string, r io.Reader) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}
