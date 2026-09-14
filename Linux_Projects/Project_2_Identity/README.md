# Project 2: Enterprise Identity Provisioning & Access Hardening

## 📌 Situation
Manual onboarding of enterprise infrastructure accounts introduces human error, inconsistent group allocations, and permissions sprawling out of compliance with strict governance models (e.g., GDPR, NIS2, SOC2).

## 🛠️ Task
Build a hardened shell-based automation program capable of processing distinct user roles, checking system dependencies, isolating directories, and assigning secure group structures cleanly.

## 🚀 Action
* **IAM Core Integration:** Configured loops parsing dual values (`username:group`) to safely align human tasks to systemic roles.
* **Least-Privilege Enforcement:** Integrated explicit post-provision loops modifying access matrices (`chmod 700`), closing off cross-tenant directories to arbitrary visibility.
* **Defensive Coding Controls:** Engineered multi-layer evaluation checkpoints verifying execution privileges (`$EUID`) and preventing existing user record collisions.

## 📈 Result
* Guaranteed immediate operational compliance with Zero-Trust identity standards.
* Eradicated misconfigurations by structuring standardized user profiles via automated variables.
* Accelerated staff workspace initialization timelines to zero manual setup.
