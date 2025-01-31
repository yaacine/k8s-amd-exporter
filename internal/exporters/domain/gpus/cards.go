package gpus

/* rocm-smi output sample
{
    "card0": {
		"Card series": "Instinct MI210",
		"Card model": "0x0c34",
		"Card vendor": "Advanced Micro Devices, Inc. [AMD/ATI]",
		"Card SKU": "D67301"
        "Subsystem ID": "0x0b0c",
        "Device Rev": "0x01",
        "Node ID": "13",
        "GUID": "63755",
        "PCI Bus": "0000:31:00.0",
        "GFX Version": "gfx9010"
    },
    "card1": {
        "Device Name": "Instinct MI210",
        "Device ID": "0x740c",
        "Device Rev": "0x01",
        "Subsystem ID": "0x0b0c",
        "GUID": "27432",
        "PCI Bus": "0000:34:00.0",
        "Card Series": "AMD INSTINCT MI250 (MCM) OAM AC MBA",
		"Card model": "0x0c34",
		"Card vendor": "Advanced Micro Devices, Inc. [AMD/ATI]",
		"Card SKU": "D67301"
        "Node ID": "13",
        "GFX Version": "gfx9010"
    }
}
*/

// Card contains information about gpu cards within system.
type Card struct {
	Cardseries string `json:"cardseries"`
	Cardmodel  string `json:"cardmodel"`
	Cardvendor string `json:"cardvendor"`
	CardSKU    string `json:"cardsku"`
	PCIBus     string `json:"pcibus"`
	CardGUID   string `json:"guid"`
}

// amd constant values.
const (
	MaxNumGPUDevices               uint   = 24
	GKEVirtualGPUDeviceIDSeparator string = "/vgpu"
	AMDVirtualGPUDeviceIDSeparator string = "/mxgpu"
	AMDResourceName                string = "amd.com/gpu"
)
