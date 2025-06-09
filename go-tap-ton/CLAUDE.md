# Claude Development Notes

## Project Overview
A terminal-based tap tempo analyzer that measures timing between spacebar presses and calculates BPM with intelligent averaging.

## Key Implementation Details

### Architecture
- Built using Bubble Tea TUI framework for real-time key input
- Model-View-Update pattern with separated concerns
- Structured into focused functions for maintainability

### Core Features
1. **Real-time timing measurement** using Go's `time.Now()`
2. **Ao8 averaging algorithm** (excludes highest/lowest from up to 8 samples)
3. **Outlier detection** (2x slower than recent average triggers skip)
4. **Automatic reset** after 3 consecutive outliers
5. **60fps frame conversion** (16.67ms per frame)

### Code Structure
- `model`: Main state structure with timing history
- `handleFirstTap()`: Initialize measurement
- `calculateMetrics()`: Convert duration to ms/frames/BPM
- `isOutlierBpm()`: Detect timing anomalies
- `calculateAo8()`: Implement averaging algorithm
- `handleValidTap()`/`handleOutlier()`: Process measurements

### Development Commands
```bash
go run main.go          # Run the application
go build               # Build binary
go mod tidy            # Update dependencies
```

### Dependencies
- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- Standard library: `fmt`, `sort`, `time`

### UI Design Philosophy
- Minimal, no decorative borders or excessive formatting
- Fixed-width columns for clean alignment
- Informative but not cluttered output
- Unix tool aesthetic - functional over flashy

### Testing Notes
- Test with various tempo ranges (slow/fast)
- Verify outlier detection with intentional pauses
- Confirm Ao8 calculation accuracy
- Check reset behavior after consecutive outliers