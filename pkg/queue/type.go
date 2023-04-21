package queue

type AlertMetrics struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Current string `json:"current"`
	Expect  string `json:"expect"`
	Cluster string `json:"cluster"`
}

type CollectMetrics struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Cluster string `json:"cluster"`
}
