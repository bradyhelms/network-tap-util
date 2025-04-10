#!/bin/bash

# Default kj
PCAP_FILE="capture.pcap"

source claude-env/bin/activate

if [ ! -z "$1" ]; then
  PCAP_FILE="$1"
  echo "Using specified PCAP file: ${PCAP_FILE}"
else
  echo "Using PCAP file: ${PCAP_FILE}"
fi

echo "Checking if PCAP file exists..."
if [ ! -f "${PCAP_FILE}" ]; then
  echo "Error: PCAP file '${PCAP_FILE}' not found."
  exit 1
fi

echo "PCAP file found."

echo "Starting analysis with summary.py..."
echo "Sending data to Anthropic API..."
python3 summary.py -j "${PCAP_FILE}" | python3 claude_testing.py

if [ $? -eq 0 ]; then
  echo "Processing completed successfully."
else
  echo "Processing completed with errors."
fi
