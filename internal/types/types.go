package types

type WriteDataRequest struct {
	UserID   string
	Filepath string
	Datatype string
}

func NewWriteDataRequest(userID, filepath, datatype string) *WriteDataRequest {
	return &WriteDataRequest{
		UserID:   userID,
		Filepath: filepath,
		Datatype: datatype,
	}
}
