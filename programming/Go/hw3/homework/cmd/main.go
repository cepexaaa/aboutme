package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	BlockSize int
	Conv      string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset in input file")
	flag.IntVar(&opts.Limit, "limit", -1, "limit of bytes to read")
	flag.IntVar(&opts.BlockSize, "block-size", 1024, "block size for reading")
	flag.StringVar(&opts.Conv, "conv", "", "conversion options: upper_case, lower_case, trim_spaces")

	flag.Parse()

	if opts.From != "" {
		if _, err := os.Stat(opts.From); os.IsNotExist(err) {
			return nil, fmt.Errorf("file specified in 'from' does not exist: %s", opts.From)
		}
	}

	if opts.To != "" {
		if _, err := os.Stat(opts.To); err == nil {
			return nil, fmt.Errorf("file specified in 'to' already exists: %s", opts.To)
		}
	}

	if opts.Conv != "" {
		convs := strings.Split(opts.Conv, ",")
		hasUpper := false
		hasLower := false
		for _, conv := range convs {
			switch conv {
			case "upper_case":
				hasUpper = true
			case "lower_case":
				hasLower = true
			case "trim_spaces":
				// Ok, do nothing
			default:
				return nil, fmt.Errorf("invalid conversion option: %s", conv)
			}
		}
		if hasUpper && hasLower {
			return nil, fmt.Errorf("cannot use both upper_case and lower_case conversions")
		}
	}

	if opts.Offset < 0 {
		return nil, fmt.Errorf("invalid offset: %d", opts.Offset)
	}

	if opts.Limit < -1 {
		return nil, fmt.Errorf("invalid value of limit: %d", opts.Offset)
	}

	if opts.BlockSize <= 0 {
		return nil, fmt.Errorf("invalid value of block size: %d", opts.Offset)
	}

	return &opts, nil
}

func copyData(input io.Reader, output io.Writer, opts *Options) error {
	buf := make([]byte, opts.BlockSize)
	bytesCopied := 0

	if opts.Offset > 0 {
		_, err := io.CopyN(io.Discard, input, int64(opts.Offset))
		if err != nil {
			return fmt.Errorf("failed to skip offset: %w", err)
		}
	}

	var fullData []byte
	if strings.Contains(opts.Conv, "trim_spaces") {
		var err error
		if opts.Limit == -1 {
			fullData, err = io.ReadAll(input)
		} else {
			fullLimitData := make([]byte, opts.Limit)
			_, err = input.Read(fullLimitData)
			fullData = fullLimitData
		}
		if err != nil {
			return fmt.Errorf("failed to read input data: %w", err)
		}
		fullData, err := modifyText(fullData, opts.Conv)
		if err != nil {
			return fmt.Errorf("failed to apply conversions: %w", err)
		}

		for len(fullData) > 0 {
			chunkSize := opts.BlockSize
			if chunkSize > len(fullData) {
				chunkSize = len(fullData)
			}
			_, err := output.Write(fullData[:chunkSize])
			if err != nil {
				return fmt.Errorf("failed to write data: %w", err)
			}
			fullData = fullData[chunkSize:]
		}

		return nil
	}

	for {
		n, err := input.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to read data: %w", err)
		}

		if n == 0 {
			break
		}

		if opts.Limit != -1 && bytesCopied+n > opts.Limit {
			n = opts.Limit - bytesCopied
		}

		processedData, err := modifyText(buf[:n], opts.Conv)
		if err != nil {
			return fmt.Errorf("failed to apply conversions: %w", err)
		}

		_, err = output.Write(processedData)
		if err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}

		bytesCopied += n

		if opts.Limit != -1 && bytesCopied >= opts.Limit {
			break
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return nil
}

func modifyText(data []byte, conv string) ([]byte, error) {
	if conv == "" {
		return data, nil
	}

	if !utf8.Valid(data) {
		return data, nil
	}

	str := string(data)
	convs := strings.Split(conv, ",")

	for _, c := range convs {
		switch c {
		case "upper_case":
			str = strings.ToUpper(str)
		case "lower_case":
			str = strings.ToLower(str)
		case "trim_spaces":
			str = strings.TrimSpace(str)
		}
	}

	return []byte(str), nil
}

func openInput(opts *Options) (io.Reader, func() error, error) {
	if opts.From != "" {
		file, err := os.Open(opts.From)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open input file: %w", err)
		}
		return file, file.Close, nil
	}
	return os.Stdin, func() error { return nil }, nil
}

func openOutput(opts *Options) (io.Writer, func() error, error) {
	if opts.To != "" {
		file, err := os.Create(opts.To)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create output file: %w", err)
		}
		return file, file.Close, nil
	}
	return os.Stdout, func() error { return nil }, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	input, closeInput, err := openInput(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		if err := closeInput(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to close input:", err)
		}
	}()

	output, closeOutput, err := openOutput(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		if err := closeOutput(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to close output:", err)
		}
	}()

	err = copyData(input, output, opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to copy data:", err)
		os.Exit(1)
	}
}
