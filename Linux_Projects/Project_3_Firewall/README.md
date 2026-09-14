# Project 3: Automated Host-Based Firewall Security Hardening

## 📌 Situation
Default deployments of network cloud instances are highly susceptible to malicious brute-force attempts and lateral system scans due to unhardened or overly permissive network rule pools.

## 🛠️ Task
Design an infrastructure-as-code shell mechanism to automatically clear unverified traffic channels and lock server ingress points to a strict, minimal operational baseline.

## 🚀 Action
* **Zero-Trust Network Profile:** Enforced drop-all security rules natively (`ufw default deny incoming`) to ensure implicit denial by default.
* **Role-Based Port Management:** Selectively whitelisted critical service lanes including Transport Control Protocol ports `22` (SSH), `80` (HTTP), and `443` (HTTPS) with built-in descriptive comments.
* **Automation Assurance:** Scripted `--force` overrides to allow complete automated provisioning without requiring interactive human confirmation.

## 📈 Result
* Substantially minimized host-based external exploit vectors.
* Replaced manually applied network topologies with a programmatic, verifiable deployment file.
* Created a highly portable script that can be used across multi-region server nodes to scale network edge security uniformly.
