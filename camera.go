package gp

// #cgo linux pkg-config: libgphoto2
// #include <gphoto2/gphoto2.h>
// #include <string.h>
import "C"
import (
	"unsafe"
	"fmt"
)

const (
	CAPTURE_IMAGE = C.GP_CAPTURE_IMAGE
	CAPTURE_MOVIE = C.GP_CAPTURE_MOVIE
	CAPTURE_SOUND = C.GP_CAPTURE_SOUND
)

type Camera C.Camera
type CameraCaptureType int

func NewCamera() (*Camera, error) {
	var _cam *C.Camera

	if ret := C.gp_camera_new(&_cam); ret != 0 {
		return nil, e(ret)
	}

	return (*Camera)(_cam), nil
}

func (camera *Camera) Init(ctx *Context) error {
	if ret := C.gp_camera_init(camera.c(), ctx.c()); ret != 0 {
		return e(ret)
	}
	return nil
}

func (camera *Camera) CapturePreview(ctx *Context) ([]byte, error) {
	var _file C.CameraFile
	if ret := C.gp_camera_capture_preview(camera.c(), &_file, ctx.c()); ret != 0 {
		fmt.Printf("Error gp_camera_capture_preview: %s\n", e(ret))
		return nil, e(ret)
	}

	var data *byte
	_data := (*C.char)(unsafe.Pointer(data))
	var size uint64
	_size := C.ulong(size)

	if ret := C.gp_file_get_data_and_size(&_file, &_data, &_size); ret != 0 {
		return nil, e(ret)
	}

	return C.GoBytes(unsafe.Pointer(_data), C.int(_size)), nil
}

func (camera *Camera) CaptureImage(ctx *Context) ([]byte, error) {
	var path CameraFilePath
	var _path C.CameraFilePath

	_captureType := C.CameraCaptureType(C.GP_CAPTURE_IMAGE)
	if ret := C.gp_camera_capture(camera.c(), _captureType, &_path, ctx.c()); ret != 0 {
		return nil , e(ret)
	}
	path.Name = C.GoString(&_path.name[0])
	path.Folder = C.GoString(&_path.folder[0])

	var _file *C.CameraFile
	C.gp_file_new(&_file)
	_filetype := (C.CameraFileType)(FILE_TYPE_NORMAL)
	if ret := C.gp_camera_file_get(camera.c(), &_path.folder[0], &_path.name[0], _filetype, _file, ctx.c()); ret != 0 {
		return nil, e(ret)
	}

	var data *byte
	_data := (*C.char)(unsafe.Pointer(data))
	var size uint64
	_size := C.ulong(size)

	if ret := C.gp_file_get_data_and_size(_file, &_data, &_size); ret != 0 {
		return nil, e(ret)
	}

	return C.GoBytes(unsafe.Pointer(_data), C.int(_size)), nil
}

func (camera *Camera) File(folder, name string, filetype CameraFileType, context *Context) (*CameraFile, error) {
	var _file *C.CameraFile
	C.gp_file_new(&_file)

	_camera := (*C.Camera)(unsafe.Pointer(camera))
	_folder := C.CString(folder)
	_name := C.CString(name)
	_context := (*C.GPContext)(unsafe.Pointer(context))
	_filetype := (C.CameraFileType)(filetype)
	if ret := C.gp_camera_file_get(_camera, _folder, _name, _filetype, _file, _context); ret != 0 {
		return nil, e(ret)
	}

	return (*CameraFile)(unsafe.Pointer(_file)), nil
}



func (camera *Camera) Free() error {
	if ret := C.gp_camera_free(camera.c()); ret != 0 {
		return e(ret)
	}
	return nil
}

func (camera *Camera) c() *C.Camera {
	return (*C.Camera)(camera)
}

func (path *CameraFilePath) c() *C.CameraFilePath {
	return (*C.CameraFilePath)(unsafe.Pointer(path))
}