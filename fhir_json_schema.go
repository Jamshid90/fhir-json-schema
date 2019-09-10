package fhir_json_schema

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetJsonMap(data []byte) (map[string]json.RawMessage, error) {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func RemoveLowerDash(input_file, output_file string) {

	start := time.Now()

	f, err := os.Create(output_file)
	Check(err)

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(f)

	data, err := ioutil.ReadFile(input_file)
	Check(err)

	jsonMap, err := GetJsonMap(data)
	Check(err)

	fmt.Fprintln(w, `{`)

	i_ := 1
	for k, v := range jsonMap {

		point_ := ","
		if len(jsonMap) == i_ {
			point_ = fmt.Sprintf("")
		}

		if k == "definitions" {

			fmt.Fprintf(w, ` "%s" : { `, k)

			definitions, _ := GetJsonMap(v)

			i := 1
			for _k, _v := range definitions {

				point := ","
				if len(definitions) == i {
					point = fmt.Sprintf("")
				}

				resource, _ := GetJsonMap(_v)

				if _, ok := resource["properties"]; ok == false {
					fmt.Fprintf(w, ` "%s" : %s `+point, _k, _v)
					i++
					continue
				}

				fmt.Fprintf(w, ` "%s" : { `, _k)

				_i := 1
				for r_k, r_v := range resource {

					_point := ","
					if _i == len(resource) {
						_point = fmt.Sprintf("")
					}

					if r_k == "description" {
						fmt.Fprintf(w, ` "description" : %s `+_point, string(r_v))
						_i++
						continue
					}

					if r_k == "properties" {

						fmt.Fprintf(w, ` "%s" : { `, r_k)

						properties, _ := GetJsonMap(r_v)

						__i := 1
						for p_k, p_v := range properties {

							if string(p_k[0]) == "_" {
								continue
							}

							__point := ""
							if __i > 1 {
								__point = fmt.Sprintf(",")
							}

							fmt.Fprintf(w, __point+` "%s" : { `, p_k)

							propertr_value, _ := GetJsonMap(p_v)

							___i := 1
							for pv_k, pv_v := range propertr_value {

								___point := ","
								if ___i == len(propertr_value) {
									___point = fmt.Sprintf("")
								}

								var j []byte
								j, err = json.Marshal(pv_v)

								fmt.Fprintf(w, ` "%s" : %s `+___point, pv_k, string(j))

								___i++
							}

							fmt.Fprintln(w, ` } `)

							__i++
						}

						fmt.Fprintln(w, ` } `+_point)
						_i++
						continue

					}

					fmt.Fprintf(w, ` "%s" : %s `+_point, r_k, r_v)

					_i++
				}

				fmt.Fprintln(w, ` } `+point)
				i++
			}

			fmt.Fprintln(w, ` } `+point_)
			i_++
			continue
		}

		fmt.Fprintf(w, ` "%s" : %s `+point_, k, v)

		i_++
	}

	fmt.Fprintln(w, ` } `)

	err = w.Flush()
	Check(err)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf(" -------- Successful executed %s -------- ", elapsed)
}

func GenerateStructMap(input_file, output_file, file_package, map_name string) {
	start := time.Now()

	f, err := os.Create(output_file)
	Check(err)

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(f)

	data, err := ioutil.ReadFile(input_file)
	Check(err)

	jsonMap, err := GetJsonMap(data)
	Check(err)

	fmt.Fprintf(w, `
package %s

import (
	"errors"
)
`, file_package)

	fmt.Fprintln(w, `
func GetFhirResourceMap() map[string]interface{}{
	return map[string]interface{}{`)
	for k, v := range jsonMap {
		if k == "definitions" {
			definitions, _ := GetJsonMap(v)
			for _k, _v := range definitions {
				resource, _ := GetJsonMap(_v)
				if _, ok := resource["properties"]; ok == true {
					properties, _ := GetJsonMap(resource["properties"])
					if _, ok := properties["resourceType"]; ok == true {
						_k := strings.Replace(_k, "_", "", -1)
						fmt.Fprintf(w, `		"%s" : &%s{}, `+"\n", _k, _k)
					}
				}
			}
		}
	}
	fmt.Fprintln(w, `		}
}`)

	fmt.Fprintf(w, `
func GetFhirResourceInstance(resource_name string) (interface{}, error) {
	%s := GetFhirResourceMap()
	resource, ok := %s[resource_name]
	if ok == true {
		return resource, nil
	}
	return resource, errors.New("Resource not found")
}`, map_name, map_name)

	err = w.Flush()
	Check(err)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf(" -------- Successful executed %s -------- ", elapsed)
}
