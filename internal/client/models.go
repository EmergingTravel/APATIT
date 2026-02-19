package client

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// --- Raw API Structures (direct JSON parsing) ---

// EntryRaw
// is a 'task_graph_stat' API response.
// It contains TM (tochka monitoringa) info with TmID; TmName and TmRes (results with time, speed).
// P.S. TM  will be used as MP (Monitoring Point) after processing.
type EntryRaw struct {
	TmID     string      `json:"tm_id"`
	TmName   string      `json:"tm_name"`
	TmNameEn string      `json:"tm_name_en"`
	TmRes    []*TmResRaw `json:"tm_res"`
}

// TmResRaw
// is a 'tm_res' JSON from 'task_graph_stat' API response.
// It contains TM (tochka monitoringa) info with speed and times for task.
// P.S. TM (tochka monitoringa) will be used as MP (Monitoring Point) after processing.
type TmResRaw struct {
	Connect *string `json:"connect"`
	DNS     *string `json:"dns"`
	Server  *string `json:"server"`
	TmStamp *string `json:"tmstamp"`
	Speed   *string `json:"speed"`
	Total   *string `json:"total"`
}

// TaskRaw
// is a result of 'tasks' API request.
// It contains info about specified task.
type TaskRaw struct {
	// Console Task Status (enabled/disabled)
	Status int `json:"status"`
	// TaskID
	ID int `json:"tid"`
	// Service Name (in task)
	SName string `json:"nazv"`
	// IP / DNS-name of service
	Address string `json:"name"`
	// Task Status (1 — works; 0 — doesn't work)
	TaskStatus int `json:"log_status"`
	// Blacklist status (if address was found in RKN, Spamhause, etc. blacklists)
	BlackListStatus int `json:"rk_log_status"`
	// Virus status
	VirusStatus int `json:"sb_log_status"`
	// This extra data will be obtained via API and may be used in the future.
	// Datetime of last check
	LastData string `json:"last_data"`
	// Datetime of last service status change
	LogData string `json:"log_data"`
	// Checking period (default)
	Period int `json:"period"`
	// Checking period (during error status)
	PeriodError int `json:"period_error"`
	// Blacklist status check settings
	Rk        int         `json:"rk"`
	RkIp      int         `json:"rk_ip"`
	RkLogData interface{} `json:"rk_log_data"`
	Rrd       int         `json:"rrd"`
	// Virus status check settings
	Sb        int         `json:"sb"`
	SbLogData interface{} `json:"sb_log_data"`
	// Contacts
	TasksImsIcqList      []interface{} `json:"tasks_ims_icq_list"`
	TasksImsJabberList   []interface{} `json:"tasks_ims_jabber_list"`
	TasksImsSkypeList    []interface{} `json:"tasks_ims_skype_list"`
	TasksImsTelegramList []interface{} `json:"tasks_ims_telegram_list"`
	// Check type
	Tip string `json:"tip"`
	// Unavailability time
	UptimeNw int `json:"uptime_nw"`
	// Availability time
	UptimeW int `json:"uptime_w"`
	// ???
	Uveddva int `json:"uveddva"`
}

// MonitoringPointRaw
// is a result of 'tm' API request.
// It contains information about specified Monitoring Point.
type MonitoringPointRaw struct {
	// Monitoring Point ID
	ID string `json:"id"`
	// Monitoring Point Name
	Name string `json:"name"`
	// Monitoring Point NameEn
	NameEn string `json:"name_en"`
	// Monitoring Point IP
	IP string `json:"ip"`
	// Monitoring Point GPS
	GPS string `json:"gps"`
	// Monitoring point availability status
	Status string `json:"status"`
}

// TaskStatRaw
// is a result of 'task_stat' API request.
// Contains info about last task events.
type TaskStatRaw struct {
	TasksLogs []*TasksLogsRaw `json:"tasks_logs"`
	Uptime    string          `json:"uptime"`
	UptimeNw  int             `json:"uptime_nw"`
	UptimeW   int             `json:"uptime_w"`
}

// TasksLogsRaw
// is a 'TasksLogs' array element in 'TaskStatRaw'.
type TasksLogsRaw struct {
	Comment    *any    `json:"comment"`
	Data       *string `json:"data"`
	Descr      *string `json:"descr"`
	Status     *int    `json:"status"`
	Tm         *string `json:"tm"`
	TmEn       *string `json:"tm_en"`
	TmID       *string `json:"tm_id"`
	Traceroute *string `json:"traceroute"`
}

// --- Processed Data Structures ---

// TaskInfo
// is a processed TaskRaw
type TaskInfo struct {
	EnabledStatus   int
	ID              int
	ServiceName     string
	URL             string
	TaskStatus      int
	BlackListStatus int
	VirusStatus     int
	Timestamp       time.Time
}

// MonitoringPointInfo
// is a processed MonitoringPointRaw.
type MonitoringPointInfo struct {
	ID     string
	Name   string
	NameEn string
	IP     string
	GPS    string
	Status int64
}

// MonitoringPointEntry
// is a processed EntryRaw.
type MonitoringPointEntry struct {
	ID     string
	Name   string
	NameEn string
	Status int
	Result []*MonitoringPointConnectionResult
}

// MonitoringPointConnectionResult
// is a processed TmResRaw.
type MonitoringPointConnectionResult struct {
	Connect   float64
	DNS       float64
	Server    float64
	Timestamp int64
	Speed     int64
	Total     float64
}

// TaskStatEntry
// is a processed TaskStatRaw.
type TaskStatEntry struct {
	TaskID    string
	TaskName  string
	Timestamp time.Time
	TaskLogs  []*TaskLog
}

// TaskLog
// is a processed TasksLogsRaw.
type TaskLog struct {
	Data        string
	Description string
	Status      int64
	MPNameRu    string
	MPName      string
	MPID        string
	Traceroute  string
}

// ProcessMonitoringPointInfo
// converts TaskRaw to TaskInfo.
func (mp *MonitoringPointRaw) ProcessMonitoringPointInfo() *MonitoringPointInfo {
	return &MonitoringPointInfo{
		ID:     mp.ID,
		Name:   mp.Name,
		NameEn: mp.NameEn,
		IP:     mp.IP,
		GPS:    mp.GPS,
		Status: parseInt(&mp.Status, "status"),
	}
}

// ProcessTaskInfo
// converts TaskRaw to TaskInfo.
func (t *TaskRaw) ProcessTaskInfo() *TaskInfo {
	return &TaskInfo{
		EnabledStatus:   t.Status,
		ID:              t.ID,
		ServiceName:     t.SName,
		URL:             t.Address,
		TaskStatus:      t.TaskStatus,
		BlackListStatus: t.BlackListStatus,
		VirusStatus:     t.VirusStatus,
		Timestamp:       time.Now(),
	}
}

// ProcessMonitoringPointEntry
// converts "raw" structure EntryRaw into Entry with correct types of data.
func (e *EntryRaw) ProcessMonitoringPointEntry() *MonitoringPointEntry {
	entry := &MonitoringPointEntry{
		ID:     e.TmID,
		Name:   e.TmName,
		NameEn: e.TmNameEn,
		Result: make([]*MonitoringPointConnectionResult, 0, len(e.TmRes)),
	}

	for _, resRaw := range e.TmRes {
		MPRes := &MonitoringPointConnectionResult{
			Connect:   parseFloat(resRaw.Connect, "connect"),
			DNS:       parseFloat(resRaw.DNS, "dns"),
			Server:    parseFloat(resRaw.Server, "server"),
			Total:     parseFloat(resRaw.Total, "total"),
			Speed:     parseInt(resRaw.Speed, "speed"),
			Timestamp: parseInt(resRaw.TmStamp, "timestamp"),
		}
		entry.Result = append(entry.Result, MPRes)
	}
	return entry
}

// ProcessTaskEntry
// converts "raw" structure TaskStatRaw to TaskStatEntry
func (t *TaskStatRaw) ProcessTaskEntry() *TaskStatEntry {
	entry := &TaskStatEntry{
		TaskLogs: make([]*TaskLog, 0, len(t.TasksLogs)),
	}

	for _, resRaw := range t.TasksLogs {
		TaskStatRes := &TaskLog{
			Data:        *resRaw.Data,
			Description: *resRaw.Descr,
			Status:      int64(*resRaw.Status),
			MPName:      *resRaw.TmEn,
			MPNameRu:    *resRaw.Tm,
			MPID:        *resRaw.TmID,
			Traceroute:  *resRaw.Traceroute,
		}
		entry.TaskLogs = append(entry.TaskLogs, TaskStatRes)
	}
	return entry
}

// parseFloat safely parse string to float64.
func parseFloat(s *string, fieldName string) float64 {
	if s == nil || *s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"field": fieldName, "value": *s}).
			Warn("Failed to parse float value")
		return 0
	}
	return f
}

// parseInt safely parse string to int64.
func parseInt(s *string, fieldName string) int64 {
	if s == nil || *s == "" {
		return 0
	}
	i, err := strconv.ParseInt(*s, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"field": fieldName, "value": *s}).
			Warn("Failed to parse int value")
		return 0
	}
	return i
}
