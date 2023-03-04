package big

import (
	"fmt"
	"io/fs"
)

func Files(fsys fs.FS, size int64) (int64, error) {
	var count int64
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		
		fmt.Printf("%s: [%s]\n", info.Name(), ByteCountIEC(info.Size()))

		if info.Size() >= size {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	
	return count, nil
}