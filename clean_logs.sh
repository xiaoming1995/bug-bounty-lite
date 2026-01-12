#!/bin/bash

# 日志目录路径
LOG_DIR="./logs"

if [ -d "$LOG_DIR" ]; then
    echo "清理目录: $LOG_DIR"
    rm -rf "$LOG_DIR"/*
    echo "[OK] 日志清理完成。"
else
    echo "[Error] 日志目录不存在。"
fi
