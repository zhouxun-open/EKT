package param

import (
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
)

// MainNet后面将通过testnet投票产生
var MainNet = []p2p.Peer{}

// TestNet，测试网内测和公测的DPoS节点
var TestNet = []p2p.Peer{
	{"df37b40e83debb2943f46432d0d2983cd52ccaa44ce16f6dae48672f08fb704d", "58.83.148.228", 19951, 4, ""},
	//{"8b020bb6af84b7c5044393d816309894e325d05c65fd0fd006587d705bd5cb17", "58.83.148.229", 19951, 4, ""},
	//{"38dadeee1dfa41ea35e37d050bd79540d641f1f7d4312c190aad1ebd1ba27186", "58.83.148.230", 19951, 4, ""},
	//{"cf881598ddb6b9b83e87e7e44a43d14eee73d099555f6307261ab0430815eff4", "119.23.14.56", 19951, 4, ""},
	//{"38dadeee1dfa41ea35e37d050bd79540d641f1f7d4312c190aad1ebd1ba27186", "139.224.133.92", 19951, 4, ""},
	//{"2f7174abff25cedd80cb4278a0e7553a3a97247a8eeaf20489e9d1b072aeece6", "45.76.50.87", 19951, 4, ""},
}
