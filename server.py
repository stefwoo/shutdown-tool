import http.server
import socketserver
import json
import subprocess
import os
import sys
from urllib.parse import urlparse

# Load configuration
CONFIG_FILE = 'config.json'

def load_config():
    try:
        with open(CONFIG_FILE, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"Config file '{CONFIG_FILE}' not found. Creating default...")
        default_config = {
            "port": 8080,
            "commands": {
                "shutdown": "shutdown /s /t 0",
                "sleep": "rundll32.exe powrprof.dll,SetSuspendState 0,1,0"
            }
        }
        with open(CONFIG_FILE, 'w') as f:
            json.dump(default_config, f, indent=2)
        return default_config

config = load_config()
PORT = config.get('port', 8080)
COMMANDS = config.get('commands', {})

class CommandHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path_parts = parsed_path.path.strip('/').split('/')

        # Handle /execute/{command}
        if len(path_parts) == 2 and path_parts[0] == 'execute':
            cmd_name = path_parts[1]
            if cmd_name in COMMANDS:
                cmd_str = COMMANDS[cmd_name]
                print(f"Executing command: {cmd_name} -> {cmd_str}")
                
                try:
                    # Execute command
                    # shell=True allows executing shell commands directly
                    subprocess.Popen(cmd_str, shell=True)
                    
                    self.send_response(200)
                    self.send_header('Content-type', 'text/plain; charset=utf-8')
                    self.end_headers()
                    self.wfile.write(f"Command '{cmd_name}' executed successfully.".encode('utf-8'))
                except Exception as e:
                    self.send_response(500)
                    self.send_header('Content-type', 'text/plain; charset=utf-8')
                    self.end_headers()
                    self.wfile.write(f"Error executing command: {str(e)}".encode('utf-8'))
            else:
                self.send_response(404)
                self.send_header('Content-type', 'text/plain; charset=utf-8')
                self.end_headers()
                self.wfile.write(f"Command '{cmd_name}' not found.".encode('utf-8'))
        
        # Handle root /
        elif self.path == '/':
            self.send_response(200)
            self.send_header('Content-type', 'text/html; charset=utf-8')
            self.end_headers()
            html = "<html><body><h1>Shutdown Tool is Running</h1>"
            html += "<p>Available commands:</p><ul>"
            for cmd in COMMANDS:
                html += f"<li><a href='/execute/{cmd}'>{cmd}</a></li>"
            html += "</ul></body></html>"
            self.wfile.write(html.encode('utf-8'))
            
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == "__main__":
    try:
        with socketserver.TCPServer(("", PORT), CommandHandler) as httpd:
            print(f"Serving on port {PORT}")
            print(f"Available commands: {list(COMMANDS.keys())}")
            httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped.")
