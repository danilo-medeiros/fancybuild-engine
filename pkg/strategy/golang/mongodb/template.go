package mongodb

func mapSort(sort string) int {
	switch sort {
	case "asc":
		return 1
	case "desc":
		return -1
	}
	return 1
}
