package go_daemon_protobuf

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// protobuf_message "google.golang.org/protobuf/message"

	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

func TestPrototest(t *testing.T) {
	var sample = "idfa\t1rfw452y52g2gq4g\t55.55\t42.42\t1423,43,567,3,7,23\ngaid\t7rfw452y52g2gq4g\t55.55\t42.42\t7423,424"
	sample_splitted := strings.Split(sample, "\n")
	expected_packed := []string{"apps:1423 apps:43 apps:567 apps:3 apps:7 apps:23 lat:55.55 lon:42.42",
		"apps:7423 apps:424 lat:55.55 lon:42.42"}

	for i_part := 0; i_part < len(sample_splitted); i_part += 1 {
		sample_part_array := strings.Split(strings.Trim(sample_splitted[i_part], " "), "\t")
		var lat, _ = strconv.ParseFloat(sample_part_array[2], 8)
		var lon, _ = strconv.ParseFloat(sample_part_array[3], 8)
		var raw_apps = strings.Split(sample_part_array[4], ",")
		var apps []uint32
		for i_app := 0; i_app < len(raw_apps); i_app += 1 {
			if valid.IsInt(raw_apps[i_app]) {
				ii_app, err := strconv.Atoi(raw_apps[i_app])
				if err == nil {
					apps = append(apps, uint32(ii_app))
				}
			}
		}
		var ua = UserApps{}
		ua.Apps = apps
		ua.Lat = &lat
		ua.Lon = &lon
		var packed = ua.String()
		assert.Equal(t, expected_packed[i_part], packed, "compare packed data")
	}
}
