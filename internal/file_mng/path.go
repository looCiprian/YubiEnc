package file_mng

func AddSlashIfNotPresent(destination string) string {
	if destination[len(destination)-1:] != "/" {
		return destination + "/"
	}
	return destination
}
