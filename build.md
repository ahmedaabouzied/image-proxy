# Building the Image Proxy App

This document covers how to build the Fyne-based Image Proxy application for both macOS (`.app` bundle) and Windows (`.exe`). The same instructions can be adapted for Linux or other platforms as needed.

---

## Prerequisites

- **Go** (1.18+) installed.
- **Fyne CLI** (optional, but recommended for packaging a `.app` on macOS).
  - Install with:
    ```bash
    go install fyne.io/fyne/v2/cmd/fyne@latest
    ```
- Basic familiarity with terminal commands on macOS and Windows.

---

## 1. Getting the Source

Clone or download the source code for the Image Proxy. 

## 2. Building for macOS

### A) Building a Plain Binary

The simplest method is to compile a single binary (CLI-based). You can run it directly from the terminal:

```bash
# On macOS (Intel/AMD)
go build -o ImageProxy .

# On macOS (Apple Silicon)
GOARCH=arm64 go build -o ImageProxy .
```

This produces a single executable named `ImageProxy` in the current directory. When you double-click or run it in the terminal (`./ImageProxy`), it will open the Fyne UI window.

### B) Creating a .app Bundle

To produce a “proper” macOS `.app` bundle with an icon, name, etc., use the **Fyne CLI** packaging feature:

1. **Install the Fyne CLI** (if you haven’t already):
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```
2. **Run the packaging command** from within your project directory:
   ```bash
   fyne package -os darwin -icon youricon.png -name "Image Proxy"
   ```
   - `-os darwin` specifies macOS.  
   - `-icon youricon.png` is your desired icon file (PNG).  
   - `-name "Image Proxy"` gives your app the display name “Image Proxy.”  

Fyne will create a `Image Proxy.app` in your current directory. You can **drag** this `.app` to your **Applications** folder or **zip** and distribute it.

**Note**:  
- If you distribute widely, you may need to [sign and notarize](https://developer.apple.com/developer-id/) the app to avoid macOS Gatekeeper warnings. For internal usage, you can skip that or manually allow the app to run via “System Preferences > Security & Privacy.”

---

## 3. Building for Windows

### A) Simple .exe

On Windows, open a **Command Prompt** or **PowerShell** and run:

```powershell
go build -o ImageProxy.exe .
```

If you’re on macOS or Linux and want to cross-compile for Windows (64-bit), you can do:

```bash
GOOS=windows GOARCH=amd64 go build -o ImageProxy.exe .
```

You’ll get a `ImageProxy.exe` that you can double-click on Windows. When launched, it will open the Fyne UI window.

### B) Fyne CLI Packaging for Windows

Similarly, you can package a Windows executable with an icon resource by using the Fyne CLI on your Windows machine (or cross-compiling from another OS):

```bash
fyne package -os windows -icon youricon.png -name "Image Proxy"
```

This creates `Image Proxy.exe` with metadata like icon, version info, etc., embedded.

---

## 4. Running the App

Once you have your built application:

- **macOS**: 
  - If it’s a `.app` bundle, double-click `Image Proxy.app`.
  - If it’s a plain binary (`ImageProxy`), open Terminal, `cd` to the directory, and run `./ImageProxy`.
- **Windows**:
  - Double-click `ImageProxy.exe` (or run `ImageProxy.exe` from Command Prompt).

On launch, the Fyne window will appear, showing:
- Your local IP address.
- Instructions to download the root CA from your specified URL.
- Steps for installing it on an iOS device.
- Steps to set the iPhone’s manual proxy settings to point to your local IP and port.

---

## 5. Common Troubleshooting

1. **Firewall Issues**:  
   - Ensure your macOS or Windows firewall allows inbound connections on the chosen port (e.g., 8080).
2. **Missing CA Certificate**:  
   - If you’re doing HTTPS interception, the iOS device needs to trust the certificate.  
3. **No Local IP Found**:  
   - If you’re on a VPN or have multiple network interfaces, the code might pick the wrong interface. Adjust your IP detection logic if needed.
4. **Gatekeeper on macOS**:  
   - You may need to right-click the `.app` and select “Open” the first time to override Gatekeeper if it is unsigned/notarized.

---

## 6. Further Reading

- [Fyne Documentation](https://developer.fyne.io/) – official docs on creating GUI apps in Go.  
- [Cross-Compiling Go](https://go.dev/doc/install/source#environment) – official docs on targeting different platforms and architectures.  
- [Apple Developer ID / Notarization](https://developer.apple.com/developer-id/) – if you plan to distribute a macOS app widely without warnings.  

---

*That’s it! You should now have a functional proxy app for macOS and Windows.*  
