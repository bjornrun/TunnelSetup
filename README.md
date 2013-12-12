TunnelSetup
===========

Setup dynamic SSH tunnels 


Usage:


	%./TunnelSetup
	Tunnel Setup
	-c="tunnels.cfg": Tunnel config setup file
	-e="help": Execute command: help|attach|detatch|config|forward <local port:ip:remote port>|remote <remote port:ip:local port>

	%./TunnelSetup -e attach
	Tunnel Setup
	bjorn@10.0.1.136's password:
	Server 10.0.1.136 is now attached

	%./TunnelSetup -e forward 10022:10.0.1.136:22
	Tunnel Setup
	Forward tunnel 10022:10.0.1.136:22 active

	%./TunnelSetup -e remote 10022:10.0.1.4:22
	Tunnel Setup
	Remote tunnel 10022:10.0.1.4:22 active

	%./TunnelSetup -e config
	Tunnel Setup
	Configuration:
	Instance: 1
	Server: 10.0.1.136:10080
	Attached to Proxy 10.0.1.136
	Tunnels:
	Forward 10022:10.0.1.136:22
	Remote 10022:10.0.1.4:22

	%./TunnelSetup -e detach
	Tunnel Setup
	Detach Stop listening request sent.
	Server 10.0.1.136 is now detached


