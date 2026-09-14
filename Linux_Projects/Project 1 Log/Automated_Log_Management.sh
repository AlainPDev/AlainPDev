#!/bin/bash
# ==============================================================================
# Script Name:  manage_logs.sh
# Description:  Automates compression of logs older than 7 days and cleans space.
# Author:       Your Name
# ==============================================================================

# Configuration
LOG_DIR="./mock_logs"
BACKUP_DIR="./archived_logs"
DAYS_TO_KEEP=7
REPORT_FILE="./log_cleanup_report.txt"

# Ensure directories exist
mkdir -p "$LOG_DIR" "$BACKUP_DIR"

echo "=== Log Management Execution: $(date) ===" > "$REPORT_FILE"

# Find files older than X days and compress them into the backup directory
echo "[*] Archiving logs older than $DAYS_TO_KEEP days..." >> "$REPORT_FILE"
find "$LOG_DIR" -type f -mtime +"$DAYS_TO_KEEP" -name "*.log" | while read -r file; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        tar -czf "$BACKUP_DIR/${filename}_$(date +%F).tar.gz" -C "$LOG_DIR" "$filename"
        rm "$file"
        echo "Successfully archived and removed: $filename" >> "$REPORT_FILE"
    fi
done

# Output current disk status of the backup folder
echo "" >> "$REPORT_FILE"
echo "[*] Current Backup Storage Status:" >> "$REPORT_FILE"
du -sh "$BACKUP_DIR" >> "$REPORT_FILE"

echo "=== Execution Finished ===" >> "$REPORT_FILE"
cat "$REPORT_FILE"
