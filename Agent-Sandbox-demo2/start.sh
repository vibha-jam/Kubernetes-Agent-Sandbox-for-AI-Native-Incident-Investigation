#!/bin/sh


ollama serve &

sleep 5

ollama pull qwen3.5:2b &

./ai-agent