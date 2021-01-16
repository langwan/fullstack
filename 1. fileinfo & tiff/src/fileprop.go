package passtime

import (
	"context"
	"fmt"
	"github.com/dsoprea/go-exif"
	"github.com/flotzilla/pdf_parser"
	"github.com/h2non/filetype"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

type PropType string

const (
	TYPE_STRING PropType = "string"
	TYPE_TIME   PropType = "time"
	TYPE_LONG   PropType = "long"
	TYPE_FLOAT  PropType = "float"
)

type Name string

const (
	FILESIZE        Name = "filesize" //文件大小
	MTIME           Name = "mtime"
	ATIME           Name = "atime"
	CTIME           Name = "ctime"
	BTIME           Name = "btime"
	CAMERA          Name = "camera"
	WIDTH           Name = "width"
	HEIGHT          Name = "height"
	LENS            Name = "lens"
	FOCALLENGTH     Name = "focallength"
	ISOSPEEDRATINGS Name = "isospeedratings"
	EXPOSURETIME    Name = "exposuretime"
	FNUMBER         Name = "fnumber"
	PHOTOTIIME      Name = "phototime"
	DURATION        Name = "duration"
	AUDIO_BITRATE   Name = "audio_bitrate"
	PAGES           Name = "pages"
	TITLE           Name = "title"
)

type Prop struct {
	Name  Name
	Type  PropType
	Value interface{}
}

func Props(filepath string) (map[Name]Prop, map[Name]Prop, error) {

	props := make(map[Name]Prop)
	fi, err := os.Stat(filepath)

	prop := newProp(FILESIZE, TYPE_LONG, fi.Size())
	props[FILESIZE] = prop

	prop = newProp(MTIME, TYPE_TIME, fi.ModTime().Format("2006-01-02 15:04:05"))
	props[MTIME] = prop

	if err != nil {
		log.Fatal().Err(err).Msg("stat")
		return nil, nil, err
	}

	props, err = PropsOs(props, filepath)

	if err != nil {
		log.Fatal().Err(err).Send()
		return nil, nil, err
	}

	mps := make(map[Name]Prop, 0)

	mps, err = moreProps(mps, filepath)

	return props, mps, nil
}

func newProp(name Name, tp PropType, value interface{}) (prop Prop) {
	prop.Name = name
	prop.Type = tp
	prop.Value = value
	return prop
}

func moreProps(props map[Name]Prop, filepath string) (map[Name]Prop, error) {

	file, _ := os.Open(filepath)
	defer file.Close()
	head := make([]byte, 261)

	file.Read(head)

	if filetype.IsImage(head) {

		props, _ = imageProps(props, filepath)
	} else if filetype.IsVideo(head) {
		props, _ = videoProps(props, filepath)
	} else if filetype.IsExtension(head, "pdf") {

		props, _ = pdfProps(props, filepath)
	}



	return props, nil
}

func pdfProps(props map[Name]Prop, filepath string) (map[Name]Prop, error) {

	pdf, _ := pdf_parser.ParsePdf(filepath)

	prop := Prop{
		Name:  PAGES,
		Type:  TYPE_LONG,
		Value: pdf.GetPagesCount(),
	}
	props[PAGES] = prop

	return props, nil
}

func videoProps(props map[Name]Prop, filepath string) (map[Name]Prop, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	data, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return props, err
	}
	duration := data.Format.Duration().Seconds()
	prop := Prop{
		Name:  DURATION,
		Type:  TYPE_FLOAT,
		Value: duration,
	}
	props[DURATION] = prop

	prop = Prop{
		Name:  WIDTH,
		Type:  TYPE_LONG,
		Value: data.FirstVideoStream().Width,
	}
	props[WIDTH] = prop

	prop = Prop{
		Name:  HEIGHT,
		Type:  TYPE_LONG,
		Value: data.FirstVideoStream().Height,
	}
	props[HEIGHT] = prop

	prop = Prop{
		Name:  AUDIO_BITRATE,
		Type:  TYPE_STRING,
		Value: data.FirstAudioStream().BitRate,
	}
	props[AUDIO_BITRATE] = prop
	return props, nil
}

func imageProps(props map[Name]Prop, filepath string) (map[Name]Prop, error) {

	file, _ := os.Open(filepath)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	imageExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return props, err
	}

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	var list []Prop
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {

		defer func() {
			if state := recover(); state != nil {

			}
		}()

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
		it, err := ti.Get(ifdPath, tagId)
		if err != nil {

		}

		valueString := ""
		var value interface{}
		if tagType.Type() == exif.TypeUndefined {
			var err error
			value, err = valueContext.Undefined()
			if err != nil {

			}

			valueString = fmt.Sprintf("%v", value)
		} else {
			valueString, err = valueContext.FormatFirst()
			if err != nil {
			}
			value = valueString
		}
		var prop Prop

		switch it.Id {
		case 0x0110:
			prop = Prop{
				Name:  CAMERA,
				Type:  TYPE_STRING,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x0100:
			prop = Prop{
				Name:  WIDTH,
				Type:  TYPE_LONG,
				Value: value,
			}
			list = append(list, prop)

			break
		case 0xa002:
			prop = Prop{
				Name:  WIDTH,
				Type:  TYPE_LONG,
				Value: value,
			}
			list = append(list, prop)

			break
		case 0x0101:
			prop = Prop{
				Name:  HEIGHT,
				Type:  TYPE_LONG,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0xa003:
			prop = Prop{
				Name:  HEIGHT,
				Type:  TYPE_LONG,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0xa434:
			prop = Prop{
				Name:  LENS,
				Type:  TYPE_STRING,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x920a:
			prop = Prop{
				Name:  FOCALLENGTH,
				Type:  TYPE_STRING,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x8827:
			prop = Prop{
				Name:  ISOSPEEDRATINGS,
				Type:  TYPE_LONG,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x829a:
			prop = Prop{
				Name:  EXPOSURETIME,
				Type:  TYPE_STRING,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x829d:
			prop = Prop{
				Name:  FNUMBER,
				Type:  TYPE_STRING,
				Value: value,
			}
			list = append(list, prop)
			break
		case 0x0132:
			prop = Prop{
				Name:  PHOTOTIIME,
				Type:  TYPE_TIME,
				Value: value,
			}
			list = append(list, prop)
			break
		default:

			return nil
		}

		return nil
	}

	_, err = exif.Visit(exif.IfdStandard, im, ti, imageExif, visitor)
	if err != nil {
		return props, err
	}
	var w int64 = 0
	var h int64 = 0
	for _, el := range list {

		if el.Name == WIDTH {
			var n int64 = 0
			if _, ok := el.Value.(int64); ok {
				n = el.Value.(int64)
			} else {
				n, _ = strconv.ParseInt(el.Value.(string), 10, 64)
			}
			if n > w {
				props[el.Name] = el
				w = n
			}
		} else if el.Name == HEIGHT {
			var n int64 = 0
			if _, ok := el.Value.(int64); ok {
				n = el.Value.(int64)
			} else {
				n, _ = strconv.ParseInt(el.Value.(string), 10, 64)
			}
			if n > h {
				props[el.Name] = el
				h = n
			}
		} else {
			props[el.Name] = el
		}
	}

	return props, nil
}

func remove(a []Prop, i int) []Prop {
	a[i] = a[len(a)-1]
	return a[:len(a)-1]
}
