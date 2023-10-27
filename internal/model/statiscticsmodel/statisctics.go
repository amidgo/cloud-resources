package statiscticsmodel

type Statistics struct {
	Availability  float32 `json:"availability"`
	CostTotal     float32 `json:"cost_total"`
	DBCPU         float32 `json:"db_cpu"`
	DBCPULoad     float32 `json:"db_cpu_load"`
	DBRAM         float32 `json:"db_ram"`
	DBRAMLoad     float32 `json:"db_ram_load"`
	ID            float32 `json:"id"`
	Last1         float32 `json:"last1"`
	Last15        float32 `json:"last15"`
	Last5         float32 `json:"last5"`
	LastDay       float32 `json:"lastDay"`
	LastHour      float32 `json:"lastHour"`
	LastWeek      float32 `json:"lastWeek"`
	OfflineTime   float32 `json:"offline_time"`
	Online        bool    `json:"online"`
	OnlineTime    int     `json:"online_time"`
	Requests      int     `json:"requests"`
	RequestsTotal int     `json:"requests_total"`
	ResponseTime  int     `json:"response_time"`
	Timestamp     string  `json:"timestamp"`
	UserID        int     `json:"user_id"`
	UserName      string  `json:"user_name"`
	VMCPU         float32 `json:"vm_cpu"`
	VMCPULoad     float32 `json:"vm_cpu_load"`
	VMRAM         float32 `json:"vm_ram"`
	VMRAMLoad     float32 `json:"vm_ram_load"`
}
