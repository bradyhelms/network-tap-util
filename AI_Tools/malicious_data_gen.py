import socket
import random
import time
import threading
import ipaddress

def generate_random_ip():
    # Generate a random public IP address (avoiding private ranges)
    while True:
        ip = str(ipaddress.IPv4Address(random.randint(1, 2**32-1)))
        # Skip private ranges
        if not ipaddress.ip_address(ip).is_private:
            return ip

def port_scan_simulation(target_ip):
    """Simulate a port scanning pattern"""
    common_ports = [21, 22, 23, 25, 80, 139, 445, 3389, 8080, 8443]
    for port in random.sample(common_ports, 5):  # Scan 5 random ports
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.settimeout(0.1)  # Very short timeout
            s.connect((target_ip, port))
            s.close()
        except:
            pass
        time.sleep(0.05)  # Small delay between port attempts

def bruteforce_simulation(target_ip):
    """Simulate a login brute force attempt"""
    auth_ports = [22, 23, 3389]  # SSH, Telnet, RDP
    port = random.choice(auth_ports)
    for _ in range(5):  # Send multiple login attempts
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.settimeout(0.5)
            s.connect((target_ip, port))
            # Send random "authentication" data
            s.send(random.randbytes(random.randint(32, 64)))
            s.close()
        except:
            pass
        time.sleep(0.2)  # Delay between login attempts

def dos_simulation(target_ip):
    """Simulate a high-volume traffic pattern"""
    port = random.choice([80, 443])
    for _ in range(10):  # Send burst of packets
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.settimeout(0.3)
            s.connect((target_ip, port))
            # Send HTTP-like request with large payload
            s.send(b"GET / HTTP/1.1\r\nHost: target\r\n" + b"X-Data: " + b"A" * 1000 + b"\r\n\r\n")
            s.close()
        except:
            pass

def generate_malicious_traffic(duration=60, intensity=1):
    """
    Generate simulated malicious traffic patterns
    
    Args:
        duration: Duration to run in seconds
        intensity: 1-5 scale of traffic intensity (higher = more traffic)
    """
    start_time = time.time()
    print(f"Generating simulated malicious traffic patterns for {duration} seconds...")
    
    local_targets = ["127.0.0.1"]  # Default to localhost
    
    attack_types = [port_scan_simulation, bruteforce_simulation, dos_simulation]
    
    while time.time() - start_time < duration:
        attack = random.choice(attack_types)
        target = random.choice(local_targets)
        
        # Launch attack simulation in a thread
        thread = threading.Thread(target=attack, args=(target,))
        thread.start()
        
        # Wait between attack patterns based on intensity (lower = longer wait)
        wait_time = 5 / intensity
        time.sleep(wait_time)
    
    print("Traffic generation completed")

if __name__ == "__main__":
    generate_malicious_traffic(duration=30, intensity=2)
