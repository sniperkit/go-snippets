#include <linux/module.h>
#include <linux/version.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kprobes.h>
#include <linux/kallsyms.h>
#include <linux/kernel.h>
#include <linux/jiffies.h>
#include <linux/module.h>
#include <linux/interrupt.h>
#include <linux/fs.h>
#include <linux/types.h>
#include <linux/string.h>
#include <linux/socket.h>
#include <linux/errno.h>
#include <linux/fcntl.h>
#include <linux/in.h>

#include <linux/uaccess.h>
#include <asm/io.h>

#include <linux/inet.h>
#include <linux/netdevice.h>
#include <linux/etherdevice.h>
#include <linux/skbuff.h>
#include <linux/ethtool.h>
#include <net/sock.h>
#include <net/checksum.h>
#include <linux/if_ether.h>	/* For the statistics structure. */
#include <linux/if_arp.h>	/* For ARPHRD_ETHER */
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/percpu.h>
#include <net/net_namespace.h>
#include <linux/u64_stats_sync.h>
#include <linux/sched.h>
#include <linux/delay.h>
#include <linux/types.h>
#include <linux/workqueue.h>
#include <linux/netlink.h>
#include <linux/net_dropmon.h>
#include <linux/percpu.h>
#include <linux/timer.h>
#include <uapi/linux/if_arp.h>
#include <linux/skbuff.h>
#include <linux/netdevice.h>


struct pcpu_lstats {
	u64			packets;
	u64			bytes;
	struct u64_stats_sync	syncp;
};


struct timer_list	send_timer;
struct work_struct	send_work;
struct sk_buff *packet_create(struct net_device *dev);

/*
 * The higher levels take care of making this non-reentrant (it's
 * called with bh's disabled).
 */
static netdev_tx_t kkkmmu_xmit(struct sk_buff *skb,
				 struct net_device *dev)
{
	struct pcpu_lstats *lb_stats;
	int len;

	printk(KERN_EMERG "Sending Pakcet from kkkmmu\n");
	skb_orphan(skb);

	/* Before queueing this packet to netif_rx(),
	 * make sure dst is refcounted.
	 */
	skb_dst_force(skb);

	skb->protocol = eth_type_trans(skb, dev);

	/* it's OK to use per_cpu_ptr() because BHs are off */
	lb_stats = this_cpu_ptr(dev->lstats);

	len = skb->len;
	if (likely(netif_rx(skb) == NET_RX_SUCCESS)) {
		u64_stats_update_begin(&lb_stats->syncp);
		lb_stats->bytes += len;
		lb_stats->packets++;
		u64_stats_update_end(&lb_stats->syncp);
	}

	return NETDEV_TX_OK;
}

static void kkkmmu_get_stats64(struct net_device *dev,
				 struct rtnl_link_stats64 *stats)
{
	u64 bytes = 0;
	u64 packets = 0;
	int i;

	for_each_possible_cpu(i) {
		const struct pcpu_lstats *lb_stats;
		u64 tbytes, tpackets;
		unsigned int start;

		lb_stats = per_cpu_ptr(dev->lstats, i);
		do {
			start = u64_stats_fetch_begin_irq(&lb_stats->syncp);
			tbytes = lb_stats->bytes;
			tpackets = lb_stats->packets;
		} while (u64_stats_fetch_retry_irq(&lb_stats->syncp, start));
		bytes   += tbytes;
		packets += tpackets;
	}
	stats->rx_packets = packets;
	stats->tx_packets = packets;
	stats->rx_bytes   = bytes;
	stats->tx_bytes   = bytes;
}

static u32 always_on(struct net_device *dev)
{
	return 1;
}

static const struct ethtool_ops kkkmmu_ethtool_ops = {
	.get_link		= always_on,
};

static int kkkmmu_dev_init(struct net_device *dev)
{
	dev->lstats = netdev_alloc_pcpu_stats(struct pcpu_lstats);
	if (!dev->lstats)
		return -ENOMEM;
	return 0;
}

static void kkkmmu_dev_free(struct net_device *dev)
{
	free_percpu(dev->lstats);
	free_netdev(dev);
}

static const struct net_device_ops kkkmmu_ops = {
	.ndo_init      = kkkmmu_dev_init,
	.ndo_start_xmit= kkkmmu_xmit,
	.ndo_get_stats64 = kkkmmu_get_stats64,
	.ndo_set_mac_address = eth_mac_addr,
};

int kkkmmu_header(struct sk_buff *skb, struct net_device *dev,
	       unsigned short type,
	       const void *daddr, const void *saddr, unsigned int len)
{
	struct ethhdr *eth = skb_push(skb, ETH_HLEN);

	if (type != ETH_P_802_3 && type != ETH_P_802_2)
		eth->h_proto = htons(type);
	else
		eth->h_proto = htons(len);

	/*
	 *      Set the source hardware address.
	 */

	if (!saddr)
		saddr = dev->dev_addr;
	memcpy(eth->h_source, saddr, ETH_ALEN);

	if (daddr) {
		memcpy(eth->h_dest, daddr, ETH_ALEN);
		return ETH_HLEN;
	}

	/*
	 *      Anyway, the loopback-device should never use this function...
	 */

	if (dev->flags & (IFF_LOOPBACK | IFF_NOARP)) {
		eth_zero_addr(eth->h_dest);
		return ETH_HLEN;
	}

	return -ETH_HLEN;
}

int kkkmmu_header_parse(const struct sk_buff *skb, unsigned char *haddr)
{
	const struct ethhdr *eth = eth_hdr(skb);
	memcpy(haddr, eth->h_source, ETH_ALEN);
	return ETH_ALEN;
}

int kkkmmu_header_cache(const struct neighbour *neigh, struct hh_cache *hh, __be16 type)
{
	struct ethhdr *eth;
	const struct net_device *dev = neigh->dev;

	eth = (struct ethhdr *)
	    (((u8 *) hh->hh_data) + (HH_DATA_OFF(sizeof(*eth))));

	if (type == htons(ETH_P_802_3))
		return -1;

	eth->h_proto = type;
	memcpy(eth->h_source, dev->dev_addr, ETH_ALEN);
	memcpy(eth->h_dest, neigh->ha, ETH_ALEN);
	hh->hh_len = ETH_HLEN;
	return 0;
}

void kkkmmu_header_cache_update(struct hh_cache *hh,
			     const struct net_device *dev,
			     const unsigned char *haddr)
{
	memcpy(((u8 *) hh->hh_data) + HH_DATA_OFF(sizeof(struct ethhdr)),
	       haddr, ETH_ALEN);
}

const struct header_ops kkkmmu_header_ops ____cacheline_aligned = {
	.create		= kkkmmu_header,
	.parse		= kkkmmu_header_parse,
	.cache		= kkkmmu_header_cache,
	.cache_update	= kkkmmu_header_cache_update,
};

static void kkkmmu_setup(struct net_device *dev)
{
	dev->mtu		= 64 * 1024;
	dev->hard_header_len	= ETH_HLEN;	/* 14	*/
	dev->min_header_len	= ETH_HLEN;	/* 14	*/
	dev->addr_len		= ETH_ALEN;	/* 6	*/
	dev->type		= ARPHRD_LOOPBACK;	/* 0x0001*/
	dev->flags		= IFF_LOOPBACK;
	dev->priv_flags		|= IFF_LIVE_ADDR_CHANGE | IFF_NO_QUEUE;
	netif_keep_dst(dev);
	dev->hw_features	= NETIF_F_GSO_SOFTWARE;
	dev->features 		= NETIF_F_SG | NETIF_F_FRAGLIST
		| NETIF_F_GSO_SOFTWARE
		| NETIF_F_HW_CSUM
		| NETIF_F_RXCSUM
		| NETIF_F_SCTP_CRC
		| NETIF_F_HIGHDMA
		| NETIF_F_LLTX
		| NETIF_F_NETNS_LOCAL
		| NETIF_F_VLAN_CHALLENGED;
	dev->ethtool_ops	= &kkkmmu_ethtool_ops;
	dev->header_ops		= &kkkmmu_header_ops;
	dev->netdev_ops		= &kkkmmu_ops;
	//dev->destructor		= kkkmmu_dev_free;
}
/* Setup and register a new device. */
static int kkkmmu_init(struct net *net)
{
	struct net_device *dev;
	int err;

	err = -ENOMEM;
	dev = alloc_netdev(0, "kkkmmu", NET_NAME_UNKNOWN, kkkmmu_setup);
	if (!dev)
		goto out;

	dev_net_set(dev, net);
	err = register_netdev(dev);
	if (err)
		goto out_free_netdev;

	return 0;


out_free_netdev:
	free_netdev(dev);
out:
	if (net_eq(net, &init_net))
		panic("kkkmmu: Failed to register netdevice: %d\n", err);
	return err;
}

static void sched_send_work(unsigned long _data)
{
	struct net_device *dev = NULL;
	struct sk_buff *skb = NULL;

	dev = __dev_get_by_name(&init_net, "kkkmmu");
	if (dev == NULL) {
		printk(KERN_EMERG "Cannot find netdevice kkkmmu\n");
		return;
	}

	skb = packet_create(dev);
	if (skb == NULL) {
		printk(KERN_EMERG "Cannot create new packet\n");
		return;
	}
	netif_rx(skb);

	mod_timer(&send_timer, jiffies + HZ * 5);
}

struct sk_buff *packet_create(struct net_device *dev)
{

	struct sk_buff *skb;
	struct arphdr *arp;
	unsigned char *arp_ptr;
	int hlen = LL_RESERVED_SPACE(dev);
	int tlen = dev->needed_tailroom;

	unsigned char dest_hw[6] = {0xff, 0xff, 0xff, 0xff, 0xff};
	unsigned char target_hw[6] = {0xff, 0xff, 0xff, 0xff, 0xff};
	unsigned char src_hw[6] = {0xff, 0xff, 0xff, 0xff, 0xff};
	__be32 dest_ip = in_aton("10.0.0.1");
	__be32 src_ip = in_aton("10.0.0.2");
	int ptype = ETH_P_ARP;
	int type = ARPOP_REQUEST;	
	/*
	 *	Allocate a buffer
	 */

	skb = alloc_skb(arp_hdr_len(dev) + hlen + tlen, GFP_ATOMIC);
	if (!skb)
		return NULL;

	skb_reserve(skb, hlen);
	skb_reset_network_header(skb);
	arp = skb_put(skb, arp_hdr_len(dev));
	skb->dev = dev;
	skb->protocol = htons(ETH_P_ARP);

	/*
	 *	Fill the device header for the ARP frame
	 */
	if (dev_hard_header(skb, dev, ptype, dest_hw, src_hw, skb->len) < 0)
		goto out;

	/*
	 * Fill out the arp protocol part.
	 *
	 * The arp hardware type should match the device type, except for FDDI,
	 * which (according to RFC 1390) should always equal 1 (Ethernet).
	 */
	/*
	 *	Exceptions everywhere. AX.25 uses the AX.25 PID value not the
	 *	DIX code for the protocol. Make these device structure fields.
	 */
	switch (dev->type) {
	default:
		arp->ar_hrd = htons(dev->type);
		arp->ar_pro = htons(ETH_P_IP);
		break;

#if IS_ENABLED(CONFIG_AX25)
	case ARPHRD_AX25:
		arp->ar_hrd = htons(ARPHRD_AX25);
		arp->ar_pro = htons(AX25_P_IP);
		break;

#if IS_ENABLED(CONFIG_NETROM)
	case ARPHRD_NETROM:
		arp->ar_hrd = htons(ARPHRD_NETROM);
		arp->ar_pro = htons(AX25_P_IP);
		break;
#endif
#endif

#if IS_ENABLED(CONFIG_FDDI)
	case ARPHRD_FDDI:
		arp->ar_hrd = htons(ARPHRD_ETHER);
		arp->ar_pro = htons(ETH_P_IP);
		break;
#endif
	}

	arp->ar_hln = dev->addr_len;
	arp->ar_pln = 4;
	arp->ar_op = htons(type);

	arp_ptr = (unsigned char *)(arp + 1);

	memcpy(arp_ptr, src_hw, dev->addr_len);
	arp_ptr += dev->addr_len;
	memcpy(arp_ptr, &src_ip, 4);
	arp_ptr += 4;

	switch (dev->type) {
#if IS_ENABLED(CONFIG_FIREWIRE_NET)
	case ARPHRD_IEEE1394:
		break;
#endif
	default:
		if (target_hw)
			memcpy(arp_ptr, target_hw, dev->addr_len);
		else
			memset(arp_ptr, 0, dev->addr_len);
		arp_ptr += dev->addr_len;
	}
	memcpy(arp_ptr, &dest_ip, 4);

	return skb;

out:
	kfree_skb(skb);
	return NULL;
}

static int __init net_dev_init(void)
{
	int ret = -1;
	printk("Test net_dev module init\n");
	kkkmmu_init(&init_net);
	init_timer(&send_timer);
	send_timer.data = (unsigned long)10000;
	send_timer.function = sched_send_work;
	send_timer.expires = jiffies + 1 * HZ;
	add_timer(&send_timer);
	return 0;
}

static void __exit net_dev_exit(void)
{
	printk("Test net_dev module removed\n");
}

module_init(net_dev_init);
module_exit(net_dev_exit);

MODULE_AUTHOR("kkkmmu");
MODULE_DESCRIPTION("netdevice");
MODULE_LICENSE("GPL");

