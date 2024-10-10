package resample

// LinearResampler resamples the input samples from the original sample rate to the target sample rate.
type LinearResampler struct {
	inputRate  float64
	outputRate float64
}

func NewLinearResampler(inputRate, outputRate int) *LinearResampler {
	return &LinearResampler{
		inputRate:  float64(inputRate),
		outputRate: float64(outputRate),
	}
}

func (lr *LinearResampler) ResampleS16(samples []byte) []byte {
	if len(samples)%2 != 0 {
		return nil
	}

	ratio := lr.outputRate / lr.inputRate
	outputLength := int(float64(len(samples)/2)*ratio) * 2
	output := make([]byte, outputLength)

	for i := 0; i < outputLength; i += 2 {
		inputIndex := float64(i/2) / ratio
		index := int(inputIndex)
		frac := inputIndex - float64(index)

		if index*2+2 < len(samples) {
			sample1 := int16(samples[index*2]) | int16(samples[index*2+1])<<8
			sample2 := int16(samples[index*2+2]) | int16(samples[index*2+3])<<8
			interpolated := int16(float64(sample1)*(1-frac) + float64(sample2)*frac)
			output[i] = byte(interpolated)
			output[i+1] = byte(interpolated >> 8)
		} else {
			sample := int16(samples[index*2]) | int16(samples[index*2+1])<<8
			output[i] = byte(sample)
			output[i+1] = byte(sample >> 8)
		}
	}
	return output
}

func LinearResamplePCM(data []byte, inputRate, outputRate int) ([]byte, error) {
	if inputRate == outputRate {
		return data, nil
	}
	lr := NewLinearResampler(inputRate, outputRate)
	return lr.ResampleS16(data), nil
}
