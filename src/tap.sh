#!/bin/bash

usage() {
    echo "Usage: $0 [-h | --help] [-p | --program] [-c | --capture] [-a | --all]"
    echo ""
    echo "-h | --help      Show this dialog"
    echo "-p | --program   Start Vivado hw_server and program the FPGA board."
    echo "-c | --capture   Run a packet capture."
    echo "-a | --all       Start hw_server, program the board, and run a packet capture."
    exit 0
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -p|--program)
            echo "Starting the hardware server and programming the FPGA board..."
            go run main.go
            exit 0
            ;;
        -c|--capture)
            echo "Starting packet capture."
            go run capture/capture.go
            exit 0
            ;;
        -a|--all)
            go run main.go
            go run capture/capture.go
            exit 0
            ;;
        *)
            echo "Invalid option: $1"
            usage
            ;;
    esac
    shift
done

# Print usage if no args provided
usage
