package types

type WriteDataRequest struct {
	UserID   int64
	Filepath string
	Datatype string
}

const (
	Photo = "photo"
	Video = "video"
)

func NewWriteDataRequest(userID int64, filepath, datatype string) *WriteDataRequest {
	return &WriteDataRequest{
		UserID:   userID,
		Filepath: filepath,
		Datatype: datatype,
	}
}
