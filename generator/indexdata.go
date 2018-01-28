package generator

type IndexData struct {
	IndexPath string
}

func NewIndexData(indexPath string) (indexData *IndexData) {
	return &IndexData{IndexPath: indexPath}
}
