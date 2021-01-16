package passtime

import (
	"github.com/rs/zerolog/log"
	"os"
	"syscall"
	"time"
)

func PropsOs(props map[Name]Prop, filepath string) (map[Name]Prop, error) {
	fi, err := os.Stat(filepath)

	if err != nil {
		log.Fatal().Err(err).Msg("stat")
		return props, err
	}

	stat := fi.Sys().(*syscall.Stat_t)

	//t := time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec)).Format("2006-01-02 15:04:05")
	//prop := newProp(ATIME, TYPE_TIME, t)
	//props[ATIME] = prop

	t := time.Unix(int64(stat.Birthtimespec.Sec), int64(stat.Birthtimespec.Nsec)).Format("2006-01-02 15:04:05")
	prop := newProp(BTIME, TYPE_TIME, t)
	props[BTIME] = prop

	//t = time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)).Format("2006-01-02 15:04:05")
	//prop = newProp(CTIME, TYPE_TIME, t)
	//props[CTIME] = prop

	return props, nil
}
