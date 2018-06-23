#include <linux/module.h>
#include <linux/version.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/kprobes.h>
#include <linux/kallsyms.h>
#include <linux/inet.h>
#include <linux/skbuff.h>
#include <net/ip.h>
#include <linux/types.h>
#include <linux/ip.h>
#include <linux/netdevice.h>

static unsigned int counter = 0;

static long pre_ip_rcv(struct sk_buff *skb, struct net_device *dev, struct packet_type *pt)
{
	printk(KERN_EMERG "[%s]Received packet from %s\n", "ip_rcv", dev->name);
	jprobe_return();
	return 0;
}

static long pre_netif_rx(struct sk_buff *skb)
{
	printk(KERN_EMERG "[%s]Received packet from %s\n", "netif_rx", skb->dev->name);
	jprobe_return();
	return 0;
}

static long pre_dev_queue_xmit(struct sk_buff *skb)
{
	    printk(KERN_EMERG "[%s]Send Packet from: %s\n", "dev_queue_xmit", skb->dev->name);
		jprobe_return();
		return 0;
}

static long pre_arp_rcv(struct sk_buff *skb, struct net_device *dev,
		    struct packet_type *pt, struct net_device *orig_dev)
{
	    printk(KERN_EMERG "[%s] Received packet from : %s\n", "arp_rcv", skb->dev->name);
	        jprobe_return();
		    return 0;
}

static struct jprobe jprobe_netif_rx = 
{
	    .entry = pre_netif_rx,
	        .kp = {
				.symbol_name = "netif_rx",
				    },
};

static struct jprobe jprobe_arp_rcv = 
{
	    .entry = pre_arp_rcv,
	        .kp = {
				.symbol_name = "arp_rcv",
				    },
};

static struct jprobe jprobe_ip_rcv = 
{
	    .entry = pre_ip_rcv,
	        .kp = {
				.symbol_name = "ip_rcv",
				    },
};

static struct jprobe jprobe_dev_queue_xmit = 
{
	    .entry = pre_dev_queue_xmit,
	        .kp = {
				.symbol_name = "dev_queue_xmit",
				    },
};

static int __init jp_init(void)
{
	int ret = -1;
	printk("Test jp module init\n");

	ret = register_jprobe(&jprobe_netif_rx);
	ret = register_jprobe(&jprobe_ip_rcv);
	ret = register_jprobe(&jprobe_arp_rcv);
	ret = register_jprobe(&jprobe_dev_queue_xmit);
	if (ret < 0)
	{
		printk (KERN_EMERG "Register jprobe failed, return: %x\n", ret);
		return -1;
	}

#if 0
	printk(KERN_EMERG "Planted jprobe at: %p, handler addr: %p\n", 
			jprobe_do_fork.kp.addr, jprobe_do_fork.entry);
#endif

	return 0;
}

static void __exit jp_exit(void)
{
	unregister_jprobe(&jprobe_arp_rcv);
	unregister_jprobe(&jprobe_netif_rx);
	unregister_jprobe(&jprobe_ip_rcv);
	unregister_jprobe(&jprobe_dev_queue_xmit);
	printk("Test jp module removed\n");
}

module_init(jp_init);
module_exit(jp_exit);

MODULE_AUTHOR("kkkmmu");
MODULE_DESCRIPTION("Jprobe_Mudule");
MODULE_LICENSE("GPL");

