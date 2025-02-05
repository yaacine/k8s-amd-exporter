package amd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"

	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/gpus"
)

const (
	cardKeyName = "card"
)

func GetGpuProductNames() ([24]gpus.Card, error) {
	var result [24]gpus.Card

	rocmsmiRawJSON, err := runROCMSMI()
	if err != nil {
		slog.Error("running rocm-smi --showproductname --showid --showbus --json", slog.String("error", err.Error()))

		return result, fmt.Errorf("unable to get GPU product names: %w", err)
	}

	slog.Debug("rocm-smi product information", slog.String("json", string(rocmsmiRawJSON)))

	// format for json unmarshalling
	productnamesJSON := strings.ReplaceAll(strings.ToLower(string(rocmsmiRawJSON)), " ", "")

	dec := json.NewDecoder(strings.NewReader(productnamesJSON))

	for {
		var cards map[string]gpus.Card

		err := dec.Decode(&cards)
		if err == io.EOF {
			break
		}

		if err != nil {
			return result, fmt.Errorf("decoding card information from rocm-smi: %w", err)
		}
		// iterate over each card
		for k := range cards {
			deviceID, err := strconv.Atoi(strings.TrimPrefix(k, cardKeyName))
			if err == nil {
				result[deviceID] = cards[k]
			}
		}
	}

	return result, nil
}

// runROCMSMI executes rocm-smi to get information about GPU cards and PCI bus addresses.
func runROCMSMI() ([]byte, error) {
	// rocm-smi output in json format
	slog.Info("running rocm-smi --showproductname --showid --showbus --json")

	output, err := exec.Command("rocm-smi", "--showproductname", "--showid", "--showbus", "--json").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to run rocm-smi: %w", err)
	}

	return output, nil
}
