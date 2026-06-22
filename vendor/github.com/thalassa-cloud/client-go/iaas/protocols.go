package iaas

type LoadbalancerProtocol string

const (
	ProtocolTCP   LoadbalancerProtocol = "tcp"
	ProtocolUDP   LoadbalancerProtocol = "udp"
	ProtocolHTTP  LoadbalancerProtocol = "http"
	ProtocolHTTPS LoadbalancerProtocol = "https"
	ProtocolGRPC  LoadbalancerProtocol = "grpc"
	ProtocolQUIC  LoadbalancerProtocol = "quic"
)

type FirewallAclProtocol string

const (
	FirewallAclProtocolTCP  FirewallAclProtocol = "tcp"
	FirewallAclProtocolUDP  FirewallAclProtocol = "udp"
	FirewallAclProtocolICMP FirewallAclProtocol = "icmp"
	FirewallAclProtocolIP   FirewallAclProtocol = "ip"
	FirewallAclProtocolAll  FirewallAclProtocol = "all"
)
