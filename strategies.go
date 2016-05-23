package chaosmonkey

// These are the default chaos strategies supported by Chaos Monkey. Each
// strategy "breaks" an EC2 instance in a different way, to simulate different
// types of failures.
//
// Source: https://github.com/Netflix/SimianArmy/wiki/The-Chaos-Monkey-Army
const (
	// StrategyShutdownInstance shuts down the instance using the EC2 API.
	// This is the classic chaos monkey strategy.
	StrategyShutdownInstance = "ShutdownInstance"

	// StrategyBlockAllNetworkTraffic removes all security groups from the
	// instance, and moves it into a security group that does not allow any
	// access. Thus the instance is running, but cannot be reached via the
	// network. This can only work on VPC instances.
	StrategyBlockAllNetworkTraffic = "BlockAllNetworkTraffic"

	// StrategyDetachVolumes force-detaches all EBS volumes from the
	// instance, simulating an EBS failure. Thus the instance is running,
	// but EBS disk I/O will fail.
	StrategyDetachVolumes = "DetachVolumes"

	// StrategyBurnCPU runs CPU intensive processes, simulating a noisy
	// neighbor or a faulty CPU. The instance will effectively have a much
	// slower CPU. Requires SSH to be configured.
	StrategyBurnCPU = "BurnCpu"

	// StrategyBurnIO runs disk intensive processes, simulating a noisy
	// neighbor or a faulty disk. The instance will effectively have a much
	// slower disk. Requires SSH to be configured.
	StrategyBurnIO = "BurnIo"

	// StrategyKillProcesses kills any Java or Python programs it finds
	// every second, simulating a faulty application, corrupted
	// installation or faulty instance. The instance is fine, but the
	// Java/Python application running on it will fail continuously.
	// Requires SSH to be configured.
	StrategyKillProcesses = "KillProcesses"

	// StrategyNullRoute null-routes the 10.0.0.0/8 network, which is used
	// by the EC2 internal network. All EC2 <-> EC2 network traffic will
	// fail. Requires SSH to be configured.
	StrategyNullRoute = "NullRoute"

	// StrategyFailEC2 puts dummy host entries into /etc/hosts, so that all
	// EC2 API communication will fail. This simulates a failure of the EC2
	// control plane. Of course, if your application doesn't use the EC2
	// API from the instance, then it will be completely unaffected.
	// Requires SSH to be configured.
	StrategyFailEC2 = "FailEc2"

	// StrategyFailDNS uses iptables to block port 53 for TCP & UDP; those
	// are the DNS traffic ports. This simulates a failure of your DNS
	// servers. Requires SSH to be configured.
	StrategyFailDNS = "FailDns"

	// StrategyFailDynamoDB puts dummy host entries into /etc/hosts, so
	// that all DynamoDB communication will fail. This simulates a failure
	// of DynamoDB. As with its monkey relatives, this will only affect
	// instances that use DynamoDB. Requires SSH to be configured.
	StrategyFailDynamoDB = "FailDynamoDb"

	// StrategyFailS3 puts dummy host entries into /etc/hosts, so that all
	// S3 communication will fail. This simulates a failure of S3. Of
	// course, if your application doesn't use S3, then it will be
	// completely unaffected. Requires SSH to be configured.
	StrategyFailS3 = "FailS3"

	// StrategyFillDisk writes a huge file to the root device, filling up
	// the (typically relatively small) EC2 root disk. Requires SSH to be
	// configured.
	StrategyFillDisk = "FillDisk"

	// StrategyNetworkCorruption uses the traffic shaping API to corrupt a
	// large fraction of network packets. This simulates degradation of the
	// EC2 network. Requires SSH to be configured.
	StrategyNetworkCorruption = "NetworkCorruption"

	// StrategyNetworkLatency uses the traffic shaping API to introduce
	// latency (1 second +- 50%) to all network packets. This simulates
	// degradation of the EC2 network. Requires SSH to be configured.
	StrategyNetworkLatency = "NetworkLatency"

	// StrategyNetworkLoss uses the traffic shaping API to drop a fraction
	// of all network packets. This simulates degradation of the EC2
	// network. Requires SSH to be configured.
	StrategyNetworkLoss = "NetworkLoss"
)

// Strategies is a list of default chaos strategies supported by Chaos Monkey.
var Strategies = []ChaosStrategy{
	StrategyShutdownInstance,
	StrategyBlockAllNetworkTraffic,
	StrategyDetachVolumes,
	StrategyBurnCPU,
	StrategyBurnIO,
	StrategyKillProcesses,
	StrategyNullRoute,
	StrategyFailEC2,
	StrategyFailDNS,
	StrategyFailDynamoDB,
	StrategyFailS3,
	StrategyFillDisk,
	StrategyNetworkCorruption,
	StrategyNetworkLatency,
	StrategyNetworkLoss,
}
