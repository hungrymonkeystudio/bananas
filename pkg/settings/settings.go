// setting screen for bananas
// core functionality: change time control

package settings 

import (
	"strconv"
    "os"
    "bufio"
	resourcepath "bananas/pkg/resourcepath"
	colors "bananas/pkg/colors"
	tea "github.com/charmbracelet/bubbletea"
)

const settingInstructions = "UP/DOWN/LEFT/RIGHT to move\nENTER to select\nESC to close settings page"

type SettingsModel struct {
    Show bool;
    options []string;
    optionIdx int;
    times []int;
    timeIdx int;
    ActiveTime int;
}

func NewSettingsModel() SettingsModel {
    s := SettingsModel{
        Show: false,
        options: []string{"timer", "restart", "quit"},
        optionIdx: 0,
        times: []int{15, 30, 60, 120},
        timeIdx: 1,
        ActiveTime: 30,
    }
    readSettings(&s)
    return s
}

func (m SettingsModel) Init() tea.Cmd {
    return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            m.optionIdx = (m.optionIdx-1+len(m.options)) % len(m.options)
        case "down":
            m.optionIdx = (m.optionIdx+1) % len(m.options)
        case "left":
            if m.options[m.optionIdx] == "timer" {
                m.timeIdx = (m.timeIdx-1+len(m.times)) % len(m.times)
            }
        case "right":
            if m.options[m.optionIdx] == "timer" {
                m.timeIdx = (m.timeIdx+1) % len(m.times)
            }
        case "enter":
            if m.options[m.optionIdx] == "quit" {
                return m, tea.Quit
            } else if m.options[m.optionIdx] == "restart" {
                return m, func () tea.Msg { return m }
            } else if m.options[m.optionIdx] == "timer" {
                m.ActiveTime = m.times[m.timeIdx]
                m.writeSettings()
                return m, func() tea.Msg { return m }
            }
        }
    }
    return m, nil
}

func (m SettingsModel) View() string {
    output := colors.Yellow.Render("Settings Page") + "\n"
    if m.options[m.optionIdx] == "timer" {
        output += colors.White.Render("timer: ")
        for timeIdx, times := range m.times {
            if timeIdx == m.timeIdx {
                if m.times[m.timeIdx] == m.ActiveTime {
                    output += colors.White.Underline(true).Render(strconv.Itoa(times)) + " "
                } else {
                    output += colors.Gray.Underline(true).Render(strconv.Itoa(times)) + " "
                }
            } else {
                if m.times[timeIdx] == m.ActiveTime {
                    output += colors.White.Render(strconv.Itoa(times)) + " "
                } else {
                    output += colors.Gray.Render(strconv.Itoa(times)) + " "
                }
            }
        }
    } else {
        output += colors.Gray.Render("timer: ")
        for timeIdx, times := range m.times {
            if timeIdx == m.timeIdx {
                if m.times[m.timeIdx] == m.ActiveTime {
                    output += colors.Gray.Underline(true).Render(strconv.Itoa(times)) + " "
                } else {
                    output += colors.Gray.Render(strconv.Itoa(times)) + " "
                }
            } else {
                output += colors.Gray.Render(strconv.Itoa(times)) + " "
            }
        }
    }
    output += "\n"
    if m.options[m.optionIdx] == "restart" {
        output += colors.White.Render("restart")
    } else {
        output += colors.Gray.Render("restart")
    }
    output += "\n"
    if m.options[m.optionIdx] == "quit" {
        output += colors.White.Render("quit")
    } else {
        output += colors.Gray.Render("quit")
    }
    output += "\n"
    output += "\n" + colors.Instructions.Render(settingInstructions)
    return output
}

func (m SettingsModel) writeSettings() {
    // create file
    file, err := os.Create("user_settings.json")
	if err != nil {
		return
	}
    defer file.Close()

	// Write to the file
	_, err = file.WriteString(strconv.Itoa(m.ActiveTime))
	if err != nil {
		return
	}
}

func readSettings(m *SettingsModel) {
    // read from settings file
	basePath := resourcepath.GetResourcePath()
    file, err := os.Open(basePath+"/settings.json")
    if err != nil {
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    textVal := ""
    for scanner.Scan() {
        textVal = scanner.Text()
    }
    num, err := strconv.Atoi(textVal)
    if err != nil {
        return
    }
    // check that num matches at least 1 of the values in times
    matches := false
    for _, time := range m.times {
        if time == num {
            matches = true
        }
    }
    if !matches {
        return
    }
    m.ActiveTime = num
}
