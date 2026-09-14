# Project 1: Automated System Log Management & Archive

## 📌 Situation
Unmanaged log files accumulate indefinitely on Linux servers, leading to storage exhaustion, cloud budget overruns, and system degradation. Manual intervention creates operational drag.

## 🛠️ Task
Develop a lightweight, high-reliability Bash automation engine to scan log spaces, archive stale payloads (files > 7 days old) via secure `gzip` standard compression, and clear raw execution space automatically.

## 🚀 Action
* **Scripting Language:** Bash 4.x+
* **Core Utilities:** Built utilizing native `find` architecture with `-mtime` constraints to locate targeted structures dynamically.
* **Storage Optimization:** Implemented `tar -czf` configurations to safely shrink storage footings before cleaning raw assets.
* **Observability:** Programmed logging hooks that compile an isolated audit execution summary (`log_cleanup_report.txt`) recording every systemic storage reduction loop.

## 📈 Result
* Eliminated manual administration oversight for log retention pipelines.
* Drastically reduced data footprints through automated compression routines.
* Provided system administrators with a scannable, append-only execution history.
