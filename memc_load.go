package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"go_daemon_protobuf"
	"log"
	"os"
	filepath "path/filepath"
	"reflect"
	"strconv"
	"strings"

	memcache "github.com/bradfitz/gomemcache/memcache"
)

var NORMAL_ERR_RATE float32 = 0.01

type AppsInstalled struct {
	dev_type string
	dev_id   string
	lat      float64
	lon      float64
	apps     []uint32
}

type arguments struct {
	l       *bool
	log     *bool
	dry     *bool
	pattern *string
	idfa    *string
	gaid    *string
	adid    *string
	dvid    *string
}

type device struct {
	idfa string
	gaid string
	adid string
	dvid string
}

func convert_pointer_struct_to_json(data arguments) (result string) {
	result = fmt.Sprintf("{'log': %v, 'dry': %v, 'pattern': %v, 'idfa': %v, 'gaid': %v, 'adid': %v, 'dvid': %v}", *data.l || *data.log, *data.dry, *data.pattern, *data.idfa, *data.gaid, *data.adid, *data.dvid)
	return result
}

func main_exec(arg_str arguments) {
	var device_memc device
	device_memc.idfa = *arg_str.idfa
	device_memc.gaid = *arg_str.gaid
	device_memc.adid = *arg_str.adid
	device_memc.dvid = *arg_str.dvid
	files, _ := filepath.Glob(*arg_str.pattern)
	for _, fn := range files {
		processed := 0
		errors := 0
		log.Println("Processing " + fn)
		reader, err := os.Open(fn)
		if err != nil {
			log.Fatalf(err.Error())
		}
		uncompressed_data, err := gzip.NewReader(reader)
		if err != nil {
			log.Fatal("can't read gzip file:" + fn)
		}

		scanner := bufio.NewScanner(uncompressed_data)
		scanner.Split(bufio.ScanLines)
		uncompressed_data.Close()
		reader.Close()

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}
			appsinstalled := parse_appsinstalled(line)
			if appsinstalled == nil {
				errors += 1
				continue
			}
			memc_addr := getField(&device_memc, appsinstalled.dev_type)
			if memc_addr.Equal(reflect.Value{}) {
				errors += 1
				log.Fatal("Unknow device type:" + appsinstalled.dev_type)
				continue
			}
			ok := insert_appsinstalled(memc_addr.String(), appsinstalled, arg_str.dry)
			if ok == true {
				processed += 1
			} else {
				errors += 1
			}
			if processed == 0 {
				dot_rename(fn)
				continue
			}
		}
		err_rate := float32(errors) / float32(processed)
		if err_rate < NORMAL_ERR_RATE {
			log.Println(fmt.Sprintf("Acceptable error rate (%v). Successfull load", err_rate))
		} else {
			log.Println(fmt.Sprintf("High error rate (%v > %v). Failed load", err_rate, NORMAL_ERR_RATE))
		}
		dot_rename(fn)
	}
}

func insert_appsinstalled(memc_addr string, appsinstalled *AppsInstalled, dry_run *bool) (result bool) {
	var ua = go_daemon_protobuf.UserApps{}
	ua.Lat = &appsinstalled.lat
	ua.Lon = &appsinstalled.lon
	key := fmt.Sprintf("%v:%v", appsinstalled.dev_type, appsinstalled.dev_id)
	ua.Apps = appsinstalled.apps
	var packed = ua.String()
	if *dry_run == true {
		log.Println(fmt.Sprintf("%v - %v -> %v", memc_addr, key, strings.ReplaceAll(packed, "\n", " ")))
	} else {
		memc := memcache.New(memc_addr)
		// log.Println("Save key in memcache: " + key)
		err := memc.Set(&memcache.Item{Key: key, Value: []byte(packed)})
		if err != nil {
			log.Fatalln(fmt.Sprintf("Cannot write to memc %v: %v", memc_addr, err.Error()))
			return false
		}
	}
	return true
}

func dot_rename(old_filename string) {
	file := filepath.Base(old_filename)
	dir := filepath.Dir(old_filename)
	// time.Sleep(2 * time.Second)
	err := os.Rename(old_filename, filepath.Join(dir, "."+file))
	if err != nil {
		log.Fatal(err)
	}
}

func getField(v *device, field string) (result reflect.Value) {
	r := reflect.ValueOf(v)
	result = reflect.Indirect(r).FieldByName(field)
	return result
}

func parse_appsinstalled(line string) (app_install *AppsInstalled) {
	line_parts := strings.Split(strings.TrimSpace(line), "\t")
	if len(line_parts) < 5 {
		return nil
	}
	dev_type := line_parts[0]
	dev_id := line_parts[1]
	lat := line_parts[2]
	lon := line_parts[3]
	raw_apps := line_parts[4]
	if dev_type == "" || dev_id == "" {
		return nil
	}
	var apps []uint32
	for _, i_app := range strings.Split(raw_apps, ",") {
		i_app = strings.TrimSpace(i_app)
		if _, err := strconv.Atoi(i_app); err == nil {
			int_app, _ := strconv.Atoi(i_app)
			apps = append(apps, uint32(int_app))
		}
	}
	lon_float, err := strconv.ParseFloat(lon, 32)
	if err != nil {
		log.Fatal(lon)
		panic("can't convert lon")
	}
	lat_float, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		log.Fatal(lon)
		panic("can't convert lon")
	}

	var result AppsInstalled
	result.dev_type = dev_type
	result.dev_id = dev_id
	result.lon = float64(lon_float)
	result.lat = float64(lat_float)
	result.apps = apps
	app_install = &result
	return app_install
}

func main() {
	var arg_str arguments

	arg_str.l = flag.Bool("l", false, "log")
	arg_str.log = flag.Bool("log", false, "log")
	arg_str.dry = flag.Bool("dry", false, "dry")
	arg_str.pattern = flag.String("pattern", "/data/appsinstalled/*.tsv.gz", "pattern")
	arg_str.idfa = flag.String("idfa", "127.0.0.1:33013", "idfa")
	arg_str.gaid = flag.String("gaid", "127.0.0.1:33014", "gaid")
	arg_str.adid = flag.String("adid", "127.0.0.1:33015", "adid")
	arg_str.dvid = flag.String("dvid", "127.0.0.1:33016", "dvid")
	flag.Parse()
	log.Println(fmt.Sprintln("Memc loader started with options: " + convert_pointer_struct_to_json(arg_str)))
	main_exec(arg_str)
}
