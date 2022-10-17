package setting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var settingData map[string]string

func init() {
	if settingData == nil {
		var err error
		settingData, err = readSettingFile("../config/config.json")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}
func ReadFieldsFromSettingData(fields_name []string) map[string]string {
	var fields_value = make(map[string]string, len(fields_name))
	for _, fn := range fields_name {
		value, exists := settingData[fn]
		if !exists {
			panic(fmt.Sprintf("%s field name isn't exist in the setting data", fn))
		}
		fields_value[fn] = value
	}
	return fields_value
}
func ReadFieldInSettingData(field_name string) string {
	value, exists := settingData[field_name]
	if !exists {
		panic(fmt.Sprintf("%s field name isn't exist in the setting data", field_name))
	}
	return value
}

func readSettingFile(setting_path string) (map[string]string, error) {
	file_bytes, err := ioutil.ReadFile(setting_path)
	if err != nil {
		return nil, fmt.Errorf("error when opening setting file")
	}
	var data map[string]string
	err = json.Unmarshal(file_bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error occurred during unmarshal setting file")
	}
	err = validateSettingData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func validateSettingData(data map[string]string) error {
	var expect_fields_name = []string{
		"ADMIN_TELEGRAM_ID",
		"ADMIN_TELEGRAM_USERNAME",
		"DSN",
	}

	var exists bool
	var data_value string
	for _, data_name := range expect_fields_name {
		data_value, exists = data[data_name]
		if !exists {
			return fmt.Errorf("%s doesn't exists in setting file", data_name)
		}
		if data_value == "" {
			return fmt.Errorf("%s is empty in setting file", data_name)
		}

	}
	return nil
}
