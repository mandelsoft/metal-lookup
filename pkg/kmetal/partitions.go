package kmetal

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gardener/controller-manager-library/pkg/logger"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-go/api/models"
)

type Partition struct {
	*models.V1PartitionResponse
	Networks []*net.IPNet
}

type Partitions struct {
	lock       sync.RWMutex
	driver     *metalgo.Driver
	ttl        time.Duration
	next       time.Time
	partitions map[string]*Partition
}

func NewPartitions(driver *metalgo.Driver, ttl time.Duration) *Partitions {
	return &Partitions{driver: driver, ttl: ttl}
}

func (this *Partitions) Update() error {
	if this.requireUpdate() {
		return this.UpdateCache(true)
	}
	return nil
}

func (this *Partitions) requireUpdate() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.partitions == nil || time.Now().After(this.next)
}

func (this *Partitions) UpdateCache(check ...bool) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	now := time.Now()
	if this.partitions != nil && !now.After(this.next) {
		for _, b := range check {
			if b {
				return nil
			}
		}
	}
	plist, err := this.driver.PartitionList()
	if err != nil {
		return fmt.Errorf("cannot get partition list: %s", err)
	}

	partitions := map[string]*Partition{}
	if plist != nil {
		for _, p := range plist.Partition {
			partitions[s(p.ID)] = &Partition{
				V1PartitionResponse: p,
			}
		}
	}

	nlist, err := this.driver.NetworkList()
	if err != nil {
		return err
	}
	if nlist != nil {
		for _, n := range nlist.Networks {
			/*
			 */
			fmt.Printf("%s: parent: %s, partition: %s, priv: %t, under: %t, nat: %t, prefix: %+v\n", s(n.ID), s(n.Parentnetworkid),
				n.Partitionid,
				b(n.Privatesuper), b(n.Underlay), b(n.Nat),
				n.Prefixes)
			if n.Parentnetworkid == nil && b(n.Nat) && !b(n.Privatesuper) && !b(n.Underlay) {
				p := partitions[n.Partitionid]
				if p != nil {
					for _, c := range n.Prefixes {
						_, cidr, err := net.ParseCIDR(c)
						if err == nil {
							p.Networks = append(p.Networks, cidr)
						}
					}
				}
			}
		}
	}

	this.partitions = partitions
	this.next = now.Add(this.ttl)
	return nil
}

func (this *Partitions) LookupForRequester(ips ...net.IP) (*Partition, error) {
	err := this.Update()
	if err != nil {
		return nil, err
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, p := range this.partitions {
		for _, c := range p.Networks {
			for _, ip := range ips {
				if c.Contains(ip) {
					return p, nil
				}
			}
		}
	}
	return nil, nil
}

func (this *Partitions) Lookup(search string) (*Partition, error) {
	err := this.Update()
	if err != nil {
		return nil, err
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	for id, p := range this.partitions {
		if id == search {
			return p, nil
		}
		if p.Name == search {
			return p, nil
		}
	}
	return nil, nil
}

func FillPartitionMetadata(logger logger.LogContext, m *Partition, metadata map[string]interface{}) {
	partition := map[string]interface{}{}
	set(CMDLINE, m.Bootconfig.Commandline, partition)
	set(KERNEL, m.Bootconfig.Kernelurl, partition)
	set(INITRD, m.Bootconfig.Imageurl, partition)
	set(NAME, m.Name, partition)
	set(ID, m.ID, partition)
	set(PARTITION_IN, m.ID, metadata)
	if len(partition) > 0 {
		metadata[PARTITION] = partition
	}
}
