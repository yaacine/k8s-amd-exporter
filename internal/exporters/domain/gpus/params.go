package gpus

// AMDParamsHandler defines function signature to return amd metrics data.
type AMDParamsHandler func() AMDParams

// AMDParams contains all metrics to be collected from environment.
type AMDParams struct {
	CoreEnergy     [768]float64
	CoreBoost      [768]float64
	SocketEnergy   [8]float64
	SocketPower    [8]float64
	PowerLimit     [8]float64
	ProchotStatus  [8]float64
	Sockets        uint
	Threads        uint
	ThreadsPerCore uint
	NumGPUs        uint
	GPUDevID       [MaxNumGPUDevices]float64
	GPUDevPCIId    [MaxNumGPUDevices]float64
	GPUPowerCap    [MaxNumGPUDevices]float64
	GPUPower       [MaxNumGPUDevices]float64
	GPUTemperature [MaxNumGPUDevices]float64
	GPUSCLK        [MaxNumGPUDevices]float64
	GPUMCLK        [MaxNumGPUDevices]float64
	GPUUsage       [MaxNumGPUDevices]float64
	GPUMemoryUsage [MaxNumGPUDevices]float64
}

// Init initializes amd metrics.
func (amdParams *AMDParams) Init() {
	amdParams.Sockets = 0
	amdParams.Threads = 0
	amdParams.ThreadsPerCore = 0

	amdParams.NumGPUs = 0

	for socketLoopCounter := 0; socketLoopCounter < len(amdParams.SocketEnergy); socketLoopCounter++ {
		amdParams.SocketEnergy[socketLoopCounter] = -1
		amdParams.SocketPower[socketLoopCounter] = -1
		amdParams.PowerLimit[socketLoopCounter] = -1
		amdParams.ProchotStatus[socketLoopCounter] = -1
	}

	for logicalCoreLoopCounter := 0; logicalCoreLoopCounter < len(amdParams.CoreEnergy); logicalCoreLoopCounter++ {
		amdParams.CoreEnergy[logicalCoreLoopCounter] = -1
		amdParams.CoreBoost[logicalCoreLoopCounter] = -1
	}

	for gpuLoopCounter := 0; gpuLoopCounter < len(amdParams.GPUDevID); gpuLoopCounter++ {
		amdParams.GPUDevID[gpuLoopCounter] = -1
		amdParams.GPUPowerCap[gpuLoopCounter] = -1
		amdParams.GPUPower[gpuLoopCounter] = -1
		amdParams.GPUTemperature[gpuLoopCounter] = -1
		amdParams.GPUSCLK[gpuLoopCounter] = -1
		amdParams.GPUMCLK[gpuLoopCounter] = -1
		amdParams.GPUUsage[gpuLoopCounter] = -1
		amdParams.GPUMemoryUsage[gpuLoopCounter] = -1
	}
}
