package zont

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseUrl = "https://zont-online.ru/api/"
)

type Client struct {
	httpClient        *http.Client
	httpTimeout       time.Duration
	debug             bool
	AuthTokenResponse *AuthTokenResponse
	clientName        string
	xZontClient       string
	login             string
	password          string
}

type AuthTokenResponse struct {
	Token             string        `json:"token"`
	IsGuest           bool          `json:"is_guest"`
	Permissions       []interface{} `json:"permissions"`
	IsStaff           bool          `json:"is_staff"`
	IsSuperuser       bool          `json:"is_superuser"`
	AllowTerminal     bool          `json:"allow_terminal"`
	ID                int           `json:"id"`
	Username          string        `json:"username"`
	Fullname          string        `json:"fullname"`
	Phone             interface{}   `json:"phone"`
	Email             string        `json:"email"`
	EmailConfirmed    bool          `json:"email_confirmed"`
	IsDemo            bool          `json:"is_demo"`
	GuestPasswordSet  bool          `json:"guest_password_set"`
	ShowNetworkEvents bool          `json:"show_network_events"`
	AutoscanSystem    interface{}   `json:"autoscan_system"`
	PublicAPIAllowed  bool          `json:"public_api_allowed"`
	ZthDashboard      struct {
		Show              bool `json:"show"`
		EnableModes       bool `json:"enable_modes"`
		EnableZones       bool `json:"enable_zones"`
		EnableTemperature bool `json:"enable_temperature"`
		EnableAlert       bool `json:"enable_alert"`
	} `json:"zth_dashboard"`
	Sims []struct {
		SimCard struct {
			ID struct {
				Operator string `json:"operator"`
				ID       string `json:"id"`
			} `json:"id"`
			Iccid  string `json:"iccid"`
			Msisdn string `json:"msisdn"`
		} `json:"sim_card"`
		Active        bool        `json:"active"`
		UserID        int         `json:"user_id"`
		PaidUntil     int         `json:"paid_until"`
		AutoProlong   bool        `json:"auto_prolong"`
		Terminated    bool        `json:"terminated"`
		FirstActivate interface{} `json:"first_activate"`
		Limit         string      `json:"limit"`
		Tariff        string      `json:"tariff"`
	} `json:"sims"`
	ForeignSims         []interface{} `json:"foreign_sims"`
	DateLastNewsReading int           `json:"date_last_news_reading"`
	WebuiBeta           bool          `json:"webui_beta"`
	AllowNewUI          bool          `json:"allow_new_ui"`
	UseNewUI            bool          `json:"use_new_ui"`
	Ok                  bool          `json:"ok"`
}

type DevicesResponse struct {
	Devices    []Device `json:"devices"`
	DeviceTree []struct {
		ZontID int `json:"ZontId"`
	} `json:"device_tree"`
	Ok bool `json:"ok"`
}

type Device struct {
	Access       []interface{} `json:"access,omitempty"`
	Capabilities []string      `json:"capabilities,omitempty"`
	DeviceType   struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"device_type,omitempty"`
	WidgetType   string `json:"widget_type,omitempty"`
	HardwareType struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"hardware_type,omitempty"`
	ID                       int         `json:"id,omitempty"`
	IP                       string      `json:"ip,omitempty"`
	IsActive                 bool        `json:"is_active,omitempty"`
	Online                   bool        `json:"online,omitempty"`
	OwnerUsername            string      `json:"owner_username,omitempty"`
	UserID                   int         `json:"user_id,omitempty"`
	LastReceiveTime          int         `json:"last_receive_time,omitempty"`
	LastReceiveTimeRelative  int         `json:"last_receive_time_relative,omitempty"`
	Name                     string      `json:"name,omitempty"`
	Color                    string      `json:"color,omitempty"`
	Notes                    string      `json:"notes,omitempty"`
	Serial                   string      `json:"serial,omitempty"`
	VisibleDeviceType        interface{} `json:"visible_device_type,omitempty"`
	FirmwareVersion          []int       `json:"firmware_version,omitempty"`
	IsConfiguredByWizard     interface{} `json:"is_configured_by_wizard,omitempty"`
	IsConfiguredByZontWizard interface{} `json:"is_configured_by_zont_wizard,omitempty"`
	SuggestFirmwareUpgrade   bool        `json:"suggest_firmware_upgrade,omitempty"`
	GraphsConfig             struct {
		Blocks []struct {
			Sources []struct {
				Class  string `json:"class"`
				Params struct {
					ZoneNo int `json:"zone_no"`
				} `json:"params,omitempty"`
				Params0 struct {
					ThermID string `json:"therm_id"`
				} `json:"params,omitempty"`
			} `json:"sources"`
			Height       interface{} `json:"height"`
			HeightMobile string      `json:"height_mobile"`
		} `json:"blocks"`
	} `json:"graphs_config,omitempty"`
	CamsShow            bool `json:"cams_show,omitempty"`
	ShowHeatingTab      bool `json:"show_heating_tab,omitempty"`
	ServerNotifications struct {
		Events struct {
			Enabled bool `json:"enabled"`
		} `json:"events"`
		Offline struct {
			Enabled bool `json:"enabled"`
			Timeout int  `json:"timeout"`
		} `json:"offline"`
	} `json:"server_notifications,omitempty"`
	DebugTextMessagesRegexp interface{} `json:"debug_text_messages_regexp,omitempty"`
	StationaryLocation      struct {
		Loc []float64 `json:"loc"`
	} `json:"stationary_location,omitempty"`
	Cams                 []interface{} `json:"cams,omitempty"`
	SpecialistInfo       interface{}   `json:"specialist_info,omitempty"`
	InstallationState    interface{}   `json:"installation_state,omitempty"`
	Maintenances         interface{}   `json:"maintenances,omitempty"`
	WorkState            interface{}   `json:"work_state,omitempty"`
	IndividualDeviceName string        `json:"individual_device_name,omitempty"`
	IndividualDeviceNote string        `json:"individual_device_note,omitempty"`
	UISettings           struct {
		HeatTab struct {
		} `json:"heat_tab"`
		Icons struct {
		} `json:"icons"`
	} `json:"ui_settings,omitempty"`
	Timezone                     int         `json:"timezone,omitempty"`
	LastGuardEvent               interface{} `json:"last_guard_event,omitempty"`
	ThermostatErrorInputPolarity string      `json:"thermostat_error_input_polarity,omitempty"`
	ThermostatInputconfig        struct {
		Num1 string `json:"1"`
		Num2 string `json:"2"`
	} `json:"thermostat_inputconfig,omitempty"`
	ThermostatEnableGuard bool `json:"thermostat_enable_guard,omitempty"`
	BoilerInfo            struct {
		Model  string `json:"model"`
		Vendor string `json:"vendor"`
	} `json:"boiler_info,omitempty"`
	ThermostatExtMode int    `json:"thermostat_ext_mode,omitempty"`
	ThermostatMode    string `json:"thermostat_mode,omitempty"`
	ThermostatGate    bool   `json:"thermostat_gate,omitempty"`
	Tempschedule      struct {
		Day  []float64 `json:"day"`
		Week struct {
			Num0 []float64 `json:"0"`
			Num1 []float64 `json:"1"`
			Num2 []float64 `json:"2"`
			Num3 []float64 `json:"3"`
			Num4 []float64 `json:"4"`
			Num5 []float64 `json:"5"`
			Num6 []float64 `json:"6"`
		} `json:"week"`
	} `json:"tempschedule,omitempty"`
	Tempstep      int `json:"tempstep,omitempty"`
	Notifications struct {
		Alarm struct {
			Numbers    string `json:"numbers"`
			PowerOff   string `json:"power-off"`
			PowerOn    string `json:"power-on"`
			Blackout   string `json:"blackout"`
			Doors      string `json:"doors"`
			DriverCall string `json:"driver_call"`
			Ignition   string `json:"ignition"`
			Moving     string `json:"moving"`
			Shock      string `json:"shock"`
			Tilt       string `json:"tilt"`
			TrunkHood  string `json:"trunk-hood"`
		} `json:"alarm"`
		Guard struct {
			Numbers string `json:"numbers"`
			Off     string `json:"off"`
			On      string `json:"on"`
		} `json:"guard"`
		Info struct {
			Balance    string `json:"balance"`
			EcuError   string `json:"ecu_error"`
			Numbers    string `json:"numbers"`
			FobBattery string `json:"fob_battery"`
		} `json:"info"`
		Thermostat struct {
			BoilerFail string `json:"boiler_fail"`
			TempHigh   string `json:"temp_high"`
			TempLow    string `json:"temp_low"`
			ThermMalf  string `json:"therm_malf"`
		} `json:"thermostat"`
		Autoignition struct {
			Breakdown string `json:"breakdown"`
			Fail      string `json:"fail"`
			Numbers   string `json:"numbers"`
			Success   string `json:"success"`
		} `json:"autoignition"`
	} `json:"notifications,omitempty"`
	ThermostatHysteresis  float64 `json:"thermostat_hysteresis,omitempty"`
	ThermostatTempsLimits struct {
		Num0 struct {
			Max interface{} `json:"max"`
			Min interface{} `json:"min"`
		} `json:"0"`
		Num1 struct {
			Max interface{} `json:"max"`
			Min interface{} `json:"min"`
		} `json:"1"`
	} `json:"thermostat_temps_limits,omitempty"`
	TemperatureAlarm struct {
		High int `json:"high"`
		Low  int `json:"low"`
	} `json:"temperature_alarm,omitempty"`
	SimInDevice struct {
		SimType string `json:"sim_type"`
		SimID   struct {
			Operator string `json:"operator"`
			ID       string `json:"id"`
		} `json:"sim_id"`
		ForeignMsisdn interface{} `json:"foreign_msisdn"`
	} `json:"sim_in_device,omitempty"`
	Balance struct {
		Limit   int    `json:"limit"`
		Ussd    string `json:"ussd"`
		Warning bool   `json:"warning"`
	} `json:"balance,omitempty"`
	TrustedPhones    string      `json:"trusted_phones,omitempty"`
	GsmRoaming       interface{} `json:"gsm_roaming,omitempty"`
	Imei             string      `json:"imei,omitempty"`
	Iccid            interface{} `json:"iccid,omitempty"`
	OtEnabled        bool        `json:"ot_enabled,omitempty"`
	OtSaveParams     []string    `json:"ot_save_params,omitempty"`
	OtMinSetpoint    float64     `json:"ot_min_setpoint,omitempty"`
	OtMaxSetpoint    float64     `json:"ot_max_setpoint,omitempty"`
	OtMaxMl          float64     `json:"ot_max_ml,omitempty"`
	OtDhwSetpoint    float64     `json:"ot_dhw_setpoint,omitempty"`
	OtMinWp          float64     `json:"ot_min_wp,omitempty"`
	OtConfig         []string    `json:"ot_config,omitempty"`
	OtMode           string      `json:"ot_mode,omitempty"`
	OtBoilerType     string      `json:"ot_boiler_type,omitempty"`
	OtShowDhwControl bool        `json:"ot_show_dhw_control,omitempty"`
	RfStatus         interface{} `json:"rf_status,omitempty"`
	BoilerDelay      struct {
		Off int `json:"off"`
		On  int `json:"on"`
	} `json:"boiler_delay,omitempty"`
	Pza struct {
		Enabled bool `json:"enabled"`
		Curve   int  `json:"curve"`
	} `json:"pza,omitempty"`
	PzaMaxDelta struct {
		Enabled bool `json:"enabled"`
	} `json:"pza_max_delta,omitempty"`
	ThermostatExtModesConfig struct {
		Num0 struct {
			Active         bool        `json:"active"`
			Name           string      `json:"name"`
			ScheduleNumber interface{} `json:"schedule_number"`
			ZoneSensors    struct {
				Num1 interface{} `json:"1"`
			} `json:"zone_sensors"`
			ZoneTemp struct {
				Num1 float64 `json:"1"`
			} `json:"zone_temp"`
		} `json:"0"`
		Num1 struct {
			Active         bool        `json:"active"`
			Name           string      `json:"name"`
			ScheduleNumber interface{} `json:"schedule_number"`
			ZoneSensors    struct {
				Num1 interface{} `json:"1"`
			} `json:"zone_sensors"`
			ZoneTemp struct {
				Num1 float64 `json:"1"`
			} `json:"zone_temp"`
		} `json:"1"`
		Num2 struct {
			Active         bool   `json:"active"`
			Name           string `json:"name"`
			ScheduleNumber int    `json:"schedule_number"`
			ZoneSensors    struct {
			} `json:"zone_sensors"`
			ZoneTemp struct {
			} `json:"zone_temp"`
		} `json:"2"`
		Num3 struct {
			Active         bool        `json:"active"`
			Name           string      `json:"name"`
			ScheduleNumber interface{} `json:"schedule_number"`
			ZoneSensors    struct {
				Num1 interface{} `json:"1"`
			} `json:"zone_sensors"`
			ZoneTemp struct {
				Num1 float64 `json:"1"`
			} `json:"zone_temp"`
		} `json:"3"`
		Num4 struct {
			Active   bool   `json:"active"`
			Name     string `json:"name"`
			ZoneTemp struct {
				Num1 float64 `json:"1"`
			} `json:"zone_temp"`
			ScheduleNumber interface{} `json:"schedule_number"`
			ZoneSensors    struct {
				Num1 int `json:"1"`
			} `json:"zone_sensors"`
		} `json:"4"`
		Num5 struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		} `json:"5"`
		Num6 struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		} `json:"6"`
		Num7 struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		} `json:"7"`
		Num8 struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		} `json:"8"`
		Num9 struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		} `json:"9"`
	} `json:"thermostat_ext_modes_config,omitempty"`
	ThermostatTargetTemps      *map[string]ThermostatTargetTemps `json:"thermostat_target_temps,omitempty"`
	ThermostatExtModesAdvanced bool                              `json:"thermostat_ext_modes_advanced,omitempty"`
	ThermostatRelayMode        string                            `json:"thermostat_relay_mode,omitempty"`
	OtGateEnabled              bool                              `json:"ot_gate_enabled,omitempty"`
	UseInternetWeatherForPza   bool                              `json:"use_internet_weather_for_pza,omitempty"`
	Thermometers               []struct {
		IsAssignedToSlot bool   `json:"is_assigned_to_slot"`
		Slot             int    `json:"slot"`
		UUID             string `json:"uuid"`
		Serial           string `json:"serial"`
		Type             string `json:"type"`
		Name             string `json:"name"`
		Color            string `json:"color"`
		Limits           struct {
			Low  int `json:"low"`
			High int `json:"high"`
		} `json:"limits"`
		Function  string `json:"function"`
		Functions []struct {
			F    string `json:"f"`
			Zone int    `json:"zone"`
		} `json:"functions"`
		Sort          int     `json:"sort"`
		LastState     string  `json:"last_state"`
		LastValue     float64 `json:"last_value"`
		LastValueTime int     `json:"last_value_time"`
	} `json:"thermometers,omitempty"`
	Filetransfers   []interface{} `json:"filetransfers,omitempty"`
	InternetWeather float64       `json:"internet_weather,omitempty"`
	AspBilling      struct {
		InService      bool `json:"in_service"`
		AllowedForUser bool `json:"allowed_for_user"`
	} `json:"asp_billing,omitempty"`
}

type ThermostatTargetTemps struct {
	Manual bool    `json:"manual"`
	Temp   float64 `json:"temp"`
}

type LoadDataResponse struct {
	Ok        bool `json:"ok,omitempty"`
	Responses []struct {
		DeviceID      int  `json:"device_id,omitempty"`
		Ok            bool `json:"ok,omitempty"`
		TimeTruncated bool `json:"time_truncated,omitempty"`
		Temperature   map[string]struct {
			Name        string      `json:"name,omitempty"`
			Color       string      `json:"color,omitempty"`
			Sort        int         `json:"sort,omitempty"`
			Temperature [][]float64 `json:"temperature,omitempty"`
		} `json:"temperature,omitempty"`
		Timings struct {
			Temperature struct {
				Wall float64 `json:"wall,omitempty"`
				Proc float64 `json:"proc,omitempty"`
			} `json:"temperature,omitempty"`
		} `json:"timings,omitempty"`
	} `json:"responses,omitempty"`
}

// https://zont-online.ru/api/docs/#thermostat_work
type LoadDataThermostatWorkResponse struct {
	Ok        bool `json:"ok"`
	Responses []struct {
		DeviceID       int  `json:"device_id"`
		Ok             bool `json:"ok"`
		TimeTruncated  bool `json:"time_truncated"`
		ThermostatWork struct {
			ThermostatMode [][]int     `json:"thermostat_mode"`
			DhwT           [][]float64 `json:"dhw_t"`
			Power          [][]any     `json:"power"`
			Fail           [][]any     `json:"fail"`
			Gate           [][]any     `json:"gate"`
			Ot             struct {
				Cs  [][]int     `json:"cs"`
				Bt  [][]int     `json:"bt"`
				Dt  [][]int     `json:"dt"`
				Rwt [][]int     `json:"rwt"`
				Rml [][]int     `json:"rml"`
				Wp  [][]float64 `json:"wp"`
				S   [][]any     `json:"s"`
			} `json:"ot"`
			Zones struct {
				Num1 struct {
					TargetTemp [][]int `json:"target_temp"`
					Worktime   [][]int `json:"worktime"`
				} `json:"1"`
			} `json:"zones"`
			BoilerWorkTime [][]int `json:"boiler_work_time"`
			TargetTemp     [][]int `json:"target_temp"`
		} `json:"thermostat_work"`
		Timings struct {
			ThermostatWork struct {
				Wall float64 `json:"wall"`
				Proc float64 `json:"proc"`
			} `json:"thermostat_work"`
		} `json:"timings"`
	} `json:"responses"`
}

// NewClient return new client
func NewClient(clientName, xZontClient, login, password string) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	standardClient := retryClient.StandardClient() // *http.Client

	httpTimeout := 20 * time.Second

	return &Client{
		standardClient,
		httpTimeout,
		false,
		nil,
		clientName,
		xZontClient,
		login,
		password,
	}
}

// PostRequestHandler handle all post requests
func (cl *Client) PostRequestHandler(data interface{}, uri string, basic bool) []byte {
	jsonValue, err := json.Marshal(data)
	if err != nil {
		ContextLogger.Error(err)
	}

	b := bytes.NewBuffer(jsonValue)
	req, err := http.NewRequest("POST", uri, b)
	if err != nil {
		ContextLogger.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ZONT-Client", cl.xZontClient)
	if basic {
		req.SetBasicAuth(cl.login, cl.password)
	} else {
		req.Header.Set("X-ZONT-Token", cl.AuthTokenResponse.Token)
	}
	res, err := cl.httpClient.Do(req)
	if err != nil {
		ContextLogger.Error(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		ContextLogger.Error(err)
	}

	defer res.Body.Close()

	return body
}

// GetAuthToken return token base on login and password
func (cl *Client) GetAuthToken() {
	method := "get_authtoken"
	uri := baseUrl + method
	values := map[string]string{"client_name": cl.clientName}

	body := cl.PostRequestHandler(values, uri, true)

	err := json.Unmarshal(body, &cl.AuthTokenResponse)
	if err != nil {
		ContextLogger.Error(err)
	}
}

// GetDevices return device info in DevicesResponse format
func (cl *Client) GetDevices() *DevicesResponse {
	if len(cl.AuthTokenResponse.Token) < 1 {
		ContextLogger.Infoln("AuthToken not exist!")
		return nil
	}

	method := "devices"
	uri := baseUrl + method
	values := map[string]string{"client_name": cl.clientName}

	body := cl.PostRequestHandler(values, uri, false)

	devices := DevicesResponse{}

	err := json.Unmarshal(body, &devices)
	if err != nil {
		ContextLogger.Error(err)
	}

	return &devices
}

// UpdateDevice send update request based on data interface
func (cl *Client) UpdateDevice(data interface{}) error {
	if len(cl.AuthTokenResponse.Token) < 1 {
		ContextLogger.Infoln("AuthToken not exist!")
		return nil
	}

	method := "update_device"
	uri := baseUrl + method

	body := cl.PostRequestHandler(data, uri, false)

	devices := DevicesResponse{}

	err := json.Unmarshal(body, &devices)
	if err != nil {
		ContextLogger.Error(err)
		return err
	}

	return nil
}

// LoadData return device information and metrics based on data interface
func (cl *Client) LoadData(data interface{}) *LoadDataResponse {
	if len(cl.AuthTokenResponse.Token) < 1 {
		ContextLogger.Infoln("AuthToken not exist!")
		return nil
	}

	method := "load_data"
	uri := baseUrl + method

	body := cl.PostRequestHandler(data, uri, false)

	loadData := LoadDataResponse{}

	err := json.Unmarshal(body, &loadData)
	if err != nil {
		ContextLogger.Error(err)
	}

	return &loadData
}

// LoadDataThermostatWork return device information and metrics based on data interface
func (cl *Client) LoadDataThermostatWork(data interface{}) *LoadDataThermostatWorkResponse {
	if len(cl.AuthTokenResponse.Token) < 1 {
		ContextLogger.Infoln("AuthToken not exist!")
		return nil
	}

	method := "load_data"
	uri := baseUrl + method

	body := cl.PostRequestHandler(data, uri, false)

	loadData := LoadDataThermostatWorkResponse{}

	err := json.Unmarshal(body, &loadData)
	if err != nil {
		ContextLogger.Error(err)
	}

	return &loadData
}

// GetCurrentTemp return current temperature from first thermometer on device with deviceId
func (cl *Client) GetCurrentTemp(deviceId int) (temp float64) {
	dataRequest := struct {
		DeviceID  int      `json:"device_id"`
		DataTypes []string `json:"data_types"`
		MinTime   int64    `json:"mintime"`
		MaxTime   int64    `json:"maxtime"`
	}{}

	dataRequest.DeviceID = deviceId
	dataRequest.DataTypes = append(dataRequest.DataTypes, "temperature")
	dataRequest.MaxTime = time.Now().Unix()
	dataRequest.MinTime = time.Now().Add(-180 * time.Second).Unix()

	dataLoad := struct {
		Requests []struct {
			DeviceID  int      `json:"device_id"`
			DataTypes []string `json:"data_types"`
			MinTime   int64    `json:"mintime"`
			MaxTime   int64    `json:"maxtime"`
		} `json:"requests"`
	}{}

	dataLoad.Requests = append(dataLoad.Requests, dataRequest)

	loadResp := &LoadDataResponse{}
	loadResp = cl.LoadData(dataLoad)
	for k := range loadResp.Responses[0].Temperature {
		return loadResp.Responses[0].Temperature[k].Temperature[0][1]
	}
	return temp
}

// GetCurrentHotWaterTemp return current hot water temperature
func (cl *Client) GetCurrentHotWaterTemp(deviceId int) (temp float64) {
	dataRequest := struct {
		DeviceID  int      `json:"device_id"`
		DataTypes []string `json:"data_types"`
		MinTime   int64    `json:"mintime"`
		MaxTime   int64    `json:"maxtime"`
	}{}

	dataRequest.DeviceID = deviceId
	dataRequest.DataTypes = append(dataRequest.DataTypes, "thermostat_work")
	dataRequest.MaxTime = time.Now().Unix()
	dataRequest.MinTime = time.Now().Add(-180 * time.Second).Unix()

	dataLoad := struct {
		Requests []struct {
			DeviceID  int      `json:"device_id"`
			DataTypes []string `json:"data_types"`
			MinTime   int64    `json:"mintime"`
			MaxTime   int64    `json:"maxtime"`
		} `json:"requests"`
	}{}

	dataLoad.Requests = append(dataLoad.Requests, dataRequest)

	loadResp := cl.LoadDataThermostatWork(dataLoad)
	return loadResp.Responses[0].ThermostatWork.DhwT[0][1]
}

type ThermostatData struct {
	DeviceID              int                              `json:"device_id"`
	ThermostatTargetTemps map[string]ThermostatTargetTemps `json:"thermostat_target_temps"`
}

// SetTargetTemp send temperature update request based on deviceId and targetTemp
func (cl *Client) SetTargetTemp(deviceId int, termostatid string, targetTemp float64) error {
	// Update data
	data := ThermostatData{}
	data.DeviceID = deviceId
	data.ThermostatTargetTemps = map[string]ThermostatTargetTemps{
		termostatid: {Manual: true, Temp: targetTemp},
	}

	err := cl.UpdateDevice(data)
	if err != nil {
		ContextLogger.Error(err)
		return err
	}

	return nil
}
