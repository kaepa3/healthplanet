package healthplanet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type HealthPlanetClient struct {
	client *http.Client
	ctx    context.Context
	token  *oauth2.Token
}

type Status int

const (
	Innerscan = iota
	Sphygmomanometer
	Pedometer
)

const (
	apiRoot = "https://www.healthplanet.jp/status/"
)

type DateType int

const (
	Registration = iota
	Measuring
)

type FormatType int

const (
	Xml = iota
	Json
)

type TagType int

const (
	Weight = iota
	Fat
	PressureHigh //622E : 最高血圧 (mmHg)
	PressureLow  ///622F : 最低血圧 (mmHg)
	Pulse        //6230 : 脈拍 (bpm)
	Steps        // 6331 : 歩数 (歩)
)

type HealthPlanetOption struct {
	Date   DateType
	Format FormatType
	From   time.Time
	To     time.Time
	Tags   []TagType
}

func (c *HealthPlanetClient) Get(st Status, opt *HealthPlanetOption) (*http.Response, error) {
	status, err := getStatusStr(st)
	if err != nil {
		return nil, err
	}
	date := getDateType(opt.Date)
	tagStr := getTag(opt.Tags)
	from := getFromTo(opt.From, true)
	to := getFromTo(opt.To, false)

	request := apiRoot + status + getFormatType(opt.Format) + "?access_token=" + c.token.AccessToken + date + tagStr + from + to
	fmt.Println(request)
	return c.client.Get(request)
}
func getFromTo(t time.Time, isFrom bool) string {
	if time.Now().IsZero() {
		return ""
	}
	opt := "&to="
	if isFrom {
		opt = "&from="
	}
	return opt + t.Format("20060102150405")
}

func getTag(tags []TagType) string {
	if len(tags) == 0 {
		return ""
	}
	tagList := make([]string, len(tags))
	for i, v := range tags {
		tagList[i] = tagToStr(v)
	}
	return "&tag=" + strings.Join(tagList, ",")
}

func tagToStr(tag TagType) string {
	tagStr := ""
	switch tag {
	case Weight:
		tagStr = "6021"
	case Fat:
		tagStr = "6022"
	case PressureHigh:
		tagStr = "602E"
	case PressureLow:
		tagStr = "602F"
	case Pulse:
		tagStr = "6230"
	case Steps:
		tagStr = "6331"
	}
	return tagStr
}

func getDateType(d DateType) string {
	date := "0"

	if d == Measuring {
		date = "1"
	}
	return "&date=" + date
}

func getStatusStr(st Status) (string, error) {
	if st == Innerscan {
		return "innerscan", nil
	} else if st == Sphygmomanometer {
		return "sphygmomanometer", nil
	} else if st == Pedometer {
		return "pedometer", nil
	}
	return "", errors.New("invalid value")
}

func getFormatType(t FormatType) string {
	if t == Xml {
		return ".xml"
	}
	return ".json"
}
