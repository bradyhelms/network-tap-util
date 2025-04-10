import anthropic
import sys
import re
import ast
import subprocess

input_json_data = sys.stdin.read()
prompt = """
System Prompt:

You are a cybersecurity analyst AI that specializes in network traffic analysis and threat detection. You will be given JSON-formatted summary statistics derived from packet captures. These include details such as common protocols, source and destination IPs, ports, timestamps, and traffic patterns.

Assume that the traffic is malicious in nature. Your task is to:

    Analyze the data and hypothesize what kind of attack(s) might be occurring (e.g., DDoS, port scan, C2 beaconing, exfiltration, MITM).

    Justify your conclusions based on patterns in the input data.

    Suggest specific remediation steps that a network defender or SOC team could take to mitigate the threat (e.g., firewall rules, IDS signatures, endpoint isolation, further investigation).

Be precise and concise. Avoid generic recommendations. Focus on threat-specific insight and actionable defense strategies.
"""

client = anthropic.Anthropic()

message = client.messages.create(
    model="claude-3-haiku-20240307",
    max_tokens=1000,
    temperature=0,
    system=prompt,
    messages=[
        {
            "role":"user",
            "content": input_json_data
        }
    ]
)

def display_with_less(text: str):
    pager = subprocess.Popen(['less', '-R'], stdin=subprocess.PIPE)
    try:
        pager.communicate(input=text.encode('utf-8'))
    except KeyboardInterrupt:
        pager.terminate()

raw_output = str(message.content) 

match = re.search(r"text='(.*?)',\s*type=", raw_output)
if match:
    escaped_text = match.group(1)
    clean_text = ast.literal_eval(f"'{escaped_text}'")
    display_with_less(clean_text)
else:
    # fallback
    print(raw_output.replace("\\n", "\n"))


