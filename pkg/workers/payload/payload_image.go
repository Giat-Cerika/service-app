package payload

import "github.com/google/uuid"

type ImageUploadPayload struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Type      string    `json:"type"`                 // "single" atau "many"
	FileBytes []byte    `json:"file_bytes,omitempty"` // isi file binary
	Folder    string    `json:"folder"`               // folder di Cloudinary
	Filename  string    `json:"filename"`             // nama file upload
}

// Implementasi interface uploadable
func (p *ImageUploadPayload) GetFileBytes() []byte { return p.FileBytes }
func (p *ImageUploadPayload) GetFolder() string    { return p.Folder }
func (p *ImageUploadPayload) GetFilename() string  { return p.Filename }
func (p *ImageUploadPayload) GetType() string      { return p.Type }
