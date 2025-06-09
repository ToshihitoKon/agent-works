package main

import (
    "fmt"
    "sort"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    lastTime      time.Time
    firstPress    bool
    output        []string
    bpmHistory    []float64
    outlierCount  int
}

func initialModel() model {
    return model{
        firstPress: true,
        output:     []string{},
        bpmHistory: []float64{},
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) handleFirstTap(currentTime time.Time) model {
    m.output = append(m.output, "First tap - measurement started")
    m.lastTime = currentTime
    m.firstPress = false
    return m
}

func (m model) calculateMetrics(elapsed time.Duration) (float64, float64, float64) {
    elapsedMs := float64(elapsed.Nanoseconds()) / 1000000.0
    frames := elapsedMs / 16.666666
    elapsedSec := elapsedMs / 1000.0
    currentBpm := 60.0 / elapsedSec
    return elapsedMs, frames, currentBpm
}

func (m model) isOutlierBpm(currentBpm float64) bool {
    if len(m.bpmHistory) == 0 {
        return false
    }
    
    var recentAvg float64
    for _, bpm := range m.bpmHistory {
        recentAvg += bpm
    }
    recentAvg /= float64(len(m.bpmHistory))
    
    return currentBpm < recentAvg/2.0
}

func (m model) handleOutlier(elapsedMs, frames, currentBpm float64) model {
    m.outlierCount++
    result := fmt.Sprintf("%8.2fms  %6.1fF  %6.1f BPM  [OUTLIER - SKIP %d]", 
        elapsedMs, frames, currentBpm, m.outlierCount)
    m.output = append(m.output, result)
    
    if m.outlierCount >= 3 {
        m.bpmHistory = []float64{}
        m.outlierCount = 0
        m.output = append(m.output, "BPM history reset due to consecutive outliers")
    }
    return m
}

func (m model) calculateAo8(history []float64) float64 {
    if len(history) < 3 {
        // 3未満の場合は単純平均
        var sum float64
        for _, bpm := range history {
            sum += bpm
        }
        return sum / float64(len(history))
    }
    
    // 履歴をコピーしてソート
    sorted := make([]float64, len(history))
    copy(sorted, history)
    sort.Float64s(sorted)
    
    // 最大値と最小値を除いた中央部分の平均
    if len(sorted) <= 2 {
        return sorted[0]
    }
    
    start := 1
    end := len(sorted) - 1
    var sum float64
    for i := start; i < end; i++ {
        sum += sorted[i]
    }
    return sum / float64(end-start)
}

func (m model) handleValidTap(elapsedMs, frames, currentBpm float64) model {
    m.outlierCount = 0
    m.bpmHistory = append(m.bpmHistory, currentBpm)
    if len(m.bpmHistory) > 8 {
        m.bpmHistory = m.bpmHistory[1:]
    }
    
    avgBpm := m.calculateAo8(m.bpmHistory)
    
    avgLabel := "avg"
    if len(m.bpmHistory) >= 3 {
        avgLabel = fmt.Sprintf("Ao%d", len(m.bpmHistory))
    }
    
    result := fmt.Sprintf("%8.2fms  %6.1fF  %6.1f BPM  [%s: %6.1f]", 
        elapsedMs, frames, currentBpm, avgLabel, avgBpm)
    m.output = append(m.output, result)
    return m
}

func (m model) handleSpaceKey() model {
    currentTime := time.Now()
    
    if m.firstPress {
        return m.handleFirstTap(currentTime)
    }
    
    elapsed := currentTime.Sub(m.lastTime)
    elapsedMs, frames, currentBpm := m.calculateMetrics(elapsed)
    
    if m.isOutlierBpm(currentBpm) {
        m = m.handleOutlier(elapsedMs, frames, currentBpm)
    } else {
        m = m.handleValidTap(elapsedMs, frames, currentBpm)
    }
    
    m.lastTime = currentTime
    return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case " ":
            m = m.handleSpaceKey()
        }
    }
    return m, nil
}

func (m model) View() string {
    view := "Tap Tempo Analyzer - Press SPACE to tap, Q to quit\n\n"
    
    for _, line := range m.output {
        view += line + "\n"
    }
    
    return view
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("エラー: %v", err)
    }
}