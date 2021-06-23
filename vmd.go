// Package vmd 解码vmd格式的动作文件
package vmd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var shiftJisDecoder = japanese.ShiftJIS.NewDecoder()

type Header struct {
	Version int    // 1 or 2
	Model   string // 模型里应该是Shift_JIS编码, 我们转为utf-8
}

type BoneFrame struct {
	Bone       string     // 骨骼名字
	Time       uint32     // 帧序号
	Translate  [3]float32 // 移动
	RotateQuat [4]float32 // 旋转 四元数
	XCurve     [16]byte   // X 曲线
	YCurve     [16]byte   // Y 曲线
	ZCurve     [16]byte   // Z 曲线
	RCurve     [16]byte   // 旋转曲线
}

type MorphFrame struct {
	Morph  string  // morph动画名字
	Time   uint32  // 帧序号
	Weight float32 // 权重

}

type CameraFrame struct {
	Time      uint32     // 帧序号
	Distance  float32    // 距离
	Translate [3]float32 // 移动
	RotateXyz [3]float32 // 旋转xyz
	Curve     [24]byte
	ViewAngle float32
	Ortho     byte
}

type LightFrame struct {
	Time      uint32 // 帧序号
	Color     [3]float32
	Direction [3]float32
}

type VMD struct {
	Header

	BoneFrames   []BoneFrame
	MorphFrames  []MorphFrame
	CameraFrames []CameraFrame
	LightFrames  []LightFrame
}

func decodeString(r io.Reader, n int) (s string, err error) {
	if n == 0 {
		return
	}
	buf := make([]byte, n)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	if i := bytes.IndexByte(buf, 0); i != -1 {
		buf = buf[:i]
	}
	src := string(buf)
	if s, _, err = transform.String(shiftJisDecoder, src); err != nil {
		s = src
		err = nil
	}
	return
}

func (vm *VMD) decodeHeader(r io.Reader) (err error) {
	// 30 bytes magic and version
	magicStr, err := decodeString(r, 30)
	if err != nil {
		return
	}
	magicStr = magicStr[:25]
	if strings.HasPrefix(magicStr, "Vocaloid Motion Data") {
		if magicStr == "Vocaloid Motion Data file" {
			vm.Version = 1
		} else if magicStr == "Vocaloid Motion Data 0002" {
			vm.Version = 2
		} else {
			err = errors.New("unsupported vmd version")
			return
		}
	} else {
		err = errors.New("not a vmd format file")
		return
	}
	vm.Model, err = decodeString(r, vm.Version*10)
	return
}

func (vm *VMD) decodeBoneFrames(r io.Reader) (err error) {
	var numFrames uint32
	if err = binary.Read(r, binary.LittleEndian, &numFrames); err != nil {
		return
	}
	if numFrames > 0 {
		vm.BoneFrames = make([]BoneFrame, numFrames)
		for i := range vm.BoneFrames {
			if vm.BoneFrames[i].Bone, err = decodeString(r, 15); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].Time); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].Translate); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].RotateQuat); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].XCurve); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].YCurve); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].ZCurve); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.BoneFrames[i].RCurve); err != nil {
				return
			}
		}
	}
	return
}

func (vm *VMD) decodeMorphFrames(r io.Reader) (err error) {
	var numFrames uint32
	if err = binary.Read(r, binary.LittleEndian, &numFrames); err != nil {
		return
	}
	if numFrames > 0 {
		vm.MorphFrames = make([]MorphFrame, numFrames)
		for i := range vm.MorphFrames {
			if vm.MorphFrames[i].Morph, err = decodeString(r, 15); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.MorphFrames[i].Time); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.MorphFrames[i].Weight); err != nil {
				return
			}
		}
	}
	return
}

func (vm *VMD) decodeCameraFrames(r io.Reader) (err error) {
	var numFrames uint32
	if err = binary.Read(r, binary.LittleEndian, &numFrames); err != nil {
		return
	}
	if numFrames > 0 {
		vm.CameraFrames = make([]CameraFrame, numFrames)
		for i := range vm.CameraFrames {
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].Time); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].Distance); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].Translate); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].RotateXyz); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].Curve); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].ViewAngle); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.CameraFrames[i].Ortho); err != nil {
				return
			}
		}
	}
	return
}

func (vm *VMD) decodeLightFrames(r io.Reader) (err error) {
	var numFrames uint32
	if err = binary.Read(r, binary.LittleEndian, &numFrames); err != nil {
		return
	}
	if numFrames > 0 {
		vm.LightFrames = make([]LightFrame, numFrames)
		for i := range vm.LightFrames {
			if err = binary.Read(r, binary.LittleEndian, &vm.LightFrames[i].Time); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.LightFrames[i].Color); err != nil {
				return
			}
			if err = binary.Read(r, binary.LittleEndian, &vm.LightFrames[i].Direction); err != nil {
				return
			}
		}
	}
	return
}

func Decode(r io.Reader) (vm *VMD, err error) {
	vm = new(VMD)

	defer func() {
		if err != nil {
			vm = nil
		}
	}()

	if err = vm.decodeHeader(r); err != nil {
		err = fmt.Errorf("vmd: error decoding header: %w", err)
		return
	}
	if err = vm.decodeBoneFrames(r); err != nil {
		err = fmt.Errorf("vmd: error decoding bone frames: %w", err)
		return
	}
	if err = vm.decodeMorphFrames(r); err != nil {
		err = fmt.Errorf("vmd: error decoding morph frames: %w", err)
		return
	}
	if err = vm.decodeCameraFrames(r); err != nil {
		err = fmt.Errorf("vmd: error decoding camera frames: %w", err)
		return
	}
	if err = vm.decodeLightFrames(r); err != nil {
		err = fmt.Errorf("vmd: error decoding light frames: %w", err)
		return
	}

	return
}
