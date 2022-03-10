package gocuke

type Attachment struct {
	Body            string
	ContentEncoding AttachmentContentEncoding
	FileName        string
	MediaType       string
	Url             string
}
type AttachmentContentEncoding string

const (
	AttachmentContentEncoding_IDENTITY AttachmentContentEncoding = "IDENTITY"
	AttachmentContentEncoding_BASE64   AttachmentContentEncoding = "BASE64"
)
