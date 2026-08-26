package matcher

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/waradu/tabbr/internal/store"
)

func TestQueryMatchingAndOrdering(t *testing.T) {
	now := time.Now().Unix()
	commands := []store.Command{
		{Text: "brew doctor", UseCount: 100, LastUsed: now},
		{Text: "bun run dev", UseCount: 1, LastUsed: now},
		{Text: "bun run dog", UseCount: 10, LastUsed: now},
		{Text: "ubuntu run dev", UseCount: 1000, LastUsed: now},
	}

	got := Query(commands, "BRD")
	gotTexts := make([]string, len(got))
	for i, command := range got {
		gotTexts[i] = command.Command.Text
	}

	want := []string{"bun run dog", "bun run dev", "brew doctor"}
	if !reflect.DeepEqual(gotTexts, want) {
		t.Fatalf("Query() = %q, want %q", gotTexts, want)
	}
}

func TestQueryRequiresAbsoluteOrder(t *testing.T) {
	commands := []store.Command{{
		Text:     "bun run dev",
		UseCount: 1,
		LastUsed: time.Now().Unix(),
	}}

	if got := Query(commands, "bdr"); len(got) != 0 {
		t.Fatalf("Query() returned %d matches, want none", len(got))
	}
}

func BenchmarkQuery(b *testing.B) {
	for _, size := range []int{1_000, 10_000, 100_000} {
		b.Run(fmt.Sprintf("commands_%d", size), func(b *testing.B) {
			commands := make([]store.Command, size)
			now := time.Now().Unix()
			for i := range commands {
				text := fmt.Sprintf("tool-%d execute task", i)
				if i%10 == 0 {
					text = fmt.Sprintf("bun run dev task-%d", i)
				}
				commands[i] = store.Command{
					Text:     text,
					UseCount: i%20 + 1,
					LastUsed: now - int64(i%86_400),
				}
			}

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				Query(commands, "brd")
			}
		})
	}
}
