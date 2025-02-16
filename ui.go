package main

import (
	"fmt"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func getLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// Skip IPv6 or loopback
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no connected network interface found")
}

func startUI() {
	// Create the fyne app
	a := app.New()
	w := a.NewWindow("Image Proxy")

	// Attempt to get local IP
	localIP, err := getLocalIP()
	if err != nil {
		localIP = "Could not detect local IP."
	}

	// Instructions for installing root CA on iOS
	// (You can adjust or expand these as needed.)
	instructions := `
1. Connect your iPhone to the same Wi-Fi network as this computer.

2. On your iPhone, open Safari and go to:
   https://static.praha.aabouzied.com/certs/rootCA.pem

3. Tap "Allow" or "Install" when prompted to download the profile.

4. Open the Settings app:
   - Go to General -> VPN & Device Management (or "Profiles" on older iOS versions).
   - Find and tap the downloaded profile (rootCA).
   - Tap "Install" (enter your passcode if required).
   - Then go to General -> About -> Certificate Trust Settings.
   - Enable full trust for the installed root certificate.

5. To use this local proxy from your iPhone, go to:
   Settings -> Wi-Fi -> (Your Network) -> Configure Proxy -> Manual
   - Set Server to the IP address %s
   - Set it to 8080
`

	instructionLabel := widget.NewLabel(strings.TrimSpace(fmt.Sprintf(instructions, localIP)))
	instructionLabel.Wrapping = fyne.TextWrapWord

	// Build UI layout
	content := container.NewVBox(
		widget.NewLabel("Welcome to safe image proxy"),
		instructionLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
