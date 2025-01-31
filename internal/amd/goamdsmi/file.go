package goamdsmi

func GO_cpu_init() bool                                { return true }
func GO_gpu_init() bool                                { return true }
func GO_cpu_number_of_sockets_get() int                { return 0 }
func GO_cpu_number_of_threads_get() int                { return 0 }
func GO_cpu_threads_per_core_get() int                 { return 0 }
func GO_cpu_core_energy_get(i int) int                 { return 0 }
func GO_cpu_core_boostlimit_get(i int) int             { return 0 }
func GO_cpu_socket_energy_get(i int) int               { return 0 }
func GO_cpu_socket_power_get(i int) int                { return 0 }
func GO_cpu_socket_power_cap_get(i int) int            { return 0 }
func GO_cpu_prochot_status_get(i int) int              { return 0 }
func GO_gpu_num_monitor_devices() int                  { return 0 }
func GO_gpu_dev_id_get(i int) int                      { return 0 }
func GO_gpu_dev_power_cap_get(i int) int               { return 0 }
func GO_gpu_dev_power_get(i int) int                   { return 0 }
func GO_gpu_dev_temp_metric_get(i, j, k int) int       { return 0 }
func GO_gpu_dev_gpu_clk_freq_get_sclk(i int) int       { return 0 }
func GO_gpu_dev_gpu_clk_freq_get_mclk(i int) int       { return 0 }
func GO_gpu_dev_gpu_busy_percent_get(i int) int        { return 0 }
func GO_gpu_dev_gpu_memory_busy_percent_get(i int) int { return 0 }
