package main

import (
	"context"
	"fmt"
	"os/user"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/Depado/scarecrow/cmd"
	"github.com/Depado/scarecrow/ui"
)

// Main command that will be run when no other command is provided on the
// command-line
var rootCmd = &cobra.Command{
	Use: "scarecrow",
	Run: func(c *cobra.Command, args []string) {
		err := fx.New(
			fx.NopLogger,
			fx.Provide(
				cmd.NewConf, cmd.NewLogger,
			),
			fx.Invoke(hw),
		).Start(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("unable to start")
		}
	},
}

func drawHeader(s tcell.Screen, ram, cpu, temp float64) {
	col := 0
	for _, r := range "Scarecrow " {
		s.SetContent(col, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLimeGreen).Attributes(tcell.AttrBold))
		col++
	}

	rt := fmt.Sprintf("RAM:%3d%%", int(ram))
	ct := fmt.Sprintf("CPU:%5.2f", cpu)
	tt := fmt.Sprintf("TEMP:%3d°", int(temp))
	hsize := len(rt) + len(ct) + len(tt) + 2
	w, _ := s.Size()

	if w < hsize { // TODO: Where to display when not enough space
		return
	}
	col = w - hsize

	rs := ui.GetFloatStyle(ram)
	for _, r := range rt {
		s.SetContent(col, 0, r, nil, rs)
		col++
	}

	s.SetContent(col, 0, ' ', nil, tcell.StyleDefault)
	col++

	rs = ui.GetFloatStyle(cpu)
	for _, r := range ct {
		s.SetContent(col, 0, r, nil, rs)
		col++
	}
	s.SetContent(col, 0, ' ', nil, tcell.StyleDefault)
	col++

	rs = ui.GetFloatStyle(temp)
	for _, r := range tt {
		s.SetContent(col, 0, r, nil, rs)
		col++
	}
	for i := 0; i < w; i++ {
		s.SetContent(i, 1, '─', nil, tcell.StyleDefault)
	}
}

func hw(c *cmd.Conf, l zerolog.Logger) {
	s, err := tcell.NewScreen()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to init screen")
	}

	if err := s.Init(); err != nil {
		l.Fatal().Err(err).Msg("unable to init screen")
	}

	l.Info().Msg("Hello world")
	usr, err := user.Current()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to get user dir")
	}

	s.Clear()
	run := true
	for run {
		// Update screen
		t, _ := host.SensorsTemperatures()
		var cpu float64
		var count float64
		for _, v := range t {
			if strings.Contains(v.SensorKey, "core") || strings.Contains(v.SensorKey, "cpu") {
				count++
				cpu += v.Temperature
			}
		}
		lavg, err := load.Avg()
		if err != nil {
			l.Fatal().Err(err).Msg("unable to get load average")
		}
		v, _ := mem.VirtualMemory()
		drawHeader(s, v.UsedPercent, lavg.Load1, cpu/count)
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Clear()
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				run = false
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		}
	}
	log.Info().Str("home", usr.HomeDir).Msg("User home found")
}

func main() {
	// Initialize Cobra and Viper
	cmd.AddAllFlags(rootCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("unable to start")
	}
}
