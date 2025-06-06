-- Example hosts for iperf3 testing
-- Run these SQL commands after starting the application to add sample hosts
-- You can execute these via the API or directly in the SQLite database

-- Example LAN hosts (replace with your actual LAN IPs)
INSERT INTO hosts (name, hostname, port, type, active, description) VALUES 
('Router', '192.168.1.1', 5201, 'lan', true, 'Home router with iperf3 server'),
('NAS Server', '192.168.1.100', 5201, 'lan', true, 'Network attached storage'),
('Desktop PC', '192.168.1.50', 5201, 'lan', true, 'Main desktop computer');

-- Example VPN hosts (replace with your VPN server IPs)
INSERT INTO hosts (name, hostname, port, type, active, description) VALUES 
('VPN Gateway', '10.0.0.1', 5201, 'vpn', true, 'VPN server gateway'),
('Remote Office', '10.0.1.100', 5201, 'vpn', true, 'Remote office server');

-- Example remote hosts (replace with your remote server IPs)
INSERT INTO hosts (name, hostname, port, type, active, description) VALUES 
('Cloud Server', 'your-server.example.com', 5201, 'remote', true, 'Cloud VPS with iperf3'),
('CDN Edge', '203.0.113.10', 5201, 'remote', true, 'CDN edge server');

-- Note: Make sure iperf3 is running in server mode on these hosts:
-- iperf3 -s -p 5201 -D 