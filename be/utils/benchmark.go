package utils
import (
	"fmt"
	"log"
	"strings"
	"time"
)
func Benchmark(label string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		category := detectCategory(label)
		level := classify(category, elapsed)
		color := map[string]string{
			"good":  "\033[32m", // green
			"warn":  "\033[33m", // yellow
			"bad":   "\033[31m", // red
			"reset": "\033[0m",
		}[level]
		elapsedStr := formatDuration(elapsed)
		log.Printf("%s%s %s - %s%s",
			color,
			start.Format("2006-01-02 15:04:05"),
			elapsedStr,
			label,
			colorizeReset(),
		)
	}
}
func formatDuration(d time.Duration) string {
	ns := d.Nanoseconds()
	switch {
	case ns == 0:
		return "<100ns" // special case for ultra-fast response
	case ns < 1_000:
		return fmt.Sprintf("%dns", ns)
	case ns < 1_000_000:
		return fmt.Sprintf("%.3fµs", float64(ns)/1_000)
	default:
		return fmt.Sprintf("%.3fms", float64(ns)/1_000_000)
	}
}
func detectCategory(label string) string {
	switch {
	case strings.HasPrefix(label, "Rdb"):
		return "rdb"
	case strings.HasPrefix(label, "Pdb"):
		return "pdb"
	case strings.HasPrefix(label, "Memory") || strings.HasPrefix(label, "Mem"):
		return "mem"
	default:
		return "generic"
	}
}
func classify(category string, d time.Duration) string {
	switch category {
	case "mem":
		if d < 300*time.Microsecond {
			return "good"
		} else if d <= 800*time.Microsecond {
			return "warn"
		}
	case "rdb":
		if d < 2*time.Millisecond {
			return "good"
		} else if d <= 5*time.Millisecond {
			return "warn"
		}
	case "pdb":
		if d < 5*time.Millisecond {
			return "good"
		} else if d <= 15*time.Millisecond {
			return "warn"
		}
	case "generic":
		if d < 10*time.Millisecond {
			return "good"
		} else if d <= 30*time.Millisecond {
			return "warn"
		}
	}
	return "bad"
}
func colorizeReset() string {
	return "\033[0m"
}
