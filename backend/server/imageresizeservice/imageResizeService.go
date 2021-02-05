package imageresizeservice

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"net/http"

	"path/filepath"

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

func (ir *imageResizeService) uploadImage(image io.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("nt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Println(handler.Header.Get("Content-type"))
	}
}

func (ir *imageResizeService) uuidPath(r io.Reader, rootDir string) string {
	uid := pathFromUUID(ir.newUUID(r))
	return filepath.Join(rootDir, uid)
}

func (ir *imageResizeService) saveImageAllSizes(image io.Reader, rootDir string) error {
	var pipeBuf bytes.Buffer
	uuidReader := io.TeeReader(image, &pipeBuf)

	path := ir.uuidPath(uuidReader, rootDir)


	readers, writers := createPipesForBreakpoints(ir.breakpoints) 
	done := make(chan error)

	for size, bkpt := range ir.breakpoints {
		go func(size string, bkpt breakpoint) {
			done <- ir.saveImageAtSize(readers[size], path, size, bkpt)
		}(size, bkpt)
	}

	go func() {
		for size := range ir.breakpoints {
			defer writers[size].Close()
		}

		writerSlice := getPipeWriterMapAsWriterSlice(writers)
		mw := io.MultiWriter(writerSlice...)

		io.Copy(mw, &pipeBuf)
	}()

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

func (ir *imageResizeService) saveImageAtSize(image io.Reader, rootDir, sizeSuffix string, b breakpoint) error {
	path := filepath.Join(rootDir, sizeSuffix)

	image = resizeImage(image, b)

	return ir.writer.writeFile(path, image)
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