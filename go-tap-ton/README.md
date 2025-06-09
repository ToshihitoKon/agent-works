# go-tap-ton

A minimalist terminal-based tap tempo analyzer written in Go.

## Features

- Real-time BPM calculation from spacebar taps
- 60fps frame conversion (useful for game development)
- Ao8 averaging (Average of 8, excluding highest and lowest values)
- Outlier detection and automatic reset
- Clean, distraction-free terminal interface

## Installation

```bash
go install github.com/ToshihitoKon/go-tap-ton@latest
```

Or clone and build:

```bash
git clone https://github.com/ToshihitoKon/go-tap-ton.git
cd go-tap-ton
go build
```

## Usage

```bash
./go-tap-ton
```

- Press **SPACE** to tap and measure timing
- Press **Q** to quit
- First tap starts measurement
- Subsequent taps show:
  - Elapsed time in milliseconds
  - Frame count at 60fps
  - Current BPM
  - Average BPM using Ao8 method

## Output Format

```
  123.45ms   7.4F  120.0 BPM  [Ao5: 118.2]
```

- **123.45ms**: Time since last tap
- **7.4F**: Equivalent frames at 60fps (16.67ms per frame)
- **120.0 BPM**: Current beats per minute
- **[Ao5: 118.2]**: Average of 5 measurements (excluding highest/lowest)

## Outlier Detection

The tool automatically detects outliers (taps that are 2x slower than the recent average) and excludes them from calculations. After 3 consecutive outliers, the history is reset to allow for tempo changes.

## Technical Details

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI
- Uses Go's high-precision time measurement
- Implements Ao8 averaging algorithm (common in speedcubing)
- Maintains up to 8 recent measurements for averaging

## License

MIT License