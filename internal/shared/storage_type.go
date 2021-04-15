package shared

// StorageType represents a type of storage a paste can be stored with
type StorageType string

const (
	StorageTypeFile     = StorageType("file")
	StorageTypePostgres = StorageType("postgres")
	StorageTypeMongoDB  = StorageType("mongodb")
	StorageTypeS3       = StorageType("s3")
)
