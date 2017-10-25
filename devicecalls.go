package GoSDK

import (
	"errors"
	"fmt"
)

const (
	_DEVICE_HEADER_KEY     = "ClearBlade-DeviceToken"
	_DEVICES_DEV_PREAMBLE  = "/admin/devices/"
	_DEVICES_USER_PREAMBLE = "/api/v/2/devices/"
)

func (d *DevClient) GetDevices(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(d, _DEVICES_DEV_PREAMBLE+systemKey, qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (u *UserClient) GetDevices(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(u, _DEVICES_USER_PREAMBLE+systemKey, qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (u *DeviceClient) GetDevices(systemKey string) ([]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _DEVICES_USER_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetDevice(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DEVICES_DEV_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DeviceClient) GetDevice(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) GetDevice(systemKey, name string) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateDevice(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, _DEVICES_DEV_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) CreateDevice(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(u, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DeviceClient) CreateDevice(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DeviceClient) AuthenticateDeviceWithKey(systemKey, name, activeKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	postBody := map[string]interface{}{
		"deviceName": name,
		"activeKey":  activeKey,
	}
	resp, err := post(d, _DEVICES_USER_PREAMBLE+systemKey+"/auth", postBody, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	theJewels, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Got unexpected return value from AuthenticateDeviceWithKey: %+v", theJewels)
	}
	d.DeviceToken = theJewels["deviceToken"].(string)
	return theJewels, nil
}

func (d *DevClient) DeleteDevice(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _DEVICES_DEV_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DeviceClient) DeleteDevice(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (u *UserClient) DeleteDevice(systemKey, name string) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(u, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateDevice(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _DEVICES_DEV_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) UpdateDevice(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(u, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *DeviceClient) UpdateDevice(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(u, _DEVICES_USER_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//  This stuff is developer only -- key sets for devices

func (d *DevClient) GetKeyset(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DEVICES_DEV_PREAMBLE+"keys/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GenerateKeyset(systemKey, name string, count int) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{"count": count}
	resp, err := post(d, _DEVICES_DEV_PREAMBLE+"keys/"+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) RotateKeyset(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	resp, err := put(d, _DEVICES_DEV_PREAMBLE+"keys/"+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteKeyset(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _DEVICES_DEV_PREAMBLE+"keys/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetDeviceColumns(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, _DEVICES_DEV_PREAMBLE+systemKey+"/columns", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting device columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting device columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) CreateDeviceColumn(systemKey, columnName, columnType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"column_name": columnName,
		"type":        columnType,
	}
	resp, err := post(d, _DEVICES_DEV_PREAMBLE+systemKey+"/columns", data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating device column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating device column: %v", resp.Body)
	}

	return nil
}

func (d *DevClient) DeleteDeviceColumn(systemKey, columnName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]string{"column_name": columnName}

	resp, err := delete(d, _DEVICES_DEV_PREAMBLE+systemKey+"/columns", data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting device column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting device column: %v", resp.Body)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (dvc *DeviceClient) credentials() ([][]string, error) {
	ret := make([][]string, 0)
	if dvc.DeviceToken != "" {
		ret = append(ret, []string{
			_DEVICE_HEADER_KEY,
			dvc.DeviceToken,
		})
	}
	if dvc.SystemKey != "" && dvc.SystemSecret != "" {
		ret = append(ret, []string{
			_HEADER_SECRET_KEY,
			dvc.SystemSecret,
		})
		ret = append(ret, []string{
			_HEADER_KEY_KEY,
			dvc.SystemKey,
		})

	}

	if len(ret) == 0 {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or DeviceToken found")
	} else {
		return ret, nil
	}
}

//  User (user/device) calls (type Client)
//func (d *DeviceClient) AuthenticateDeviceWithKey(systemKey, name, activeKey string) (map[string]interface{}, error) {

// "Login and logout"
func (dvc *DeviceClient) Authenticate() error {
	_, err := dvc.AuthenticateDeviceWithKey(dvc.SystemKey, dvc.DeviceName, dvc.ActiveKey)
	return err
}

func (dvc *DeviceClient) Logout() error {
	return nil
}

// Device MQTT calls are mqtt.go

//  Developer cbClient calls

func (dvc *DeviceClient) preamble() string {
	return _DEVICES_USER_PREAMBLE
}

func (dvc *DeviceClient) setToken(tok string) {
	dvc.DeviceName = tok
}

func (dvc *DeviceClient) getToken() string {
	return dvc.DeviceToken
}

func (dvc *DeviceClient) getSystemInfo() (string, string) {
	return dvc.SystemKey, dvc.SystemSecret
}

func (dvc *DeviceClient) getMessageId() uint16 {
	return uint16(dvc.mrand.Int())
}

func (dvc *DeviceClient) getHttpAddr() string {
	return dvc.HttpAddr
}

func (dvc *DeviceClient) getMqttAddr() string {
	return dvc.MqttAddr
}
