# iOS Setup Guide

This document explains how to use the **Image Proxy** on an iPhone or iPad. By following these steps, you’ll be able to route your iOS device’s traffic through a proxy running on your computer, allowing features like HTTPS interception (MITM) if you install the proxy’s certificate.

---
**If you've deployed to proxy to a remote host, skip 1 & 2 and jump to [step 3](https://github.com/ahmedaabouzied/image-proxy/blob/main/ios_instrucitons.md#3-configure-the-ios-wi-fi-proxy-settings)**
## 1. Connect Your iPhone and Computer to the Same Network

1. **Make sure** your computer (running the Image Proxy app) and your iPhone or iPad are both on the **same Wi-Fi network**.  
2. Your computer will host the proxy on an IP address like `192.168.x.x`. You’ll need to point your iPhone’s Wi-Fi proxy settings to that IP.

---

## 2. Launch the Image Proxy App on Your Computer

1. **Double-click** or **open** the Image Proxy application you built or downloaded.  
2. A window should appear, showing:
   - Your **local IP address** (e.g., `192.168.1.10`)
   - Basic instructions on how to set the iPhone’s proxy
   - A link to download the **root CA certificate** if you plan on intercepting HTTPS traffic.

---

## 3. Configure the iOS Wi-Fi Proxy Settings

1. On your **iPhone/iPad**, open the **Settings** app.  
2. Tap **Wi-Fi** and select the Wi-Fi network you’re connected to (the same one as your computer).  
3. Scroll down to **Configure Proxy** (or **HTTP Proxy**) and tap **Manual**.  
4. **Server**: Enter the IP address shown in the Image Proxy window (for example, `192.168.1.10`).  
5. **Port**: Enter the port the proxy is listening on (commonly `8080`).  
6. Leave **Authentication** (username/password) **Off** or as required by your proxy setup.  
7. Tap **Save** or go back to confirm settings.

At this point, any **HTTP** traffic should flow through your local proxy. If you’re **not** doing HTTPS interception (MITM), you can stop here. However, you **will not** see or modify encrypted traffic. To inspect or modify HTTPS traffic, continue with certificate installation.

---

## 4. Install the Root CA Certificate (for HTTPS MITM)

### 4.1 Download the Certificate on iOS

1. In Safari on your iPhone/iPad, visit the URL provided by the Image Proxy window. For example:
   ```
   https://static.praha.aabouzied.com/certs/rootCA.pem
   ```
2. **Tap “Allow”** or **Download** if prompted to download a configuration profile.

*(Alternatively, you can **AirDrop** or **email** the `.pem` or `.crt` file to your iOS device. When you open it, iOS will prompt you to install it as a profile.)*

### 4.2 Install the Downloaded Profile

1. After the file is downloaded, you’ll see a prompt or notice saying **“Profile Downloaded.”**  
2. Open the **Settings** app on your iPhone/iPad.  
3. Tap **General** → **VPN & Device Management**.  
   - On older iOS versions, this might be **Profiles** or **Profile & Device Management**.  
4. You should see the **rootCA** or certificate name under **Downloaded Profile**.  
5. Tap the profile and then tap **Install**. Enter your passcode if prompted.  
6. Tap **Install** again to confirm.

### 4.3 Trust the CA Certificate

Starting with iOS 10.3, Apple requires you to **explicitly enable full trust** for manually installed root certificates:

1. In **Settings**, go to **General** → **About** → **Certificate Trust Settings**.  
2. You should see your installed root certificate listed.  
3. **Toggle ON** to enable full trust for that certificate.  
4. Tap **Continue** at the warning prompt.

At this point, your iPhone trusts the Image Proxy’s CA. HTTPS connections routed through the proxy can now be intercepted and re-encrypted without producing certificate errors in Safari or other apps.

---

## 5. Verification & Troubleshooting

- **Check a site** in Safari (like `https://example.com`) to confirm:
  - The page loads normally.  
  - The proxy can see or modify HTTPS requests if you’ve enabled MITM interception.  
- If Safari warns about “Untrusted Certificate”:
  - Verify you installed and trusted the root CA in **Certificate Trust Settings**.  
  - Make sure the proxy is actually using the same certificate you installed.
- If you **lose internet** or pages don’t load:
  - Double-check the **Server IP** and **Port** in Wi-Fi proxy settings.  
  - Confirm the Image Proxy app on your computer is still running and not blocked by a firewall.

---

## 6. Done!

- You have now configured your iOS device to route traffic through **Image Proxy**.  
- Enjoy the ability to inspect or modify traffic for debugging, development, or testing purposes.

**Note**: If you want to **stop** using the proxy, just go back to **Settings** → **Wi-Fi** → **Configure Proxy** → **Off** (or **Automatic**) on your iPhone.
