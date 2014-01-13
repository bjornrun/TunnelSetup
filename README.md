TunnelSetup
===========

Setup dynamic SSH tunnels 


Usage:


	%./TunnelSetup
	Tunnel Setup
	-c="tunnels.cfg": Tunnel config setup file
	-e="help": Execute command: help|attach|detatch|config|forward <local port:ip:remote port>|remote <remote port:ip:local port>

	%./TunnelSetup -s -e attach
	Tunnel Setup
	bjorn@10.0.1.136's password:
	Server 10.0.1.136 is now attached

	%./TunnelSetup -s -e forward 10022:10.0.1.136:22
	Tunnel Setup
	Forward tunnel 10022:10.0.1.136:22 active

	%./TunnelSetup -s -e remote 10022:10.0.1.4:22
	Tunnel Setup
	Remote tunnel 10022:10.0.1.4:22 active

	%./TunnelSetup.go -s -e autoforward 127.0.0.1:80
	Tunnel Setup
	Forward tunnel 10000:127.0.0.1:80 active

	%./TunnelSetup -s -e config
	Tunnel Setup
	Configuration:
	Instance: 0
	Server: 10.0.1.136:10080
	SOCKS server on localhost port 10080
	Attached to Proxy 10.0.1.136
	Tunnels:
	Forward 10022:10.0.1.136:22
	Remote 10022:10.0.1.4:22
	Forward 10000:127.0.0.1:80

	%./TunnelSetup -e detach
	Tunnel Setup
	Stop listening request sent.
	Server 10.0.1.136 is now detached


FAQ:

Where can I find a later version of ssh for acient SuSE (e.g. Version 11p2) so it handles multiplexed forward?

You download it from here:
http://download.opensuse.org/repositories/home:/H4T:/branches:/network/SLE_11_SP2/x86_64/openssh-6.2p2-68.1.x86_64.rpm

Then you run
<pre>
rpm2cpio openssh-6.2p2-68.1.x86_64.rpm | cpio -idmv
</pre>

In tunnels.cfg you set ssh parameter to the full path to the ssh binary you have extracted.
