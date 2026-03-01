package infrastructure

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImageWriter_SaveImage(t *testing.T) {
	tempDir := t.TempDir()
	writer := NewImageWriter()

	tests := []struct {
		name        string
		createImage func() *image.RGBA
		filename    string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid RGBA image",
			createImage: func() *image.RGBA {
				img := image.NewRGBA(image.Rect(0, 0, 10, 10))
				for y := 0; y < 10; y++ {
					for x := 0; x < 10; x++ {
						img.Set(x, y, color.RGBA{
							R: uint8(x * 25),
							G: uint8(y * 25),
							B: uint8((x + y) * 12),
							A: 255,
						})
					}
				}
				return img
			},
			filename: "test_rgba.png",
			wantErr:  false,
		},
		{
			name: "valid grayscale image",
			createImage: func() *image.RGBA {
				img := image.NewRGBA(image.Rect(0, 0, 5, 5))
				for y := 0; y < 5; y++ {
					for x := 0; x < 5; x++ {
						img.Set(x, y, color.Gray{Y: uint8((x + y) * 12)})
					}
				}
				return img
			},
			filename: "test_gray.png",
			wantErr:  false,
		},
		{
			name: "empty image",
			createImage: func() *image.RGBA {
				return image.NewRGBA(image.Rect(0, 0, 0, 0))
			},
			filename:    "empty.png",
			wantErr:     true,
			errContains: "failed to encode PNG:",
		},
		{
			name: "invalid directory",
			createImage: func() *image.RGBA {
				return image.NewRGBA(image.Rect(0, 0, 1, 1))
			},
			filename:    "/nonexistent/path/image.png",
			wantErr:     true,
			errContains: "failed to create file",
		},
		{
			name: "filename without extension",
			createImage: func() *image.RGBA {
				return image.NewRGBA(image.Rect(0, 0, 1, 1))
			},
			filename: "no_extension",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := tt.createImage()
			filepath := filepath.Join(tempDir, tt.filename)

			err := writer.SaveImage(img, filepath)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q doesn't contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Errorf("SaveImage() didn't create file %s", filepath)
			}

			file, err := os.Open(filepath)
			if err != nil {
				t.Errorf("Failed to open saved file: %v", err)
				return
			}
			defer file.Close()

			pngHeader := []byte{137, 80, 78, 71, 13, 10, 26, 10}
			header := make([]byte, 8)
			n, err := file.Read(header)
			if err != nil || n != 8 {
				t.Errorf("Failed to read file header: %v", err)
				return
			}

			for i := 0; i < 8; i++ {
				if header[i] != pngHeader[i] {
					t.Errorf("Invalid PNG header at byte %d: got %d, want %d", i, header[i], pngHeader[i])
					break
				}
			}
		})
	}
}

func TestImageWriter_SaveImage_Overwrite(t *testing.T) {
	tempDir := t.TempDir()
	writer := NewImageWriter()
	filepath := filepath.Join(tempDir, "overwrite_test.png")

	img1 := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img1.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	err := writer.SaveImage(img1, filepath)
	if err != nil {
		t.Fatalf("First SaveImage() error = %v", err)
	}

	img2 := image.NewRGBA(image.Rect(0, 0, 5, 5))
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			img2.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	err = writer.SaveImage(img2, filepath)
	if err != nil {
		t.Fatalf("Second SaveImage() error = %v", err)
	}
}

func TestImageWriter_SaveImage_Permissions(t *testing.T) {
	tempDir := t.TempDir()
	writer := NewImageWriter()

	readOnlyFile := filepath.Join(tempDir, "readonly.png")
	file, err := os.Create(readOnlyFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	err = os.Chmod(readOnlyFile, 0444)
	if err != nil {
		t.Fatalf("Failed to set read-only permissions: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	err = writer.SaveImage(img, readOnlyFile)

	if err == nil {
		t.Error("SaveImage() should fail for read-only file")
	}

	if !strings.Contains(err.Error(), "failed to create file") {
		t.Errorf("error message %q should contain 'failed to create file'", err.Error())
	}
}

func TestNewImageWriter(t *testing.T) {
	writer := NewImageWriter()

	if writer == nil {
		t.Fatal("NewImageWriter() returned nil")
	}
}
